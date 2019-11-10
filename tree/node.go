package tree

import (
	"fmt"

	"github.com/AsynkronIT/protoactor-go/actor"
)

/*
InsertMessage str
*/
type InsertMessage struct {
	PID     actor.PID
	Element Pair
}

/*
DeleteMessage will
*/
type DeleteMessage struct {
	PID actor.PID
	Key int
}

/*
ChangeValueMessage will
*/
type ChangeValueMessage struct {
	PID     actor.PID
	Element Pair
}

/*
FindMessage will
*/
type FindMessage struct {
	PID actor.PID
	Key int
}

/*
RespMessage will be the response to a given request
*/
type RespMessage struct {
	Ans interface{}
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
	var result interface{}
	switch msg := context.Message().(type) {
	case *InsertMessage:
		if state.IsLeft(msg.Element.Key) {
			result = state.LeftElement.(*Leaf).Insert(msg.Element.Key, msg.Element.Value)
		} else {
			result = state.RightElement.(*Leaf).Insert(msg.Element.Key, msg.Element.Value)
		}
		context.Send(&msg.PID, &RespMessage{
			Ans: result,
		})
	case *DeleteMessage:
		if state.IsLeft(msg.Key) {
			result = state.LeftElement.(*Leaf).Erase(msg.Key)
		} else {
			result = state.RightElement.(*Leaf).Erase(msg.Key)
		}
		context.Send(&msg.PID, &RespMessage{
			Ans: result,
		})
	case *ChangeValueMessage:
		if state.IsLeft(msg.Element.Key) {
			result = state.LeftElement.(*Leaf).Change(msg.Element.Key, msg.Element.Value)
		} else {
			result = state.RightElement.(*Leaf).Change(msg.Element.Key, msg.Element.Value)
		}
		context.Send(&msg.PID, &RespMessage{
			Ans: result,
		})
	case *FindMessage:
		if state.IsLeft(msg.Key) {
			result = state.LeftElement.(*Leaf).Find(msg.Key)
		} else {
			result = state.RightElement.(*Leaf).Find(msg.Key)
		}
		context.Send(&msg.PID, &RespMessage{
			Ans: result,
		})
	}

}

/*
KnownNodeBehavior Method to set the Behavoir of a Node to a Knowing Node.
So it will have to nodes as childs and know's about this childs/manged them.
*/
func (state *Nodeactor) KnownNodeBehavior(context actor.Context) {
	switch msg := context.Message().(type) {
	case *InsertMessage, *DeleteMessage, *FindMessage, *ChangeValueMessage:
		context.Send(state.LeftElement.(*actor.PID), &msg)
	default:
		fmt.Println(msg)
	}
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
getConstrain will return the Storeable of a node
*/
func (state *Nodeactor) getConstrain() int {
	return state.Storable
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
Receive will recieve some messages and direct them to the nodes
*/
func (state *Nodeactor) Receive(context actor.Context) {
	state.Behavior.Receive(context)
}
