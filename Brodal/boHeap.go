package Brodal

import (
	"errors"
	"fmt"
)

// "math"
// "fmt"

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
		return bh.t1s.GetTree()
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
		bh.t1s = minT1S
		if bh.getTree(2) == nil && otherHeap.getTree(2) == nil {
			if minT1S.GetTree().RootRank() <= otherT1S.GetTree().RootRank() {
				bh.tree2 = otherT1S.GetTree()
			} else {
				bh.insertNode(1, otherT1S.GetTree().root)
			}
		} else {
			maxRankTree, otherT2 := getMaxTree(bh.getTree(2), otherHeap.getTree(2))
			insertIndex := 2

			if maxRankTree.RootRank() == minT1S.GetTree().RootRank() {
				bh.insertNode(1, maxRankTree.root)
				bh.tree2 = nil
				insertIndex = 1
			} else {
				bh.tree2 = maxRankTree
			}

			if otherT1S.GetTree() != nil {
				bh.insertNode(insertIndex, otherT1S.GetTree().root)
				if len(otherT2) != 0 {
					if otherT2[0].RootRank() == bh.getTree(insertIndex).RootRank() {
						nodes := otherT2[0].DecRank()
						for _, n := range nodes {
							bh.insertNode(insertIndex, n)
						}
					}
					bh.insertNode(insertIndex, otherT2[0].root)
				}
			}
			bh.createAlphaSpace()
		}
	}
}


func (bh *BrodalHeap) insertNode(treeIndex int, child *node) {
	bh.heapAction(treeIndex, child, true)
}

func (bh *BrodalHeap) cutOffNode(treeIndex int, child *node) *node {
	bh.heapAction(treeIndex, child, false)
	return child
}

func (bh *BrodalHeap) heapAction(treeIndex int, child *node, insert bool) {
	if child.rank < bh.getTree(treeIndex).RootRank() - 2 {
		bh.lowRankAction(treeIndex, child, insert)
	} else {
		bh.highRankAction(treeIndex, child, insert)
	}

	bh.updateHighRank(treeIndex, bh.getTree(treeIndex).RootRank() - 2)
	bh.updateHighRank(treeIndex, bh.getTree(treeIndex).RootRank() - 1)
}

func (bh *BrodalHeap) lowRankAction(treeIndex int, child *node, insert bool) {
	actionsInc, actionsDec := bh.getTree(treeIndex).AskGuide(child.rank, bh.getTree(treeIndex).NumOfRootChildren(child.rank) + 1, insert)
	bh.handleActions(treeIndex, append(actionsInc, actionsDec...), child)
}

func (bh *BrodalHeap) highRankAction(treeIndex int, child *node, insert bool) {
	if insert {
		bh.getTree(treeIndex).insertNode(child)
		bh.mbyAddViolation(child)
	} else {
		if _, err := bh.getTree(treeIndex).cutOffNode(child); err != nil {
			panic(err.Error())
		}
		bh.mbyRemoveFromViolating(child)
	}
}

func (bh *BrodalHeap) handleActions(treeIndex int, actions []action, child *node) {
	for _, act := range actions {
		switch op := act.op; op {
		case Reduce:
			bh.doReduceAction(treeIndex, act.bound, act.index)
		default:
			bh.doDefaultAction(treeIndex, act.bound, child)
		}
	}
}

func (bh *BrodalHeap) doReduceAction(treeIndex int, guideBound int, rank int) {
	switch guideBound {
	case UPPER_BOUND:
		bh.linkNodes(treeIndex, rank)
	case LOWER_BOUND:
		child, nodes := bh.delinkNodes(treeIndex, rank)
		for _, n := range nodes {
			bh.insertNode(treeIndex, n)
		}
		bh.insertNode(treeIndex, child)
	case GUIDE_BOUND:
		if err := bh.reduceW(rank); err != nil {
			panic(err.Error())
		}
	default:
		panic(fmt.Sprintf("Wrong bound: %d", guideBound))
	}
}

func (bh *BrodalHeap) doDefaultAction(treeIndex int, guideBound int, child *node) {
	switch guideBound {
	case UPPER_BOUND:
		if err := bh.getTree(treeIndex).insertNode(child); err != nil {
			panic(err.Error())
		}
		bh.mbyAddViolation(child)
	case LOWER_BOUND:
		if _, err := bh.getTree(treeIndex).cutOffNode(child); err != nil {
			panic(err.Error())
		}
		bh.mbyRemoveFromViolating(child)
	case GUIDE_BOUND:
		if err := bh.insertNewW(child); err != nil {
			panic(err.Error())
		}
	default:
		panic(fmt.Sprintf("Wrong bound: %d", guideBound))
	}
}

func (bh *BrodalHeap) linkNodes(treeIndex int, rank int) {
    minNode := bh.getTree(treeIndex).Link(rank)
	bh.mbyAddViolation(minNode)
	bh.mbyRemoveFromViolating(minNode.LeftChild())
	bh.mbyRemoveFromViolating(minNode.LeftChild().RightBrother())
}

