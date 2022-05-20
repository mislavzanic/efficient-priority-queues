package Brodal

import (
	"fmt"
	"testing"
)

func TestTree(t *testing.T) {
	root := newTree(1, 1)

	root.InsertNodes(newNode(2), newNode(3))

	if root.RootRank() != 1 {
		t.Error(fmt.Sprintf("ERROR: Tree rank is %d, not 1", root.RootRank()))
	}

	if root.childrenRank[root.RootRank()-1].value != 2 {
		t.Error(fmt.Sprintf("ERROR: Tree first child value is %f, not 2", root.childrenRank[root.RootRank()-1].value))
	}

	root.InsertNodes(newNode(5))
	if root.childrenRank[root.RootRank()-1].RightBrother().value != 5 {
		t.Error(fmt.Sprintf("ERROR: Tree new child value is %f, not 5", root.childrenRank[root.RootRank()-1].RightBrother().value))
	}
	root.InsertNodes(newNode(6), newNode(7), newNode(8))

	if root.NumOfRootChildren(0) != 6 {
		t.Error(fmt.Sprintf("ERROR: Num of children is %d, not 6", root.NumOfRootChildren(0)))
	}

	root.Link(0)
	if root.RootRank() != 2 {
		t.Error(fmt.Sprintf("Tree rank is %d, not 2", root.RootRank()))
	}

	if root.childrenRank[root.RootRank()-1].value != 2 {
		t.Error(fmt.Sprintf("ERROR: Tree first child value is %f, not 2", root.childrenRank[root.RootRank()-1].value))
	}

	if root.childrenRank[root.RootRank()-1].LeftChild().value != 8 {
		t.Error(fmt.Sprintf("ERROR: Tree first child value is %f, not 8", root.childrenRank[root.RootRank()-1].LeftChild().value))
	}

	if root.childrenRank[root.RootRank()-1].LeftChild().RightBrother().value != 7 {
		t.Error(fmt.Sprintf("ERROR: Tree first child value is %f, not 7", root.childrenRank[root.RootRank()-1].LeftChild().value))
	}

	node1 := newNode(100)
	node1.link(newNode(100), newNode(100))
	node4 := newNode(200)
	node4.link(newNode(200), newNode(200))
	node3 := newNode(300)
	node3.link(newNode(300), newNode(300))
	node1.link(node4, node3)

	node2 := newNode(100)
	node2.link(newNode(100), newNode(100))
	node4 = newNode(200)
	node4.link(newNode(200), newNode(200))
	node3 = newNode(300)
	node3.link(newNode(300), newNode(300))
	node2.link(node4, node3)

	if node1.rank != 2 {
		t.Error(fmt.Sprintf("Node %d, not 2", node1.rank))
	}

	root.InsertNodes(node1, node2)
	if root.RootRank() != 3 {
		t.Error(fmt.Sprintf("Tree rank is %d, not 3", root.RootRank()))
	}

	if root.LeftChild().rank != 2 {
		t.Error(fmt.Sprintf("Node rank is %d, not 2", root.LeftChild().rank))
	}

	if root.LeftChild().value != 100 {
		t.Error(fmt.Sprintf("Node value is %f, not 100", root.LeftChild().value))
	}

	if root.NumOfRootChildren(2) != 2 {
		t.Error(fmt.Sprintf("ERROR: Num of children is %d, not 2", root.NumOfRootChildren(0)))
	}

	nodes := root.Delink()
	for _, n := range nodes {
		if n.rank != 2 {
			t.Error(fmt.Sprintf("node rank is %d, not 1", n.rank))
		}
		if n.value != 100 {
			t.Error(fmt.Sprintf("node value is %f, not 100", n.value))
		}
	}

	rank1Nodes := root.DecRank()
	for _, n := range rank1Nodes {
		if n.rank != 1 {
			t.Error("Nije 1")
		}
	}
}
