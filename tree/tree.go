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
	Root  *actor.Props
	Token Token
	Pid   *actor.PID
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
getPID will return the PID of the rootnode of a given Trie-ID
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
AddNewTrie will add a rootNode into the map and return the ID and the token for this trie
*/
func AddNewTrie(context *actor.RootContext) (ID, Token, actor.PID) {
	id := ID(getRandomID())
	token := Token(hashcode.String(string(id)))
	props := actor.PropsFromProducer(func() actor.Actor {
		return tree.CreateBasicNode()
	})
	pid := context.Spawn(props)
	addInToRootsMap(id, TrieContainer{
		Root:  props,
		Token: token,
		Pid:   pid,
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
deleteTrie will delete a Trie
*/
func DeleteTrie(id ID, token Token) {
	if matchIDandToken(id, token) {
		delete(RootNodes, id)
	}
}
