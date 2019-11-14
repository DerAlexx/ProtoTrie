package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/ob-vss-ws19/blatt-3-pwn/messages"
	"github.com/ob-vss-ws19/blatt-3-pwn/tree"
)

/*
ServerRemoteActor represents a RemoteActor in the service
*/
type ServerRemoteActor struct{}

/*
ID is the ID of a new created Trie
*/
type ID int

/*
Token is
*/
type Token string

/*
TrieContainer is a container for a Trie containing a Rootnode a Token and the PID of the Root
*/
type TrieContainer struct { //TODO
	Root    *actor.Props
	Token   Token
	Pid     *actor.PID
	Size    int
	Context *actor.RootContext
}

/*
RootNodes is a Map of rootnodes containing all roots --> Sorted by Key (Trie ID) and Value (Root Node --> actor.Actor)
*/
var RootNodes map[ID]TrieContainer = make(map[ID]TrieContainer)

/*
clientpid is the PID of the Client Remote Actor
*/
var clientpid actor.PID

/*
globalpid is the PID of the Service Remote Actor
*/
var globalpid actor.PID

/*
context the root actor of the service
*/
var context actor.RootContext

/*
contains will check whether a key is in the map or not
*/
func contains(preID int) bool {
	_, contains := RootNodes[ID(preID)]
	return contains
}

//TODO Comment
func containsByID(id ID) bool {
	_, contains := RootNodes[id]
	return contains
}

/*
getRandomID will return a random id which is not allready in the map
*/
func getRandomID(l int) int {
	for {
		potantialID := rand.Intn(l)
		if !contains(potantialID) {
			return potantialID
		}
	}
}

/*
getPID will return the PID of the rootnode of a given Trie-ID
*/
func getPID(id ID) *actor.PID {
	return RootNodes[id].Pid
}

/*
addInToRootsMap will add a trie into the map of roots with its id as a parameter
*/
func addInToRootsMap(id ID, trie TrieContainer) {
	RootNodes[id] = trie
}

//TODO Comment
func printMap() {
	for k, v := range RootNodes {
		fmt.Println(k, v)
	}
}

/*
DeleteTrie will delete a Trie
*/
func deleteTrie(id ID, token Token) bool {
	if MatchIDandToken(id, token) && containsByID(id) {
		delete(RootNodes, id)
		return true
	}
	return false
}

