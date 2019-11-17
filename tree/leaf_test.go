package tree

import (
	"testing"
)

func TestContains(t *testing.T) {
	data := NewLeaf()
	data.Data[1] = "test"
	if !data.Contains(1) {
		t.Errorf("Expected return after contains is = %t; got %t", true, false)
	}
}
func TestInsertMap(t *testing.T) {
	data := NewLeaf()
	testmap := make(map[int]string)
	testmap[1] = "test"
	data.insertMap(testmap)
	if !data.Contains(1) {
		t.Errorf("Expected return after insertmap is = %t; got %t", true, false)
	}
}

func TestInsert(t *testing.T) {
	data := NewLeaf()
	ret := data.Insert(1, "test")
	ret2 := data.Contains(1)
	if !ret || !ret2 {
		t.Errorf("Expected return after insert is = %t; got %t", true, false)
	}
}

func TestSize(t *testing.T) {
	data := NewLeaf()
	data.Insert(1, "test")
	ret := data.Size()
	if ret != 1 {
		t.Errorf("Expected size is = %d; got %d", 1, ret)
	}
}

func TestErase(t *testing.T) {
	data := NewLeaf()
	data.Insert(1, "test")
	ret := data.Erase(1)
	if !ret {
		t.Errorf("Expected return afer erase is = %t; got %t", true, false)
	}
}

func TestFind(t *testing.T) {
	data := NewLeaf()
	data.Insert(1, "test")
	ret2 := data.Find(1)
	if ret2 != "test" {
		t.Errorf("Expected return afer change is = %t and %s; got %t and %s", true, "test2", false, ret2)
	}
}

func TestChange(t *testing.T) {
	data := NewLeaf()
	data.Insert(1, "test")
	ret1 := data.Change(1, "test2")
	ret2 := data.Find(1)
	if !ret1 || ret2 != "test2" {
		t.Errorf("Expected return afer change is = %t and %s; got %t and %s", true, "test2", false, ret2)
	}
}

/*
func TestAddNewTrie(t *testing.T) {
	context := actor.EmptyRootContext
	props := actor.PropsFromProducer(func() actor.Actor {
		return tree.CreateBasicNode(5)}
	)
	root := context.Spawn(props)
	a, b, v := treeservice.AddNewTrie(5)
	if ret2 != "test2" {
		t.Errorf("Expected return afer change is = %t and %s; got %t and %s", true, "test2", false, ret2)
	}
}*/

func TestAddNewTrie(t *testing.T) {
	ac := CreateBasicNode(5)

	ret := ac.(*Nodeactor).getLimit()

	if ret != 5 {
		t.Errorf("Expected return afer getLimit is = %d ; got %d", 5, ret)
	}
}

func TestReturnAllKey(t *testing.T) {
	testmap := make(map[int]string)
	testmap[1] = "test1"
	testmap[3] = "test3"
	testmap[5] = "test5"
	testmap[6] = "test6"

	ret := returnAllKey(testmap)

	if ret[0] != 1 || ret[1] != 3 || ret[2] != 5 || ret[3] != 3 {
		t.Errorf("Expected return afer getLimit is = %d %d %d %d; got %d %d %d %d", 1, 3, 5, 6, ret[0], ret[1], ret[2], ret[3])
	}
}
