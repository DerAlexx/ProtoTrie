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

context := actor.EmptyRootContext

/*
Receive will receive different types of messages from the client. Each message is responsible for one type of action that the tree can execute (e.g. delete a Key, Create a Tree)
After the execution it will return a message to the client
*/
func (*ServerRemoteActor) Receive(context actor.Context) {
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
starts an actorsystem that can be reached remotely
*/
func main() {
	var wg sync.WaitGroup

	wg.Add(1)
	defer wg.Wait()

	remote.Start("localhost:8091")

	// register a name for our local actor so that it can be spawned remotely
	remote.Register("hello", actor.PropsFromProducer(NewServerRemoteActor))

}
