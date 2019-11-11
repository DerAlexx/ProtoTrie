package tree

import (
	"fmt"
	"sort"

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
		//TODO Add a basic Behavior
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
	case *DeleteMessage:
		if state.IsLeft(msg.Key) {
			context.Send(state.LeftElement.(*actor.PID), &msg)
		} else {
			context.Send(state.RightElement.(*actor.PID), &msg)
		}
	case *FindMessage:
		if state.IsLeft(msg.Key) {
			context.Send(state.LeftElement.(*actor.PID), &msg)
		} else {
			context.Send(state.RightElement.(*actor.PID), &msg)
		}
	case *InsertMessage:
		if state.IsLeft(msg.Element.Key) {
			context.Send(state.LeftElement.(*actor.PID), &msg)
		} else {
			context.Send(state.RightElement.(*actor.PID), &msg)
		}
	case *ChangeValueMessage:
		if state.IsLeft(msg.Element.Key) {
			context.Send(state.LeftElement.(*actor.PID), &msg)
		} else {
			context.Send(state.RightElement.(*actor.PID), &msg)
		}
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

//TODO Comment
func (state *Nodeactor) splitNode() []map[int]string {
	leftmapone, leftmaptwo := sortMap(*state.LeftElement.(*Leaf).getData())
	rightmapone, rightmaptwo := sortMap(*state.RightElement.(*Leaf).getData())
	return []map[int]string{leftmapone, leftmaptwo, rightmapone, rightmaptwo}
}

//TODO Comment
func sortMap(m map[int]string) (r1 map[int]string, r2 map[int]string) {
	keys := make([]int, 0, len(m))
	pairs := make([]Pair, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for k := range keys {
		pairs = append(pairs, Pair{
			Key:   k,
			Value: m[k],
		})
	}
	mapsizeleft := len(pairs) / 2

	for i := 0; i < mapsizeleft; i++ {
		r1[pairs[i].Key] = pairs[i].Value
	}

	for ih := mapsizeleft; ih < len(pairs); ih++ {
		r2[pairs[ih].Key] = pairs[ih].Value
	}

	return r1, r2
}

// TODO Insert --> und daraufhin wechsel das verhaltens und Erweiterung
// des baums und einpflege der
// Daten in die neuen Leafs.

// TODO Split
