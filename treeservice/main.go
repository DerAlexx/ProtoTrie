package main

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/ob-vss-ws19/blatt-3-pwn/messages"
	"github.com/ob-vss-ws19/blatt-3-pwn/tree"
)

/*
ServerRemoteActor w
*/
type ServerRemoteActor struct{}

/*
Receive w
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
NewServerRemoteActor will
*/
func NewServerRemoteActor() actor.Actor {
	return &ServerRemoteActor{}
}

func main() {

	remote.Start("localhost:8091")

	remote.Register("hello", actor.PropsFromProducer(NewServerRemoteActor))

}
