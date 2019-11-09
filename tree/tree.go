package tree

import (
	"errors"
	"math/rand"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/hashicorp/terraform/helper/hashcode"
)

/*
ID is the ID of a new created Trie
*/
type ID int

/*
Token is
*/
type Token string

/*
TrieContainer is a container for a Trie containing a Rootnode a Token and the PID of the Root
*/
type TrieContainer struct {
	Root    *actor.Props
	Token   Token
	Pid     *actor.PID
	Size    int
	Context *actor.RootContext
}

/*
RootNodes is a Map of rootnodes containing all roots --> Sorted by Key (Trie ID) and Value (Root Node --> actor.Actor)
*/
var RootNodes map[ID]TrieContainer = make(map[ID]TrieContainer)

/*
contains will check whether a key is in the map or not
*/
func contains(preID int) bool {
	_, contains := RootNodes[ID(preID)]
	return contains
}

/*
getRandomID will return a random id which is not allready in the map
*/
func getRandomID() int {
	for {
		potantialID := rand.Intn(5815831813581)
		if contains(potantialID) {
			return potantialID
		}
	}
}

/*
GetPID will return the PID of the rootnode of a given Trie-ID
*/
func GetPID(id ID) *actor.PID {
	return RootNodes[id].Pid
}

/*
addInToRootsMap will add a trie into the map of roots with its id as a parameter
*/
func addInToRootsMap(id ID, trie TrieContainer) {
	RootNodes[id] = trie
}

/*
iterator will iterate over the tree and will return the node with the the correct leaf
*/
func (n *Nodeactor) iteraTor(key int) *Nodeactor {
	switch t := n.LeftElement.(type) {
	case *Leaf:
		return n
	case *Nodeactor:
		if n.getConstrain() >= key { // Nach Links gehen
			return t.iteraTor(key)
		}
		switch t := n.RightElement.(type) { //
		case *Nodeactor:
			return t.iteraTor(key)
		}
	}
	return n
}

/*
InsertInLeaf will insert a value into a given trie incase the
*/
func (n *Nodeactor) InsertInLeaf(pair Pair, id ID, token Token) (bool, error) {
	if MatchIDandToken(id, token) {
		targetMasterNode := n.iteraTor(pair.Key)
		if targetMasterNode.IsLeft(pair.Key) {
			targetMasterNode.LeftElement.(*Leaf).Insert(pair.Key, pair.Value)
		}
		targetMasterNode.RightElement.(*Leaf).Insert(pair.Key, pair.Value)
	}
	return false, errors.New("The given ID and Token doesnt match up")
}

/*
AddNewTrie will add a rootNode into the map and return the ID and the token for this trie
*/
func AddNewTrie(size int) (ID, Token, actor.PID) {
	context := actor.EmptyRootContext
	id := ID(getRandomID())
	token := Token(hashcode.String(string(id)))
	props := actor.PropsFromProducer(func() actor.Actor {
		return CreateBasicNode()
	})
	pid := context.Spawn(props)
	addInToRootsMap(id, TrieContainer{
		Root:    props,
		Token:   token,
		Pid:     pid,
		Size:    size,
		Context: context,
	})
	return id, token, *pid
}

/*
MatchIDandToken will check whether a given token and id match.
Will return true in case they do otherwise false.
*/
func MatchIDandToken(id ID, gtoken Token) bool {
	return RootNodes[id].Token == gtoken
}

/*
DeleteTrie will delete a Trie
*/
func DeleteTrie(id ID, token Token) {
	if MatchIDandToken(id, token) {
		delete(RootNodes, id)
	}
}
