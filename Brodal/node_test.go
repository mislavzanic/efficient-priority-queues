package Brodal

import (
	"fmt"
	"testing"
)

func TestInsertOfFirstChildren(t *testing.T) {
	testNode := newNode(1)
	if _, err := testNode.link(newNode(2), newNode(3)); err != nil {
		t.Error(err.Error())
	}

	if testNode.value != 1 {
		t.Error(fmt.Sprintf("Test node value is %f", testNode.value))
	}

	if testNode.rank != 1 {
		t.Error(fmt.Sprintf("Test node rank is %d, not 1", testNode.rank))
	}

	if testNode.numOfChildren[0] != 2 {
		t.Error(fmt.Sprintf("Test node numOfChildren[0] is %d, not 2", testNode.numOfChildren[0]))
	}

	if testNode.children.Len() != 2 {
		t.Error(fmt.Sprintf("Test node children.Len() is %d, not 2", testNode.children.Len()))
	}

	for e := testNode.children.Back(); e != nil; e = e.Prev() {
		if e.Value.(*node).rank != 0 {
			t.Error(fmt.Sprintf("Child node rank is %d, not 0", e.Value.(*node).rank))
		}
	}

	if lc, e := testNode.leftChild(); e == nil {
		testNode.addBrother(newNode(5), lc, true)
	}

	if lc, e := testNode.leftChild(); e == nil {
		if lc.value != 2 {
			t.Error(fmt.Sprintf("First child value is %f", lc.value))
		}
	} else {
		t.Error(e.Error())
	}


	if lc, err := testNode.leftChild(); err == nil {
		if rb, _ := lc.rightBrother(); rb.value != 5 {
			t.Error(fmt.Sprintf("Second child value is %f", rb.value))
		} else {
			if rrb, _ := rb.rightBrother(); rrb.value != 3 {
				t.Error(fmt.Sprintf("Third child value is %f", rrb.value))
			}
		}

	} else {
		t.Error(err.Error())
	}

}

func TestRemoveChildren(t *testing.T) {
	testNode := newNode(1)
	testNode.link(newNode(2), newNode(3))
	testNode.addBrother(newNode(5), testNode.LeftChild(), true)

	testNode.removeChild(testNode.LeftChild())

	if testNode.numOfChildren[0] != 2 {
		t.Error(fmt.Sprintf("Test node numOfChildren[0] is %d, not 2", testNode.numOfChildren[0]))
	}

	if testNode.children.Len() != 2 {
		t.Error(fmt.Sprintf("Test node children.Len() is %d, not 2", testNode.children.Len()))
	}

	if testNode.LeftChild().value != 5 {
		t.Error(fmt.Sprintf("First child value is %f", testNode.LeftChild().value))
	}

	if testNode.rank != 1 {
		t.Error(fmt.Sprintf("Test node rank is %d, not 1", testNode.rank))
	}

	testNode.removeChild(testNode.LeftChild())

	if testNode.LeftChild().value != 3 {
		t.Error(fmt.Sprintf("First child value is %f", testNode.LeftChild().value))
	}

	testNode.removeChild(testNode.LeftChild())

	if testNode.rank != 0 {
		t.Error(fmt.Sprintf("Test node rank is %d, not 0", testNode.rank))
	}

	if testNode.numOfChildren[0] != 0 {
		t.Error(fmt.Sprintf("Test node numOfChildren[0] is %d, not 0", testNode.numOfChildren[0]))
	}

	testNode.link(newNode(2), newNode(3))

	if testNode.numOfChildren[0] != 2 {
		t.Error(fmt.Sprintf("Test node numOfChildren[0] is %d, not 2", testNode.numOfChildren[0]))
	}

	if testNode.children.Len() != 2 {
		t.Error(fmt.Sprintf("Test node children.Len() is %d, not 2", testNode.children.Len()))
	}
}

func TestLink(t *testing.T) {
	testNode := newNode(1)
	testNode.link(newNode(2), newNode(3))

	testNode2 := newNode(100)
	testNode2.link(newNode(200), newNode(300))

	testNode3 := newNode(123)
	testNode3.link(newNode(1000), newNode(2000))

	testNode.link(testNode2, testNode3)

	if testNode.numOfChildren[0] != 2 {
		t.Error(fmt.Sprintf("Test node numOfChildren[0] is %d, not 2", testNode.numOfChildren[0]))
	}

	if testNode.numOfChildren[1] != 2 {
		t.Error(fmt.Sprintf("Test node numOfChildren[0] is %d, not 2", testNode.numOfChildren[1]))
	}

	if lc,_ := testNode.leftChild(); lc.rank != 1 {
		t.Error(fmt.Sprintf("First child rank is %d, not 1", lc.rank))
	}

	if testNode.rank != 2 {
		t.Error(fmt.Sprintf("Test node rank is %d, not 2", testNode.rank))
	}
}

func TestDelink(t *testing.T) {
	testNode := newNode(1)
	testNode.link(newNode(2), newNode(3))

	testNode2 := newNode(100)
	testNode2.link(newNode(200), newNode(300))

	testNode3 := newNode(123)
	testNode3.link(newNode(1000), newNode(2000))

	testNode.link(testNode2, testNode3)

	if nodes, err := testNode.delink(); err == nil {
		if testNode.rank != 1 {
			t.Error(fmt.Sprintf("Test node rank is %d, not 1", testNode.rank))
		}

		for _, n := range nodes {
			if n.rank != 1 {
				t.Error(fmt.Sprintf("node rank is %d, not 1", n.rank))
			}
		}
	} else {
		t.Error(fmt.Sprint(err.Error()))
	}
}
