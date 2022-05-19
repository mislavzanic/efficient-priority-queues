package Brodal

import (
	// "math"
	// "fmt"
)

type valType float64

type BrodalHeap struct {
	t1s   *tree1Struct
	tree2 *tree
}

const ALPHA int = 10

func NewEmptyHeap() *BrodalHeap {
	return &BrodalHeap{
		t1s:   newEmptyT1S(),
		tree2: nil,
	}
}

func NewHeap(value valType) *BrodalHeap {
	return &BrodalHeap{
		t1s:   newT1S(value),
		tree2: nil,
	}
}

func (bh *BrodalHeap) getTree(index int) *tree {
	if index == 1 {
		return bh.t1s.getTree()
	}
	return bh.tree2
}

func (bh *BrodalHeap) Empty() bool {
	return bh.getTree(1) == nil
}

func (bh *BrodalHeap) Insert(value valType) {
	bh.Meld(NewHeap(value))
}

func (bh *BrodalHeap) Meld(otherHeap *BrodalHeap) {
	if bh.Empty() {
		bh = otherHeap
	} else {
		minT1S, otherT1S := getMinValTree(bh.t1s, otherHeap.t1s)
		if bh.getTree(2) == nil && otherHeap.getTree(2) == nil {
			bh.t1s = minT1S
			if minT1S.getTree().RootRank() <= otherT1S.getTree().RootRank() {
				bh.tree2 = otherT1S.getTree()
			} else {
				// insert otherT1S.tree1 u bh.t1s.tree1
			}
		} else {

		}
	}
}

func (bh *BrodalHeap) insertNode(treeIndex int, child *node) {
	if child.rank < bh.getTree(treeIndex).RootRank() - 2 {
		actions, _ := bh.getTree(treeIndex).askGuide(child.rank, bh.getTree(treeIndex).numOfRootChildren(child.rank) + 1, true)
		for _, act := range actions {
			if act.op == Reduce {
				bh.linkNodes(treeIndex, act.index)
			} else {
				bh.getTree(treeIndex).insertNode(child)
				bh.mbyAddViolation(child)
			}
		}
	} else {
		bh.getTree(treeIndex).insertNode(child)
		bh.mbyAddViolation(child)
	}

	bh.updateHighRank(treeIndex, bh.getTree(treeIndex).RootRank() - 2)
	bh.updateHighRank(treeIndex, bh.getTree(treeIndex).RootRank() - 1)
}

func (bh *BrodalHeap) linkNodes(treeIndex int, rank int) {
    minNode, notViolations := bh.getTree(treeIndex).linkRank(rank)
	bh.mbyAddViolation(minNode)
	for _, nv := range notViolations {
		bh.mbyRemoveViolation(nv)
	}
}

func (bh *BrodalHeap) updateHighRank(treeIndex int, rank int) {
	noChildren := bh.getTree(treeIndex).numOfRootChildren(rank)
	if noChildren < 0 { return }

	if noChildren > 7 {
		bh.linkNodes(treeIndex, rank)
		bh.linkNodes(treeIndex, rank)
	} else if noChildren < 2 {

	} else {
		return
	}

	if rank == bh.getTree(treeIndex).RootRank() - 2 {
		bh.updateHighRank(treeIndex, rank + 1)
	}
}

func (bh *BrodalHeap) mbyAddViolation(child *node) {
	if child.isGood() { return }

	if child.rank >= bh.getTree(1).RootRank() {
		bh.addToV(child)
	} else {
		bh.addToW(child)
	}
}
