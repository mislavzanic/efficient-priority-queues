package Brodal

import (
	"container/list"
)

var UPPER_BOUND uint = 7
var LOWER_BOUND uint = 4

type reduceType byte
const (
	linkReduce   reduceType = 0
	delinkReduce            = 1
)

type rootNode struct {
	root          *node

	rankPointersW *list.List
	childrenRank  []*node

	upperBoundGuide *guide
	lowerBoundGuide *guide
	mainTreeGuideW  *guide
}

func newRootNode(value float64, treeIndex uint) *rootNode {

	rootNode := &rootNode{
		root: newNode(value),
		rankPointersW: list.New(),
		childrenRank: nil,
		upperBoundGuide: newGuide(UPPER_BOUND),
		lowerBoundGuide: newGuide(LOWER_BOUND),
		mainTreeGuideW: nil,
	}

	if treeIndex == 1 {
		rootNode.mainTreeGuideW = newGuide(6)
	}

	return rootNode
}

type BrodalHeap struct {
    tree1 *rootNode
    tree2 *rootNode
}

func newHeap(value float64) *BrodalHeap {
	return &BrodalHeap{
		tree1: newRootNode(value, 1),
		tree2: nil,
	}
}

func (bh *BrodalHeap) Min() float64 {
	return bh.tree1.root.value
}

func (bh *BrodalHeap) Insert(value float64) {
	otherHeap := newHeap(value)
	bh.Meld(otherHeap)
}

func (bh *BrodalHeap) Meld(other *BrodalHeap) {
	minTree := func () *rootNode {
		if bh.tree1.root.value > other.tree1.root.value {
			return other.tree1
		} else {
			return bh.tree1
		}
	}()

	maxRankTree := func () *rootNode {
		if bh.tree1.root.rank > other.tree1.root.rank {
			return bh.tree1
		} else {
			return other.tree1
		}
	}()

	if bh.tree2 == nil && other.tree2 == nil {
		if minTree == bh.tree1 {
			minTree.insert(other.tree1.root)
		} else {
			minTree.insert(bh.tree1.root)
		}

		bh.tree1 = minTree

	} else {

	}
}

func (tree *rootNode) insert(node *node) {
	if node.rank < tree.root.rank - 2 {
		actions := tree.upperBoundGuide.forceIncrease(node.rank, 3)
		for _, act := range actions {
			tree.performeAction(node, act, linkReduce)
		}
		tree.updateHighRank()
	} else {

	}
}

func (tree *rootNode) remove(node *node) {

}

func (tree *rootNode) performeAction(node *node, action action, reduceType reduceType) {
	if reduceType == linkReduce {
		tree.link(action.index)
	} else {
		tree.delink(action.index)
	}
}

func (tree *rootNode) link(rank uint) {
	nodeX := tree.childrenRank[rank]
	nodeY, nodeZ := nodeX.rightBrother, nodeX.rightBrother.rightBrother

	if nodeZ.rightBrother.rank == rank {
		tree.childrenRank[rank] = nodeZ.rightBrother
	} else {
		tree.childrenRank[rank] = nil
	}

	minNode, nodeX, nodeY := getMinNode(nodeX, nodeY, nodeZ)

	minNode.link(nodeX, nodeY)
}

func (tree *rootNode) delink(rank uint) {

}

func (tree *rootNode) updateHighRank() {

}
