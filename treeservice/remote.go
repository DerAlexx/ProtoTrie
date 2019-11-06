package treeservice

import (
	"fmt"
	"reflect"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
)

type ServerRemoteActor struct{}

func (*ServerRemoteActor) Receive(context actor.Context) {
    switch msg := context.Message().(type) {
    case *messages.Request:
        context.Send(msg.Sender, &messages.Response{
            SomeValue: "result",
        })
    }
}

func main() {
    remote.Start("localhost:8091")

    // register a name for our local actor so that it can be spawned remotely
    remote.Register("hello", actor.PropsFromProducer(func() actor.Actor { return &ServerRemoteActor{} }))
	console.ReadLine()
	//Methode decide(message)
}

//Methode decide(message) switch case, ruft methode auf und gibt je nach message inhalt response zurueck