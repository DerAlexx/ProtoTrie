package tree

import (
	"fmt"

	"github.com/AsynkronIT/protoactor-go/actor"
)

/*
Insertmessage str
*/
type Insertmessage struct {
	Element Pair
}

/*
DeleteMessage will
*/
type DeleteMessage struct {
	Key int
}

/*
ChangeValueMessage will
*/
type ChangeValueMessage struct {
	Element Pair
}

/*
Pair is a representation of a Keys Value Combinations
*/
type Pair struct {
	Key   int
	Value string
}

/*
Nodeactor is
*/
type Nodeactor struct {
	Behavior     actor.Behavior
	Storable     int64
	LeftElement  interface{}
	RightElement interface{}
}

/*
CreateBasicNode will return a Basic node containing the parentnode and a left, right leaf.
*/
func CreateBasicNode() *Nodeactor {
	return &Nodeactor{
		LeftElement:  NewLeaf(),
		RightElement: NewLeaf(),
	}
}

/*
Insert will insert a pair(key/value) into the correct leafs.
*/
func (state *Nodeactor) Insert(pair Pair) (bool, error) {
	return false, nil
}

/*
Delete a entry from the Trie.
*/
func (state *Nodeactor) Delete(key int) (bool, error) {
	return false, nil
}

/*
ChangeValueMessage a entry from the Trie.
*/
func (state *Nodeactor) ChangeValueMessage(pair Pair) (bool, error) {
	return false, nil
}

/*
StoringNodeBehavior Method to set the Behavior of a Node to a Storing Node.
So it will have to leafs as childs and store information in this leafs.
*/
func (state *Nodeactor) StoringNodeBehavior(context actor.Context) {
	switch msg := context.Message().(type) {
	case *Insertmessage:
		state.Insert(msg.Element)
	case *DeleteMessage:
		state.Delete(msg.Key)
	case *ChangeValueMessage:
		state.ChangeValueMessage(msg.Element)
	}
}

/*
KnownNodeBehavior Method to set the Behavoir of a Node to a Knowing Node.
So it will have to nodes as childs and know's about this childs/manged them.
*/
func (state *Nodeactor) KnownNodeBehavior(context actor.Context) {
	switch msg := context.Message().(type) {
	case *Insertmessage:
		fmt.Println(msg)
	}
}

/*
Receive will recieve some messages and direct them to the nodes
*/
func (state *Nodeactor) Receive(context actor.Context) {

}
