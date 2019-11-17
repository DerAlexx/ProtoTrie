package tree

import (
	"reflect"
	"testing"
)

func TestAddNewTrie(t *testing.T) {
	ac := CreateBasicNode(5)

	ret := ac.(*Nodeactor).getLimit()

	if ret != 5 {
		t.Errorf("Expected return after getLimit is = %d ; got %d", 5, ret)
	}
}

func TestReturnAllKey(t *testing.T) {
	testmap := make(map[int]string)
	testmap[1] = "test1"
	testmap[3] = "test3"
	testmap[5] = "test5"
	testmap[6] = "test6"

	ret := returnAllKey(testmap)

	ret2 := []int{1, 3, 5, 6}

	abc := reflect.DeepEqual(ret, ret2)

	if !abc {
		t.Errorf("Expected return after returnallkeys is = %t; got %t", true, abc)
	}
}

func TestSortMap(t *testing.T) {
	testmap1 := make(map[int]string)
	testmap1[1] = "test1"
	testmap1[3] = "test3"
	testmap1[5] = "test5"

	t1, t2, t3 := sortMap(testmap1)

	testmap2 := make(map[int]string)
	testmap2[10] = "test10"
	testmap2[13] = "test13"
	testmap2[20] = "test20"
	testmap2[45] = "test45"

	r1, r2, r3 := sortMap(testmap2)

	if t1[1] != "test1" || t1[3] != "test3" || t2[5] != "test5" || t3 != 3 || r1[10] != "test10" || r1[13] != "test13" || r2[20] != "test20" || r2[45] != "test45" || r3 != 13 {
		t.Error(t1, t2, t3, r1, r2, r3)

	}
}

func TestGetConstrain(t *testing.T) {

	ac := CreateBasicNode(5)

	ret := ac.(*Nodeactor).getConstrain()

	if ret != -1 {
		t.Errorf("Expected return after getconstrain is = %d; got %d", -1, ret)
	}
}

func TestSetStoreable(t *testing.T) {

	ac := CreateBasicNode(5)
	a := ac.(*Nodeactor)
	a.SetStoreable(10)

	ret := a.getConstrain()

	if ret != 10 {
		t.Errorf("Expected return after setstorable is = %d; got %d", 10, ret)
	}
}

func TestHasValueToDecide(t *testing.T) {

	ac := CreateBasicNode(5)
	a := ac.(*Nodeactor)
	a.SetStoreable(10)

	ret1, ret2 := a.HasValueToDecide()

	if !ret1 || ret2 != 10 {
		t.Errorf("Expected return after hasvaluetodecide is = %t %d got %t %d", true, 10, ret1, ret2)
	}
}

func TestFindBiggestKey(t *testing.T) {
	testmap := make(map[int]string)
	testmap[1] = "test1"
	testmap[3] = "test3"
	testmap[5] = "test5"
	testmap[6] = "test6"

	ret := findBiggestKey(testmap)

	if ret != 6 {
		t.Errorf("Expected return after findbiggestkey is = %d; got %d", 6, ret)
	}
}

/*
func TestIsLeft(t *testing.T) {

	if !abc {
		t.Errorf("Expected return after returnallkeys is = %t; got %t", true, false)
	}
}

func TestTraverseChild(t *testing.T) {

	if !abc {
		t.Errorf("Expected return after returnallkeys is = %t; got %t", true, false)
	}
}

func TestUnionsLeft(t *testing.T) {

	if !abc {
		t.Errorf("Expected return after returnallkeys is = %t; got %t", true, false)
	}
}

func TestIsFull(t *testing.T) {

	if !abc {
		t.Errorf("Expected return after returnallkeys is = %t; got %t", true, false)
	}
}

func TestExpand(t *testing.T) {

	if !abc {
		t.Errorf("Expected return after returnallkeys is = %t; got %t", true, false)
	}
}
*/
