package tree

import (
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
	Storable     int
	LeftElement  interface{}
	RightElement interface{}
}

/*
CreateBasicNode will return a Basic node containing the parentnode and a left, right leaf.
*/
func CreateBasicNode() *Nodeactor {
	return &Nodeactor{
		Storable:     -1,
		LeftElement:  NewLeaf(),
		RightElement: NewLeaf(),
	}
}

/*
StoringNodeBehavior Method to set the Behavior of a Node to a Storing Node.
So it will have to leafs as childs and store information in this leafs.
*/
func (state *Nodeactor) StoringNodeBehavior(context actor.Context) {
	/*
		switch msg := context.Message().(type) {
			case *Insertmessage:

			case *DeleteMessage:

			case *ChangeValueMessage:

			}
	*/
}

/*
KnownNodeBehavior Method to set the Behavoir of a Node to a Knowing Node.
So it will have to nodes as childs and know's about this childs/manged them.
*/
func (state *Nodeactor) KnownNodeBehavior(context actor.Context) {
	/*
		switch msg := context.Message().(type) {
			case *Insertmessage:
				fmt.Println(msg)
			}
	*/
}

/*
HasValueToDecide will check whether a given node has a value set to decide whether the value belongs to the
left side or not. In case it does it will return true otherwise false.
*/
func (state *Nodeactor) HasValueToDecide() (bool, int) {
	if state.Storable != -1 {
		return true, state.Storable
	}
	return false, -1
}

/*
SetStoreable will set the value so it can be use to decide the target leaf
*/
func (state *Nodeactor) SetStoreable(value int) {
	state.Storable = value
}

/*
IsLeft will check whether a given Key belongs to the left leaf incase the it does it will return true
otherwise it will return false
*/
func (state *Nodeactor) IsLeft(value int) bool {
	has, is := state.HasValueToDecide()
	if !has {
		state.SetStoreable(value)
		return true
	}
	if is <= value {
		return true
	}
	return false
}

/*
getConstrain will return the Storeable of a node
*/
func (state *Nodeactor) getConstrain() int {
	return state.Storable
}

/*
Receive will recieve some messages and direct them to the nodes
*/
func (state *Nodeactor) Receive(context actor.Context) {

}
