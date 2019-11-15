package tree

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/ob-vss-ws19/blatt-3-pwn/messages"
)

/*
TraverseMessage will be message to traverse the tree.
*/
type TraverseMessage struct {
	PID actor.PID
}

/*
InsertMessage will be a insertmessage with all information to insert something
in the Trie
*/
type InsertMessage struct {
	PID        actor.PID
	Element    Pair
	PIDService actor.PID
	PIDRoot    actor.PID
}

/*
DeleteMessage will be a message in order to delete a key from the given trie
*/
type DeleteMessage struct {
	PID actor.PID
	Key int
}

/*
ChangeValueMessage will be a message in order to change a pair in the map (just the value)
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
	LeftPid  *actor.PID
	RightPid *actor.PID
	SSender  *actor.PID
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
ExpandMessage will be a message containing a MAP
*/
type ExpandMessage struct {
	NewStorable int
	LeftMap     map[int]string
	RightMap    map[int]string
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
	act := &Nodeactor{
		Behavior:     actor.NewBehavior(),
		Storable:     -1,
		LeftElement:  NewLeaf(),
		RightElement: NewLeaf(),
		Limit:        limit,
	}
	act.Behavior.Become(act.StoringNodeBehavior)
	return act
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
unionLeafs will get two leafs and later on union the
data fields of both maps.
*/
func unionLeafs(left, right *Leaf) map[int]string {
	leftmap := *left.getData()
	rightmap := right.getData()
	for k, v := range *rightmap {
		leftmap[k] = v
	}
	return leftmap
}

/*
traverseChild will traverse the subtrie by running over the child and return a Map with all pairs.
*/
func (state *Nodeactor) traverseChild() map[int]string {
	switch state.LeftElement.(type) {
	case Leaf:
		return unionLeafs(state.LeftElement.(*Leaf), state.RightElement.(*Leaf))
	default:
		return nil
	}
}

/*
StoringNodeBehavior Method to set the Behavior of a Node to a Storing Node.
So it will have to leafs as childs and store information in this leafs.
*/
func (state *Nodeactor) StoringNodeBehavior(context actor.Context) {
	var result interface{}
	switch msg := context.Message().(type) {
	case InsertMessage:
		fmt.Println("Insert Service")
		if state.IsFull() {
			fmt.Println("Insert isFull Service")
			context.RequestWithCustomSender(&msg.PIDService, &WantBasicNodeActorsMessage{
				PMessageResult: &msg,
				Size:           state.getLimit(),
			}, context.Self())
		} else {
			if state.IsLeft(msg.Element.Key) {
				fmt.Println("Insert isLeft Service")
				result = state.LeftElement.(*Leaf).Insert(msg.Element.Key, msg.Element.Value)
				fmt.Printf("[+] %t", state.RightElement.(*Leaf).Contains(msg.Element.Key))
			} else {
				fmt.Println("Insert isRight Service")
				result = state.RightElement.(*Leaf).Insert(msg.Element.Key, msg.Element.Value)
				fmt.Printf("[+] %t", state.LeftElement.(*Leaf).Contains(msg.Element.Key))
			}
			fmt.Println("Return insert Service")
			context.Send(&msg.PID, &messages.Response{
				SomeValue: fmt.Sprintf("%t", result.(bool)),
			})
		}
	case DeleteMessage:
		fmt.Printf("[+] Given Key %d \n", msg.Key)
		if state.IsLeft(msg.Key) {
			//fmt.Printf("[+] %d \n", state.Limit)
			//fmt.Printf("[+] Left %t \n", state.LeftElement)
			result = state.LeftElement.(*Leaf).Erase(msg.Key)
			//fmt.Printf("[+] %d \n", state.Limit)
		} else {
			//fmt.Printf("[+] %d \n", state.Limit)
			//fmt.Printf("[+] Right %t \n", state.LeftElement)
			result = state.RightElement.(*Leaf).Erase(msg.Key)
			//fmt.Printf("[+] %d \n", state.Limit)
		}
		fmt.Printf("[+] Contains Left: %t Contains Right: %t ", state.LeftElement.(*Leaf).Contains(msg.Key), state.RightElement.(*Leaf).Contains(msg.Key))
		context.Send(&msg.PID, &messages.Response{
			SomeValue: fmt.Sprintf("%t", result.(bool)),
		})
	case ChangeValueMessage:
		if state.IsLeft(msg.Element.Key) {
			fmt.Printf("[+] Before %s \n", state.LeftElement.(*Leaf).Find(msg.Element.Key))
			result = state.LeftElement.(*Leaf).Change(msg.Element.Key, msg.Element.Value)
			fmt.Printf("[+] After %s \n", state.LeftElement.(*Leaf).Find(msg.Element.Key))
		} else {
			result = state.RightElement.(*Leaf).Change(msg.Element.Key, msg.Element.Value)
		}
		context.Send(&msg.PID, &messages.Response{
			SomeValue: fmt.Sprintf("%t", result.(bool)),
		})
	case FindMessage:
		if state.IsLeft(msg.Key) {
			result = state.LeftElement.(*Leaf).Find(msg.Key)
			fmt.Println("left")
			for k := range *state.LeftElement.(*Leaf).getData() {
				print(k)
			}
			fmt.Println("\nright")
			for k := range *state.RightElement.(*Leaf).getData() {
				print(k)
			}
		} else {
			result = state.RightElement.(*Leaf).Find(msg.Key)
			fmt.Println("left")
			for k := range *state.LeftElement.(*Leaf).getData() {
				print(k)
			}
			fmt.Println("\nright")
			for k := range *state.RightElement.(*Leaf).getData() {
				print(k)
			}
		}
		if result == "" {
			result = fmt.Sprintf("There is no Entry in the Database with id %d", msg.Key)
		}
		context.Send(&msg.PID, &messages.Response{
			SomeValue: fmt.Sprintf("%s", result.(string)),
		})
	case GetBasicNodesMessage:
		state.expand(state.LeftElement.(*Leaf), state.LeftElement.(*Leaf), &msg, context)
	case ExpandMessage:
		state.SetStoreable(msg.NewStorable)
		state.LeftElement.(*Leaf).insertMap(msg.LeftMap)
		state.RightElement.(*Leaf).insertMap(msg.RightMap)
	case TraverseMessage:
		ret := state.traverseChild() //TODO was wird hier zurück gesendet vlt ein Array aus Array's ?
		if ret != nil {
			context.Send(&msg.PID, messages.TraverseResponse{
				Arr: ret,
			})
		}
	default:
		fmt.Println(reflect.TypeOf(msg))
	}

}

/*
findBiggestKey will return the biggest key in a map
*/
func findBiggestKey(m map[int]string) int {
	if len(m) > 0 {
		biggestKey := 0
		for key := range m {
			if key > biggestKey {
				key = biggestKey
			}
		}
		return biggestKey
	}
	return 0
}

/*
expand will expand a trie from a given node by adding both sides with a new node
and two leafs in order to get a half-balanced-Trie.
*/
func (state *Nodeactor) expand(left, right *Leaf, msg *GetBasicNodesMessage, context actor.Context) bool {
	var (
		leftmap  map[int]string
		rightmap map[int]string
	)
	if left != nil && right != nil && msg != nil {
		leftmap = *left.getData()
		rightmap = *right.getData()

		state.LeftElement = msg.LeftPid
		state.RightElement = msg.RightPid

		ll, lr := sortMap(leftmap)
		rl, rr := sortMap(rightmap)

		context.Send(msg.LeftPid, ExpandMessage{
			NewStorable: findBiggestKey(ll),
			LeftMap:     ll,
			RightMap:    lr,
		})

		context.Send(msg.RightPid, ExpandMessage{
			NewStorable: findBiggestKey(rl),
			LeftMap:     rl,
			RightMap:    rr,
		})

		return true
	}
	return false
}

/*
KnownNodeBehavior Method to set the Behavoir of a Node to a Knowing Node.
So it will have to nodes as childs and know's about this childs/manged them.
*/
func (state *Nodeactor) KnownNodeBehavior(context actor.Context) {
	switch msg := context.Message().(type) {
	case DeleteMessage:
		if state.IsLeft(msg.Key) {
			context.Send(state.LeftElement.(*actor.PID), msg)
		} else {
			context.Send(state.RightElement.(*actor.PID), msg)
		}
	case FindMessage:
		if state.IsLeft(msg.Key) {
			context.Send(state.LeftElement.(*actor.PID), msg)
		} else {
			context.Send(state.RightElement.(*actor.PID), msg)
		}
	case InsertMessage:
		if state.IsLeft(msg.Element.Key) {
			context.Send(state.LeftElement.(*actor.PID), msg)
		} else {
			context.Send(state.RightElement.(*actor.PID), msg)
		}
	case ChangeValueMessage:
		if state.IsLeft(msg.Element.Key) {
			context.Send(state.LeftElement.(*actor.PID), msg)
		} else {
			context.Send(state.RightElement.(*actor.PID), msg)
		}
	case TraverseMessage:
		context.Send(state.LeftElement.(*actor.PID), msg)
		context.Send(state.RightElement.(*actor.PID), msg)
	default:
		fmt.Println(reflect.TypeOf(msg))
	}
}

/*
HasValueToDecide will check whether a given node has a value set to decide whether the value belongs to the
left side or not. In case it does it will return true otherwise false.
*/
func (state *Nodeactor) HasValueToDecide() (bool, int) {
	if state.Storable > 0 {
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
	fmt.Printf("is: %d, compare: %d, has: %t \n", is, value, has)
	if !has {
		fmt.Println("Ruf immer wieder den Setter auf")
		state.SetStoreable(value)
		return true
	}
	if is >= value {
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
func sortMap(m map[int]string) (map[int]string, map[int]string) {
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

	var (
		r1, r2 map[int]string
	)

	for i := 0; i < mapsizeleft; i++ {
		r1[pairs[i].Key] = pairs[i].Value
	}

	for ih := mapsizeleft; ih < len(pairs); ih++ {
		r2[pairs[ih].Key] = pairs[ih].Value
	}

	return r1, r2
}
