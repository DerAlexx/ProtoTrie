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
GetPID will return the PID of the rootnode of a given Trie-ID
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

		id := *messages.DeleteRequest.Id
		rootpid := GetPID(id)
		tok := *messages.DeleteRequest.Token

		if MatchIDandToken(id, tok) {
			context.Send(rootpid, tree.DeleteMessage{
				PID: globalpid,
				Key: *messages.DeleteRequest.Key,
			})
		} else {
			context.Respond(&messages.Response{
				SomeValue: "Wrong Combination of ID and Token",
			})
		}
	case *messages.ChangeRequest:
		pa := Pair{
			Key:   int(*messages.ChangeRequest.Key),
			Value: int(*messages.ChangeRequest.Value),
		}
		id := *messages.ChangeRequest.Id
		rootpid := GetPID(id)
		tok := *messages.ChangeRequest.Token

		if MatchIDandToken(id, tok) {
			context.Send(rootpid, tree.ChangeValueMessage{
				PID:     globalpid,
				Element: pa,
			})
		} else {
			context.Respond(&messages.Response{
				SomeValue: "Wrong Combination of ID and Token",
			})
		}
	case *messages.InsertRequest:
		pa := Pair{
			Key:   int(*messages.InsertRequest.Key),
			Value: int(*messages.InsertRequest.Value),
		}
		id := *messages.InsertRequest.Id
		rootpid := GetPID(id)
		tok := *messages.InsertRequest.Token

		if MatchIDandToken(id, tok) {
			context.Send(rootpid, tree.InsertMessage{
				PID:     globalpid,
				Element: pa,
			})
		} else {
			context.Respond(&messages.Response{
				SomeValue: "Wrong Combination of ID and Token",
			})
		}
	case *messages.FindRequest:
		id := *messages.FindRequest.Id
		rootpid := getPID(id)
		tok := *messages.FindRequest.Token

		if MatchIDandToken(id, tok) {
			context.Send(rootpid, tree.FindMessage{
				PID: globalpid,
				Key: *messages.FindRequest.Key,
			})
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
		if true {
			context.Respond(&messages.Response{
				SomeValue: "Success",
			})
		} else {
			context.Respond(&messages.Response{
				SomeValue: "Couldnt delete the tree",
			})
		}
	case *tree.RespMessage:
		switch val := msg.Ans.(type) {
		case bool:
			if val {
				context.Respond(&messages.Response{
					SomeValue: "Success",
				})

			} else {
				context.Respond(&messages.Response{
					SomeValue: "Operation couldnt be executed",
				})
			}
		case string:
			context.Respond(&messages.Response{
				SomeValue: msg.Ans.(string),
			})
		}
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
		return tree.CreateBasicNode()
	})
	pid := context.Spawn(props)
	addInToRootsMap(id, TrieContainer{ //TODO nach struct richten
		Root:    props,
		Token:   token,
		Pid:     pid,
		Size:    size,
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

var globalpid actor.PID
var context actor.RootContext

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
