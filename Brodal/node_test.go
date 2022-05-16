package Brodal

import (
	"fmt"
	"testing"
)

func TestInsertOfFirstChildren(t *testing.T) {
	testNode := newNode(1)
	testNode.link(newNode(2), newNode(3))
	// testNode.addFirstChildren(newNode(2), newNode(3))

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

	testNode.addBrother(newNode(5), testNode.leftSon(), true)
	if testNode.leftSon().value != 2 {
		t.Error(fmt.Sprintf("First child value is %f", testNode.leftSon().value))
	}

	if testNode.leftSon().rightBrother().value != 5 {
		t.Error(fmt.Sprintf("Second child value is %f", testNode.leftSon().rightBrother().value))
	}

	if testNode.leftSon().rightBrother().rightBrother().value != 3 {
		t.Error(fmt.Sprintf("Third child value is %f", testNode.leftSon().rightBrother().rightBrother().value))
	}
}

func TestRemoveChildren(t *testing.T) {
	testNode := newNode(1)
	testNode.link(newNode(2), newNode(3))
	testNode.addBrother(newNode(5), testNode.leftSon(), true)

	testNode.removeChild(testNode.leftSon())

	if testNode.numOfChildren[0] != 2 {
		t.Error(fmt.Sprintf("Test node numOfChildren[0] is %d, not 2", testNode.numOfChildren[0]))
	}

	if testNode.children.Len() != 2 {
		t.Error(fmt.Sprintf("Test node children.Len() is %d, not 2", testNode.children.Len()))
	}

	if testNode.leftSon().value != 5 {
		t.Error(fmt.Sprintf("First child value is %f", testNode.leftSon().value))
	}

	if testNode.rank != 1 {
		t.Error(fmt.Sprintf("Test node rank is %d, not 1", testNode.rank))
	}

	testNode.removeChild(testNode.leftSon())

	if testNode.leftSon().value != 3 {
		t.Error(fmt.Sprintf("First child value is %f", testNode.leftSon().value))
	}

	testNode.removeChild(testNode.leftSon())

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

	if testNode.leftSon().rank != 1 {
		t.Error(fmt.Sprintf("First child rank is %d, not 1", testNode.leftSon().rank))
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

	nodes := testNode.delink()

	if testNode.rank != 1 {
		t.Error(fmt.Sprintf("Test node rank is %d, not 1", testNode.rank))
	}

	for _, n := range nodes {
		if n.rank != 1 {
			t.Error(fmt.Sprintf("node rank is %d, not 1", n.rank))
		}
	}
}
