package main

import (
	"fmt"
	"math/rand"
	"sync"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/hashicorp/terraform/helper/hashcode"
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

/*
getRandomID will return a random id which is not allready in the map
*/
func getRandomID() int {
	for {
		potantialID := rand.Intn(5815831813581)
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

/*
DeleteTrie will delete a Trie
*/
func deleteTrie(id ID, token Token) bool {
	if MatchIDandToken(id, token) {
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
			context.RequestWithCustomSender(rootpid, tree.DeleteMessage{
				PID: msg.Sender(),
				Key: int(msg.GetKey()),
			}, globalpid)
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
			context.RequestWithCustomSender(rootpid, tree.ChangeValueMessage{
				PID:     msg.Sender(),
				Element: pa,
			}, globalpid)
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
		clientpid := msg.Sender()

		if MatchIDandToken(id, tok) {
			context.RequestWithCustomSender(rootpid, tree.InsertMessage{
				PID:     msg.Sender(),
				Element: pa,
			}, globalpid)
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
			context.RequestWithCustomSender(rootpid, tree.FindMessage{
				PID: msg.Sender(),
				Key: int(msg.GetKey()),
			}, globalpid)
		} else {
			context.Respond(&messages.Response{
				SomeValue: "Wrong Combination of ID and Token",
			})
		}
	case *messages.CreateTreeRequest:
		i, t, _ := AddNewTrie(int(msg.GetSize_()))
		context.Respond(&messages.Response{
			SomeValue: fmt.Sprintf("Your ID: %d, Your Token: %s", int(i), string(t)),
		})
	case *messages.DeleteTreeRequest:
		ret := deleteTrie(ID(msg.GetId()), Token(msg.GetToken()))
		if ret {
			context.Respond(&messages.Response{
				SomeValue: "Success",
			})
		} else {
			context.Respond(&messages.Response{
				SomeValue: "Couldnt delete the tree",
			})
		}
	case *tree.WantBasicNodeActorsMessage:
		size := msg.Size
		id := ID(msg.GetId())
		rootpid := getPID(id)

		context.Send(rootpid,tree.GetBasicNodesMessage{
			*Nodeactor left := tree.CreateBasicNode(size),
			propleft := actor.PropsFromProducer(left)
			pidleft = *context.Spawn(propleft)
			

			*Nodeactor right := tree.CreateBasicNode(size),
			propright := actor.PropsFromProducer(right)
			pidright = *context.Spawn(propright)

			LeftNode: left,
			LeftPID: pidleft
			RightNode: right,
			RightPID: pidright,
			Sender: clientpid
		})
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
	id := ID(getRandomID())
	token := Token(hashcode.String(string(id)))
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
	return RootNodes[id].Token == gtoken
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