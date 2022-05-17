package Brodal

import (
	"fmt"
	"testing"
)

func TestTree(t *testing.T) {
	root := newTree(1, 1)
	root.linkToRoot(newNode(2), newNode(3))

	if root.RootRank() != 1 {
		t.Error(fmt.Sprintf("Tree rank is %d, not 1", root.RootRank()))
	}

	if root.childrenRank[root.RootRank()-1].value != 2 {
		t.Error(fmt.Sprintf("Tree first child value is %f, not 2", root.childrenRank[root.RootRank()-1].value))
	}

	root.addRootChild(newNode(5))
	if root.childrenRank[root.RootRank()-1].rightBrother().value != 5 {
		t.Error(fmt.Sprintf("Tree new child value is %f, not 5", root.childrenRank[root.RootRank()-1].rightBrother().value))
	}
	root.addRootChild(newNode(6))
	root.addRootChild(newNode(7))
	root.addRootChild(newNode(8))

	root.link(0)
	if root.RootRank() != 2 {
		t.Error(fmt.Sprintf("Tree rank is %d, not 2", root.RootRank()))
	}

	root.addRootChild(newNode(6))
	root.addRootChild(newNode(7))
	root.link(0)

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

	root.incRank(node1, node2)
	if root.RootRank() != 3 {
		t.Error(fmt.Sprintf("Tree rank is %d, not 3", root.RootRank()))
	}

	nodes := root.delink()
	for _, n := range nodes {
		if n.rank != 2 {
			t.Error(fmt.Sprintf("node rank is %d, not 1", n.rank))
		}
		if n.value != 100 {
			t.Error(fmt.Sprintf("node value is %f, not 100", n.value))
		}
	}
}
