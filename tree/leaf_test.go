package tree

import (
	"testing"
)

func TestContains(t *testing.T) {
	data := NewLeaf()
	data.Data[1] = "test"
	ret := data.Contains(1)
	if !ret {
		t.Errorf("Expected return after contains is = %t; got %t", true, ret)
	}
}
func TestInsertMap(t *testing.T) {
	data := NewLeaf()
	testmap := make(map[int]string)
	testmap[1] = "test"
	data.insertMap(testmap)
	ret := data.Contains(1)
	if !ret {
		t.Errorf("Expected return after insertmap is = %t; got %t", true, ret)
	}
}

func TestInsert(t *testing.T) {
	data := NewLeaf()
	ret := data.Insert(1, "test")
	ret2 := data.Contains(1)
	if !ret || !ret2 {
		t.Errorf("Expected return after insert is = %t %t; got %t %t", true, true, ret, ret2)
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
		t.Errorf("Expected return after erase is = %t; got %t", true, ret)
	}
}

func TestFind(t *testing.T) {
	data := NewLeaf()
	data.Insert(1, "test")
	ret2 := data.Find(1)
	if ret2 != "test" {
		t.Errorf("Expected return after change is = %s; got %s", "test2", ret2)
	}
}

func TestChange(t *testing.T) {
	data := NewLeaf()
	data.Insert(1, "test")
	ret1 := data.Change(1, "test2")
	ret2 := data.Find(1)
	if !ret1 || ret2 != "test2" {
		t.Errorf("Expected return after change is = %t and %s; got %t and %s", true, "test2", ret1, ret2)
	}
}
