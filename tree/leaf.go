package tree

/*
Leaf is the representation of the LEAF of a Node at
the end of a treebranch
*/
type Leaf struct {
	Data map[int]string
}

/*
NewLeaf will return a new Leaf
*/
func NewLeaf() Leaf {
	return Leaf{
		Data: make(map[int]string),
	}
}

/*
Contains will check whether or not a key is already in the map
@return will return true in case the key is in the map otherwise false
*/
func (l *Leaf) Contains(key int) bool {
	_, ok := l.Data[key]
	return ok
}

/*
Insert will insert a pair into the map incase it is not already
in the map.
@return will return a bool depending on the insert, if it can insert the pair
it will return true otherwise false
*/
func (l *Leaf) Insert(key int, value string) bool {
	if !l.Contains(key) {
		l.Data[key] = value
		return true
	}
	return false
}

/*
Size will return the size of map of the leaf,
so who many paris are stored in the map
@return will be an int value
*/
func (l *Leaf) Size() int {
	return len(l.Data)
}

/*
Erase will erase a given key from
*/
func (l *Leaf) Erase(key int) bool {
	if l.Contains(key) {
		delete(l.Data, key)
		return true
	}
	return false
}

/*
getData will return a pointer to the data of the given leaf
*/
func (l *Leaf) getData() *map[int]string {
	return &l.Data
}

/*
Change will change the value of a given key in the leaf
*/
func (l *Leaf) Change(key int, value string) bool {
	if l.Contains(key) {
		(*l.getData())[key] = value
		return true
	}
	return false
}

/*
Find will return a value by a given key
*/
func (l *Leaf) Find(key int) string {
	v, _ := (*l.getData())[key]
	return v
}
