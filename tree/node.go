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
	PID        actor.PID
	Element    Pair
	PIDService actor.PID
	PIDRoot    actor.PID
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
GetBasicNodesMessage to get 2 BasicNodes
*/
type GetBasicNodesMessage struct {
	LeftPid  actor.PID
	RightPid actor.PID
	SSender  actor.PID
}

/*
WantBasicNodeActorsMessage is a request to send 2 basic nodes in order to
expand the Trie
*/
type WantBasicNodeActorsMessage struct {
	PMessageResult interface{}
	Size           int
}

/*
Pair is a representation of a Keys Value Combinations
*/
type Pair struct {
	Key   int
	Value string
}

/*
Nodeactor is a actor. Each Node has an actor
*/
type Nodeactor struct {
	Behavior     actor.Behavior
	Storable     int
	LeftElement  interface{}
	RightElement interface{}
	Limit        int
}

/*
CreateBasicNode will return a Basic node containing the parentnode and a left, right leaf.
*/
func CreateBasicNode(limit int) actor.Actor {
	return &Nodeactor{
		//TODO Add a basic Behavior
		Storable:     -1,
		LeftElement:  NewLeaf(),
		RightElement: NewLeaf(),
		Limit:        limit,
	}
}

/*
getLimit will return the max size of a map
*/
func (state *Nodeactor) getLimit() int {
	return state.Limit
}

/*
IsFull will check whether on of the Leafs is full with pairs and in this case
return true else false.
*/
func (state *Nodeactor) IsFull() bool {
	if state.LeftElement.(*Leaf).Size() == state.getLimit() || state.RightElement.(*Leaf).Size() == state.getLimit() {
		return true
	}
	return false
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
		if state.IsFull() {
			context.RequestWithCustomSender(&msg.PIDService, &WantBasicNodeActorsMessage{
				PMessageResult: result,
				Size:           state.getLimit(),
			}, &msg.PIDRoot)
		} else {
			fmt.Println("In the insert")
			context.Send(&msg.PID, &RespMessage{
				Ans: result,
			})
		}
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
	case *GetBasicNodesMessage:
		//works := state.expand(state.LeftElement.(*Leaf), state.LeftElement.(*Leaf), msg)
		if true { // hier eigentlich das ergebnis von works
			context.Send(&msg.SSender, &RespMessage{
				Ans: "Worked",
			})
			state.Behavior.Become(state.KnownNodeBehavior)
		} else {
			context.Send(&msg.SSender, &RespMessage{
				Ans: "Worked not",
			})
		}
	}

}

/*
NewSetBehaviorActor will be the setter for the new behavior.
This will change in case a leaf is full and musst be splitted.
*/
func NewSetBehaviorActor() actor.Actor {
	act := &Nodeactor{
		Behavior: actor.NewBehavior(),
	}
	act.Behavior.Become(act.KnownNodeBehavior)

	return act
}

/*
Will insert into a Nodes Leafs Data
*/
func (state *Nodeactor) insertSplittedMaps(left, right map[int]string) {
	state.LeftElement.(*Leaf).setData(&left)
	state.RightElement.(*Leaf).setData(&right)
}

/*
expand will expand a trie from a given node by adding both sides with a new node
and two leafs in order to get a half-balanced-Trie.
*/
/*
func (state *Nodeactor) expand(left, right *Leaf, msg *GetBasicNodesMessage) bool {
	var (
		leftmap  map[int]string
		rightmap map[int]string
	)
	if left != nil && right != nil && msg != nil {
		leftmap = *left.getData()
		rightmap = *right.getData()

		state.LeftElement = msg.LeftNode   //TODO
		state.RightElement = msg.RightNode //TODO

		state.LeftElement.(*Nodeactor).insertSplittedMaps(sortMap(leftmap))   //TODO
		state.RightElement.(*Nodeactor).insertSplittedMaps(sortMap(rightmap)) //TODO
		return true
	}
	return false
}

*/
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

/*
sortMap sorts a given map, splits it in half and returns 2 maps.
Each created map contains one half of the entries of the given map.
*/
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
