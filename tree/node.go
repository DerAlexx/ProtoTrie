package tree

import (
	"fmt"

	"github.com/AsynkronIT/protoactor-go/actor"
)

/*
Pair r
*/
type Pair struct {
	Key   int
	Value string
}

/*
Nodeactor is
*/
type Nodeactor struct {
}

/*
Pingmessage will
*/
type Pingmessage struct {
}

/*
Stacklimittracemessage will
*/
type Stacklimittracemessage struct {
}

/*
Insertmessage str
*/
type Insertmessage struct {
	Element Pair
}

/*
Receive will
*/
func (lactor *Nodeactor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *Pingmessage:
		fmt.Println("Ping")

	case *Stacklimittracemessage:
		fmt.Println("42")

	case *Insertmessage:
		fmt.Println(msg)
	}
}
