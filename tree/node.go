package tree

import (
	"fmt"
	"reflect"

	"github.com/AsynkronIT/protoactor-go/actor"
)

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
	behavior     actor.behavior
	LeftElement  interface{}
	RightElement interface{}
}

/*
CreateBasicNode will return a Basic node containing the parentnode and a left, right leaf.
*/
func CreateBasicNode() *Nodeactor {
	return Nodeactor{
		LeftElement:  NewLeaf(),
		RightElement: NewLeaf(),
	}
}

/*
StoringNodeBehavior Method to set the Behavior of a Node to a Storing Node.
So it will have to leafs as childs and store information in this leafs.
*/
func (state *Nodeactor) StoringNodeBehavior(context actor.Context) {

}

/*
KnownNodeBehavior Method to set the Behavoir of a Node to a Knowing Node.
So it will have to nodes as childs and know's about this childs/manged them.
*/
func (state *Nodeactor) KnownNodeBehavior(context actor.Context) {

}

/*
Receive will
*/
func (state *Nodeactor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *Pingmessage:
		fmt.Println("Ping")

	case *Stacklimittracemessage:
		fmt.Println("42")

	case *Insertmessage:
		fmt.Println(msg)
		n := Nodeactor{
			LeftElement:  12,
			RightElement: "Hello",
		}
		fmt.Println(reflect.TypeOf(n.LeftElement), reflect.TypeOf(n.RightElement))
		n.LeftElement = "3"
		fmt.Println(reflect.TypeOf(n.LeftElement), reflect.TypeOf(n.RightElement))
	}
}
