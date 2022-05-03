package Brodal

import (
	"container/list"
)

var UPPER_BOUND uint = 7
var LOWER_BOUND uint = 4

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
		// if operations,newVal := tree.upperBoundGuide.increase(node.rank, 1); uint(newVal) == UPPER_BOUND {

			// for i := range operations {
				// tree.reduce(i)
			// }
			// tree.upperBoundGuide.reduce(node.rank, 3)

			// xNode, yNode := tree.childrenRank[node.rank], tree.childrenRank[node.rank].rightBrother
			// minNode := node

			// if minNode.value > xNode.value {
			// 	minNode = xNode
			// 	xNode = node
			// }

			// if minNode.value > yNode.value {
			// 	temp := minNode
			// 	minNode = yNode
			// 	yNode = temp
			// }

			// minNode.link(xNode, yNode)

			// tree.remove(tree.childrenRank[node.rank].rightBrother)
			// tree.remove(tree.childrenRank[node.rank])

			// tree.insert(minNode)
		// }
	} else {

	}
}

func (tree *rootNode) remove (node *node) {

}
