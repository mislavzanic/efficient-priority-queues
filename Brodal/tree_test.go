package Brodal

import (
	"fmt"
	"testing"
)

func TestNewTree(t *testing.T) {
	root := newTree(1, 1)
	root.linkToRoot(newNode(2), newNode(3))

	if root.RootRank() != 1 {
		t.Error(fmt.Sprintf("Tree rank is %d, not 1", root.RootRank()))
	}

	if root.childrenRank[root.RootRank() - 1].value != 2 {
		t.Error(fmt.Sprintf("Tree first child value is %f, not 2", root.childrenRank[root.RootRank() - 1].value))
	}

	root.addRootChild(newNode(5))
	if root.childrenRank[root.RootRank() - 1].rightBrother().value != 5 {
		t.Error(fmt.Sprintf("Tree new child value is %f, not 5", root.childrenRank[root.RootRank() - 1].rightBrother().value))
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

	nodes := root.delink()
	for _, n := range nodes {
		if n.rank != 1 {
			t.Error(fmt.Sprintf("node rank is %d, not 1", n.rank))
		}
	}
}
