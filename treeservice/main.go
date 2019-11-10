package main

import (
	"sync"

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
func GetPID(id ID) *actor.PID { //TODO hiding
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
func DeleteTrie(id ID, token Token) { // TODO hiding
	if MatchIDandToken(id, token) {
		delete(RootNodes, id)
	}
}

/*
Receive will receive different types of messages from the client. Each message is responsible for one type of action that the tree can execute (e.g. delete a Key, Create a Tree)
After the execution it will return a message to the client
*/
func (*ServerRemoteActor) Receive(context actor.Context) { //TODO
	switch msg := context.Message().(type) {
	case *messages.DeleteRequest:
		nod := tree.GetPID(tree.ID(*messages.DeleteRequest.Id))
		//va, err := tree.Delete(*messages.DeleteRequest.Key, *messages.DeleteRequest.Id, *messages.DeleteRequest.Token)
		context.Respond(&messages.Response{
			if (va && err == nil){
				SomeValue: "Success",
			}
			SomeValue: "Couldnt Delete Entry",
		})
	case *messages.ChangeRequest:
		pa := tree.Pair{
			Key: int(*messages.DeleteRequest.Key),
			Value: int(*messages.DeleteRequest.Value),
		}
		va, err := tree.Change(pa, *messages.DeleteRequest.Id, *messages.DeleteRequest.Token)
		context.Respond(&messages.Response{
			if (va && err == nil){
				SomeValue: "Success",
			}
			SomeValue: "Couldnt Change Entry",
		})
	case *messages.InsertRequest:
		pa := tree.Pair{
			Key: int(*messages.DeleteRequest.Key),
			Value: int(*messages.DeleteRequest.Value),
		}
		va, err := tree.InsertInLeaf(pa, *messages.DeleteRequest.Id, *messages.DeleteRequest.Token)
		context.Respond(&messages.Response{
			if (va && err == nil){
				SomeValue: "Success",
			}
			SomeValue: "Couldnt Insert Entry",
		})
	case *messages.CreateTreeRequest:
		tree.AddNewTrie(int(msg.GetSize_()))
		context.Respond(&messages.Response{
			SomeValue: "Success",
		})
	case *messages.DeleteTreeRequest:
		tree.DeleteTrie(*messages.DeleteRequest.Id, *messages.DeleteRequest.Token)
		context.Respond(&messages.Response{
			SomeValue: "Success",
		})
	}
	case *messages.FindRequest:
		va, err := tree.Find(*messages.DeleteRequest.Key, *messages.FindRequest.Id, *messages.DeleteRequest.Token)
		context.Respond(&messages.Response{
			if (va && err == nil){
				SomeValue: "Success",
			}
			SomeValue: "Couldnt Find Entry",
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
func AddNewTrie(context actor.RootContext, size int) (ID, Token, actor.PID) {
	id := ID(getRandomID())
	token := Token(hashcode.String(string(id)))
	props := actor.PropsFromProducer(func() actor.Actor {
		return CreateBasicNode()
	})
	pid := context.Spawn(props)
	addInToRootsMap(id, TrieContainer{ //TODO nach struct richten
		Root:    props,
		Token:   token,
		Pid:     pid,
		Size:    size,
		Context: context,
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
	context := actor.EmptyRootContext
	var wg sync.WaitGroup

	wg.Add(1)
	defer wg.Wait()

	remote.Start("localhost:8091")

	// register a name for our local actor so that it can be spawned remotely
	remote.Register("hello", actor.PropsFromProducer(NewServerRemoteActor))

}
