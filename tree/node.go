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
func (state *Nodeactor) IsFull(isleft bool) bool {
	k := state.LeftElement.(*Leaf).Size() == state.getLimit()
	j := state.RightElement.(*Leaf).Size() == state.getLimit()
	fmt.Printf("[+] ISFULL: %t %t %t \n", k, j, isleft)
	if k && isleft {
		return true
	} else if j && !isleft {
		return true
	} else {
		return false
	}
}

/*
unionLeafs will get two leafs and later on union the
data fields of both maps.
*/
func unionLeafs(left, right *Leaf) map[int]string { //????
	leftmap := *left.getData()
	rightmap := right.getData()
	fmt.Println(leftmap)
	fmt.Println(rightmap)
	for k, v := range *rightmap {
		leftmap[k] = v
	}
	fmt.Println(leftmap)
	return leftmap
}

/*
traverseChild will traverse the subtrie by running over the child and return a Map with all pairs.
*/
func (state *Nodeactor) traverseChild() map[int]string {
	switch state.LeftElement.(type) {
	case Leaf:
		return unionLeafs(state.LeftElement.(*Leaf), state.RightElement.(*Leaf))
	case *Leaf:
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
		fmt.Println("[+] Insert into Storing Actor")
		if state.IsFull(state.IsLeft(msg.Element.Key)) {
			fmt.Println("[+] Leafs is full --> Expand")
			context.RequestWithCustomSender(&msg.PIDService, &WantBasicNodeActorsMessage{
				PMessageResult: &msg,
				Size:           state.getLimit(),
			}, context.Self())
		} else {
			if state.IsLeft(msg.Element.Key) {
				fmt.Println("[+] Is Left sorted key")
				result = state.LeftElement.(*Leaf).Insert(msg.Element.Key, msg.Element.Value)
				fmt.Printf("[+] %t \n", state.RightElement.(*Leaf).Contains(msg.Element.Key))
				fmt.Println("[+] Right Leaf MAP")
				fmt.Println(state.RightElement.(*Leaf).getData())
				fmt.Println("[+] Leaf Leaf MAP")
				fmt.Println(state.LeftElement.(*Leaf).getData())
			} else {
				fmt.Println("[+] Is Right sorted key")
				result = state.RightElement.(*Leaf).Insert(msg.Element.Key, msg.Element.Value)
				fmt.Printf("[+] %t \n", state.LeftElement.(*Leaf).Contains(msg.Element.Key))
				fmt.Println("[+] Right Leaf MAP")
				fmt.Println(state.RightElement.(*Leaf).getData())
				fmt.Println("[+] Leaf Leaf MAP")
				fmt.Println(state.LeftElement.(*Leaf).getData())
			}
			state.IsLeft(msg.Element.Key)
			context.Send(&msg.PID, &messages.Response{
				SomeValue: fmt.Sprintf("%t", result.(bool)),
			})
		}
	case DeleteMessage:
		if state.IsLeft(msg.Key) {
			result = state.LeftElement.(*Leaf).Erase(msg.Key)
		} else {
			result = state.RightElement.(*Leaf).Erase(msg.Key)
		}
		context.Send(&msg.PID, &messages.Response{
			SomeValue: fmt.Sprintf("%t", result.(bool)),
		})
	case ChangeValueMessage:
		if state.IsLeft(msg.Element.Key) {
			result = state.LeftElement.(*Leaf).Change(msg.Element.Key, msg.Element.Value)
		} else {
			result = state.RightElement.(*Leaf).Change(msg.Element.Key, msg.Element.Value)
		}
		context.Send(&msg.PID, &messages.Response{
			SomeValue: fmt.Sprintf("%t", result.(bool)),
		})
	case FindMessage:
		if state.IsLeft(msg.Key) {
			result = state.LeftElement.(*Leaf).Find(msg.Key)
			for k := range *state.LeftElement.(*Leaf).getData() {
				print(k)
			}
			for k := range *state.RightElement.(*Leaf).getData() {
				print(k)
			}
		} else {
			result = state.RightElement.(*Leaf).Find(msg.Key)
			for k := range *state.LeftElement.(*Leaf).getData() {
				print(k)
			}
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
		fmt.Print("[+] Before BH Change", &state.Behavior, "\n")
		fmt.Println(msg.LeftPid, msg.RightPid, msg.SSender)
		state.expand(state.LeftElement.(*Leaf), state.RightElement.(*Leaf), &msg, context)
		state.Behavior.Become(state.KnownNodeBehavior)
		fmt.Print("[+] After BH Change", &state.Behavior, "\n")
	case ExpandMessage:
		fmt.Println(msg.NewStorable, msg.LeftMap, msg.RightMap)
		state.SetStoreable(msg.NewStorable)
		state.LeftElement.(*Leaf).insertMap(msg.LeftMap)
		state.RightElement.(*Leaf).insertMap(msg.RightMap)
		fmt.Println(state.LeftElement.(*Leaf), state.RightElement.(*Leaf))
	case TraverseMessage:
		ret := state.traverseChild()
		if ret != nil && len(ret) > 0 {
			fmt.Println("[+] Currently traversing")
			map32 := make(map[int32]string, len(ret))
			for k, v := range ret {
				map32[int32(k)] = v
			}
			context.Send(&msg.PID, &messages.TraverseResponse{
				Arr: map32,
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
				biggestKey = key
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
	fmt.Println(*left, *right, *msg)
	if left != nil && right != nil && msg != nil {
		leftmap = *left.getData()
		rightmap = *right.getData()

		state.LeftElement = msg.LeftPid
		state.RightElement = msg.RightPid

		ll, lr, storeL := sortMap(leftmap)
		rl, rr, storeR := sortMap(rightmap)

		fmt.Println("[+] Expand Start")
		fmt.Println(len(ll), len(lr), len(rl), len(rr))
		for k, v := range ll {
			fmt.Println(k, v)
		}
		for k, v := range lr {
			fmt.Println(k, v)
		}

		for k, v := range rl {
			fmt.Println(k, v)
		}

		for k, v := range rr {
			fmt.Println(k, v)
		}
		fmt.Printf("[+] Knoten Storables L: %d, R: %d", storeL, storeR)
		fmt.Println("[+] Expand End")
		context.Send(msg.LeftPid, ExpandMessage{
			NewStorable: storeL,
			LeftMap:     ll,
			RightMap:    lr,
		})

		context.Send(msg.RightPid, ExpandMessage{
			NewStorable: storeR,
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
		fmt.Println("Knowing Insert")
		if state.IsLeft(msg.Element.Key) {
			fmt.Println("Knowing Insert Left")
			fmt.Println(msg.Element, msg.PID, msg.PIDRoot, msg.PIDService)
			context.Send(state.LeftElement.(*actor.PID), msg)
		} else {
			fmt.Println("Knowing Insert Right")
			fmt.Println(msg)
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
	fmt.Printf("[+] is: %d, compare: %d, has: %t \n", is, value, has)
	if !has {
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
returnAllKeys //TODO
*/
func returnAllKey(m map[int]string) []int {
	arr := []int{}
	for k := range m {
		arr = append(arr, k)
	}
	return arr
}

/*
sortMap sorts a given map, splits it in half and returns 2 maps.
Each created map contains one half of the entries of the given map.
*/
func sortMap(m map[int]string) (map[int]string, map[int]string, int) {
	var keys []int = returnAllKey(m)
	sort.Ints(keys)

	left := make(map[int]string, len(m)/2+1)
	right := make(map[int]string, len(m)/2+1)

	if len(keys) <= 0 {
		return map[int]string{}, map[int]string{}, -1
	} else if len(keys)%2 != 0 {
		for i := 0; i < int(len(m)/2)+1; i++ {
			left[keys[i]] = m[keys[i]]
			if i+int(len(m)/2)+1 < len(keys) {
				right[keys[i+int(len(m)/2)+1]] = m[i]
			}
		}
	} else {
		for i := 0; i < int(len(m)/2); i++ {
			left[keys[i]] = m[keys[i]]
			right[keys[i+int(len(m)/2)]] = m[keys[i+int(len(m)/2)]]
		}
	}
	return left, right, findBiggestKey(left)

}