/*
Receive will receive different types of messages from the client. Each message is responsible for one type of action that the tree can execute (e.g. delete a Key, Create a Tree)
After the execution it will return a message to the client
*/
func (*ServerRemoteActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *messages.DeleteRequest:
		id := ID(int(msg.GetId()))
		rootpid := getPID(id)
		tok := Token(msg.GetToken())
		if MatchIDandToken(id, tok) {
			fmt.Println("Sending DeleteMessage to RootNode")
			context.RequestWithCustomSender(rootpid, tree.DeleteMessage{
				PID: *context.Sender(),
				Key: int(msg.GetKey()),
			}, &globalpid)
		} else {
			context.Respond(&messages.Response{
				SomeValue: "Wrong Combination of ID and Token",
			})
		}
	case *messages.ChangeRequest:
		pa := tree.Pair{
			Key:   int(msg.GetKey()),
			Value: msg.GetValue(),
		}
		id := ID(msg.GetId())
		rootpid := getPID(id)
		tok := Token(msg.GetToken())

		if MatchIDandToken(id, tok) {
			fmt.Println("Sending ChangeValueMessage to RootNode")
			context.RequestWithCustomSender(rootpid, tree.ChangeValueMessage{
				PID:     *context.Sender(),
				Element: pa,
			}, &globalpid)
		} else {
			context.Respond(&messages.Response{
				SomeValue: "Wrong Combination of ID and Token",
			})
		}
	case *messages.InsertRequest:
		pa := tree.Pair{
			Key:   int(msg.GetKey()),
			Value: msg.GetValue(),
		}
		id := ID(msg.GetId())
		rootpid := getPID(id)
		tok := Token(msg.GetToken())

		if MatchIDandToken(id, tok) {
			fmt.Println("Sending InsertMessage to RootNode")
			context.Send(rootpid, tree.InsertMessage{
				PID:        *context.Sender(),
				Element:    pa,
				PIDService: globalpid,
				PIDRoot:    *rootpid,
			})
		} else {
			context.Respond(&messages.Response{
				SomeValue: "Wrong Combination of ID and Token",
			})
		}
	case *messages.FindRequest:
		id := ID(msg.GetId())
		rootpid := getPID(id)
		tok := Token(msg.GetToken())

		if MatchIDandToken(id, tok) {
			fmt.Println("Sending FindMessage to RootNode")
			context.RequestWithCustomSender(rootpid, tree.FindMessage{
				PID: *context.Sender(),
				Key: int(msg.GetKey()),
			}, &globalpid)
		} else {
			context.Respond(&messages.Response{
				SomeValue: "Wrong Combination of ID and Token",
			})
		}
	case *messages.CreateTreeRequest:
		i, t, _ := AddNewTrie(int(msg.GetSize_()))
		printMap()
		fmt.Printf("The size of the RootNode Map before creating the trie is %d", len(RootNodes))
		context.Respond(&messages.Response{
			SomeValue: fmt.Sprintf("Your ID: %d, Your Token: %s", int(i), string(t)),
		})
		fmt.Printf("The size of the RootNode Map after creating the trie is %d", len(RootNodes))
	case *messages.DeleteTreeRequest:
		ret := deleteTrie(ID(msg.GetId()), Token(msg.GetToken()))
		fmt.Printf("The size of the RootNode Map before deleting the trie is %d", len(RootNodes))
		if ret {
			context.Respond(&messages.Response{
				SomeValue: "Success",
			})
		} else {
			context.Respond(&messages.Response{
				SomeValue: "Couldnt delete the tree",
			})
		}
		fmt.Printf("The size of the RootNode Map after deleting the trie is %d", len(RootNodes))
	case *tree.WantBasicNodeActorsMessage:
		size := msg.Size

		propleft := actor.PropsFromProducer(func() actor.Actor { return tree.CreateBasicNode(size) })
		pidleft := *context.Spawn(propleft)

		propright := actor.PropsFromProducer(func() actor.Actor { return tree.CreateBasicNode(size) })
		pidright := *context.Spawn(propright)

		fmt.Println("Sending GetBasicNodesMessage to RootNode")
		context.Respond(tree.GetBasicNodesMessage{
			LeftPid:  pidleft,
			RightPid: pidright,
			SSender:  clientpid,
		})
		fmt.Println("Sending new InsertMessage to RootNode")
		time.Sleep(5 * time.Second)
		context.Send(msg.PMessageResult.get, tree.InsertMessage{
			PID:        msg.PMessageResult.PID,
			Element:    msg.PMessageResult.PIDRoot,
			PIDService: msg.PMessageResult.PIDRoot,
			PIDRoot:    msg.PMessageResult.PIDRoot,
		})

	default:
		fmt.Printf("default service")
	}

}

/*
NewServerRemoteActor will return a ServerRemoteActor
*/
func NewServerRemoteActor() actor.Actor {
	log.Message("Hello-Actor is up and running")
	return &ServerRemoteActor{}
}

/*
AddNewTrie will add a rootNode into the map and return the ID and the token for this trie
*/
func AddNewTrie(size int) (ID, Token, actor.PID) {
	id := ID(getRandomID(1024))
	token := Token(fmt.Sprintf("%d", getRandomID(50)))
	props := actor.PropsFromProducer(func() actor.Actor {
		return tree.CreateBasicNode(size)
	})
	pid := context.Spawn(props)
	addInToRootsMap(id, TrieContainer{ //TODO nach struct richten
		Root:    props,
		Token:   token,
		Pid:     pid,
		Size:    size, //TODO refactor
		Context: &context,
	})
	return id, token, *pid
}

/*
MatchIDandToken will check whether a given token and id match.
Will return true in case they do otherwise false.
*/
func MatchIDandToken(id ID, gtoken Token) bool {
	return RootNodes[id].Token == gtoken // Verändert schau commit #71
}

/*
starts an actorsystem that can be reached remotely
*/
func main() {

	var wg sync.WaitGroup

	wg.Add(1)
	defer wg.Wait()

	remote.Start("localhost:8091")

	prop := actor.PropsFromProducer(NewServerRemoteActor)
	globalpid = *context.Spawn(prop)
	// register a name for our local actor so that it can be spawned remotely
	remote.Register("hello", prop)
}
