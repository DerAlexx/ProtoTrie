package main

import (
	"sync"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/ob-vss-ws19/blatt-3-pwn/messages"
	"github.com/ob-vss-ws19/blatt-3-pwn/tree"
)

/*
ServerRemoteActor represents a RemoteActor in the service
*/
type ServerRemoteActor struct{}

/*
Receive will receive different types of messages from the client. Each message is responsible for one type of action that the tree can execute (e.g. delete a Key, Create a Tree)
After the execution it will return a message to the client
*/
func (*ServerRemoteActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *messages.DeleteRequest:
		//delete(*messages.DeleteRequest.Token, *messages.DeleteRequest.Id, *messages.DeleteRequest.Key)
		context.Respond(&messages.Response{
			SomeValue: "result",
		})
	case *messages.ChangeRequest:
		//change(*messages.DeleteRequest.Token, *messages.DeleteRequest.Id, *messages.DeleteRequest.P)
		context.Respond(&messages.Response{
			SomeValue: "result",
		})
	case *messages.InsertRequest:
		//insert(*messages.DeleteRequest.Token, *messages.DeleteRequest.Id, *messages.DeleteRequest.P)
		context.Respond(&messages.Response{
			SomeValue: "result",
		})
	case *messages.CreateTreeRequest:
		tree.AddNewTrie(int(msg.GetSize_()))
		// int(*(*(messages.CreateTreeRequest).Size))
		context.Respond(&messages.Response{
			SomeValue: "result",
		})
	case *messages.DeleteTreeRequest:
		//deleteTrie(*messages.DeleteRequest.Id, *messages.DeleteRequest.Token)
		context.Respond(&messages.Response{
			SomeValue: "result",
		})
	}
}

/*
NewServerRemoteActor will return a ServerRemoteActor
*/
func NewServerRemoteActor() actor.Actor {
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