func (bh *BrodalHeap) delinkNodes(treeIndex int, rank int) (*node, []*node) {
	child := bh.getTree(treeIndex).RemoveChild(rank)
	bh.mbyRemoveFromViolating(child)
	nodes, err := child.delink()
	if err != nil {
		panic(err.Error())
	}

	return child, nodes
}

func (bh *BrodalHeap) delink(treeIndex int) []*node {
	nodes := bh.getTree(treeIndex).Delink()
	for _, n := range nodes {
		bh.mbyRemoveFromViolating(n)
	}
	return nodes
}

func (bh *BrodalHeap) updateHighRank(treeIndex int, rank int) {
	if rank < 0 { return }
	noChildren := bh.getTree(treeIndex).NumOfRootChildren(rank)
	if noChildren > 7 {
		bh.linkNodes(treeIndex, rank)
		bh.linkNodes(treeIndex, rank)
	} else if noChildren < 2 && noChildren >= 0 {
		if rank == bh.getTree(treeIndex).RootRank() - 1 {
			nodes := bh.getTree(treeIndex).DecRank()
			for _, n := range nodes {
				bh.mbyRemoveFromViolating(n)
				bh.insertNode(treeIndex, n)
			}
		} else {
			bh.updateHighRank(treeIndex, rank+1)
			if rank == bh.getTree(treeIndex).RootRank() - 2 {
				nodes := bh.delink(treeIndex)
				for _, n := range nodes {
					bh.insertNode(treeIndex, n)
				}
			}
		}
	} else {
		panic("Negative noChildren")
	}

	if rank == bh.getTree(treeIndex).RootRank() - 2 {
		bh.updateHighRank(treeIndex, rank + 1)
	}
}


//////////////////////////////////////////////////////////////////////////////
////////////////////////Violation Handling////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////

func (bh *BrodalHeap) mbyAddViolation(child *node) {
	if child.isGood() { return }

	if child.rank >= bh.getTree(1).RootRank() {
		bh.addToV(child)
	} else {
		bh.addToW(child)
	}
}

func (bh *BrodalHeap) mbyRemoveFromViolating(notBad *node) {
	if !notBad.isGood() { return }

	if notBad.parentViolatingList == bh.getTree(1).root.wList {
		bh.removeFromW(notBad)
	}

	notBad.removeSelfFromViolating()
}

func (bh *BrodalHeap) addToV(child *node) {
	child.violatingSelf = bh.getTree(1).vList().PushBack(child)
	child.parentViolatingList = bh.getTree(1).vList()
}

func (bh *BrodalHeap) createAlphaSpace() {
	if bh.getTree(2) == nil || bh.getTree(1).vList().Len() <= ALPHA * bh.getTree(1).RootRank() {
		return
	}

	if bh.getTree(2).RootRank() == 0 {
		panic("Tree2 is of rank 0")
	}

	if bh.getTree(2).RootRank() <= bh.getTree(1).RootRank() + 1 {
		nodes := bh.getTree(2).DecRank()
		for _, n := range nodes {
			bh.insertNode(1, n)
		}
		bh.insertNode(1, bh.getTree(2).root)
		bh.tree2 = nil
		return
	}

	newChild := bh.cutOffNode(2, bh.getTree(2).childrenRank[bh.getTree(1).RootRank() + 1])
	nodes := newChild.DecRank()
	for _, n := range nodes {
		bh.insertNode(1, n)
	}

	bh.insertNode(1, newChild)
}

func (bh *BrodalHeap) addToW(child *node) {
	actions := bh.t1s.GetWGuide().forceIncrease(child.rank, bh.t1s.GetWNums(child.rank) + 1, 2)
	bh.handleActions(1, actions, child)
}

func (bh *BrodalHeap) removeFromW(child *node) error {
	return bh.t1s.removeFromW(child)
}

func (bh *BrodalHeap) reduceW(rank int) error {
	t2Children, others, err := bh.t1s.childrenWithParentInW(rank, bh.getTree(2).root)
	if err != nil { return err }

	if len(t2Children) > 4 {
		rmNodes := append(others, t2Children...)
		numOfRemoved := 0
		for numOfRemoved < 2 {
			err := bh.removeFromW(rmNodes[numOfRemoved])
			if err != nil {
				return err
			}
			if rmNodes[numOfRemoved].parent == bh.getTree(2).root {
				bh.insertNode(1, bh.cutOffNode(2, rmNodes[numOfRemoved]))
			} else {
				bh.rmViolatingNode(child, nil)
			}
		}
	} else {
		bh.reduceViolation(others[0], others[1])
	}
	return nil
}

func (bh *BrodalHeap) insertNewW(child *node) error {
	return bh.t1s.insertNewW(child)
}

func (bh *BrodalHeap) reduceViolation(x1 *node, x2 *node) {
	if x1.isGood() || x2.isGood() {
		bh.mbyRemoveFromViolating(x1)
		bh.mbyRemoveFromViolating(x2)
	} else {
		if x1.parent != x2.parent {
			if x1.parent.value <= x2.parent.value {
				x1.swapBrothers(x2)
			} else {
				x2.swapBrothers(x1)
			}
		}
		bh.rmViolatingNode(x1, x2)
	}
}

func (bh *BrodalHeap) rmViolatingNode(x1 *node, x2 *node) {

}
