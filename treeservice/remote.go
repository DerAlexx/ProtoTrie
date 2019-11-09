package treeservice

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/ob-vss-ws19/blatt-3-pwn/messages"
)

type ServerRemoteActor struct{}

func (*ServerRemoteActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *messages.DeleteRequest:
		//delete(*messages.DeleteRequest.Token, *messages.DeleteRequest.Id, *messages.DeleteRequest.Key)
		context.Send(msg.Sender, &messages.Response{
			SomeValue: "result",
		})
	case *messages.ChangeRequest:
		//change(*messages.DeleteRequest.Token, *messages.DeleteRequest.Id, *messages.DeleteRequest.P)
		context.Send(msg.Sender, &messages.Response{
			SomeValue: "result",
		})
	case *messages.InsertRequest:
		//insert(*messages.DeleteRequest.Token, *messages.DeleteRequest.Id, *messages.DeleteRequest.P)
		context.Send(msg.Sender, &messages.Response{
			SomeValue: "result",
		})
	case *messages.CreateTreeRequest:
		//createTree(*messages.DeleteRequest.Size)
		context.Send(msg.Sender, &messages.Response{
			SomeValue: "result",
		})
	case *messages.DeleteTreeRequest:
		deleteTrie(*messages.DeleteRequest.Id, *messages.DeleteRequest.Token)
		context.Send(msg.Sender, &messages.Response{
			SomeValue: "result",
		})
	}
}

func main() {
	remote.Start("localhost:8091")

	// register a name for our local actor so that it can be spawned remotely
	remote.Register("hello", actor.PropsFromProducer(NewServerRemoteActor))

}

func NewServerRemoteActor() actor.Actor {
	return &ServerRemoteActor{}
}
