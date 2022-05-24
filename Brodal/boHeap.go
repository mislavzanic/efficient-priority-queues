package Brodal

import (
	"container/list"
	"errors"
	"fmt"
)


type ValType float64

type BrodalHeap struct {
	t1s            *tree1Struct
	tree2          *tree
	violationNodes *list.List
}

const ALPHA int = 10

func NewEmptyHeap() *BrodalHeap {
	return &BrodalHeap{
		t1s:            newEmptyT1S(),
		tree2:          nil,
		violationNodes: list.New(),
	}
}

func NewHeap(value ValType) *BrodalHeap {
	bh := NewEmptyHeap()
	bh.t1s = newT1S(value, bh)
	return bh
}

func (bh *BrodalHeap) getTree(index int) *tree {
	if index == 1 {
		return bh.t1s.GetTree()
	}
	return bh.tree2
}

func (bh *BrodalHeap) ToString() string {
	if bh.Empty() { return "" }

	str := "Tree1:\n" + bh.getTree(1).ToString()
	if bh.getTree(2) != nil {
		str += "Tree2:\n" + bh.getTree(2).ToString()
	}

	return str
}

func (bh *BrodalHeap) Empty() bool {
	return bh.getTree(1) == nil
}

func (bh *BrodalHeap) Insert(value ValType) {
	bh.Meld(NewHeap(value))
}

func (bh *BrodalHeap) Meld(otherHeap *BrodalHeap) {
	if bh.Empty() {
		bh = otherHeap
	} else {
		minT1S, otherT1S := getMinValTree(bh.t1s, otherHeap.t1s)
		bh.t1s = minT1S
		bh.getTree(1).parentHeap = bh
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
				bh.getTree(2).parentHeap = bh
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
		}
	}

	if err := bh.updateViolations(); err != nil {
		panic(err.Error())
	}
	bh.createAlphaSpace()
	if err := bh.updateViolations(); err != nil {
		panic(err.Error())
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
	if child.rank < bh.getTree(treeIndex).RootRank()-2 {
		bh.lowRankAction(treeIndex, child, insert)
	} else {
		bh.highRankAction(treeIndex, child, insert)
	}

	bh.updateHighRank(treeIndex, bh.getTree(treeIndex).RootRank()-2)
	bh.updateHighRank(treeIndex, bh.getTree(treeIndex).RootRank()-1)
}

func (bh *BrodalHeap) lowRankAction(treeIndex int, child *node, insert bool) {
	actionsInc := bh.getTree(treeIndex).AskGuide(child.rank, bh.getTree(treeIndex).NumOfRootChildren(child.rank)+1, insert)
	// bh.handleActions(treeIndex, append(actionsInc, actionsDec...), child)
	bh.handleActions(treeIndex, actionsInc, child)
}

func (bh *BrodalHeap) highRankAction(treeIndex int, child *node, insert bool) {
	if insert {
		// r := child.rank
		if _, err := bh.getTree(treeIndex).insertNode(child); err != nil {
			panic(err.Error())
		}
		// bh.handleViolation({child: child, rank: r})
		bh.handleViolation(child)
	} else {
		if _, err := bh.getTree(treeIndex).cutOffNode(child); err != nil {
			panic(err.Error())
		}
		bh.handleViolation(child)
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
		if _, err := bh.getTree(treeIndex).insertNode(child); err != nil {
			panic(err.Error())
		}
		bh.handleViolation(child)
	case LOWER_BOUND:
		if _, err := bh.getTree(treeIndex).cutOffNode(child); err != nil {
			panic(err.Error())
		}
		bh.handleViolation(child)
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
	bh.handleViolation(minNode)
	bh.handleViolation(minNode.LeftChild())
	bh.handleViolation(minNode.LeftChild().RightBrother())
}

func (bh *BrodalHeap) delinkNodes(treeIndex int, rank int) (*node, []*node) {
	child := bh.getTree(treeIndex).RemoveChild(rank)
	bh.handleViolation(child)
	nodes, err := child.delink()
	if err != nil {
		panic(err.Error())
	}

	return child, nodes
}

func (bh *BrodalHeap) delink(treeIndex int) []*node {
	nodes := bh.getTree(treeIndex).Delink()
	for _, n := range nodes {
		bh.handleViolation(n)
	}
	return nodes
}

func (bh *BrodalHeap) updateHighRank(treeIndex int, rank int) {
	if rank < 0 {
		return
	}
	noChildren := bh.getTree(treeIndex).NumOfRootChildren(rank)
	if noChildren > 7 {
		bh.linkNodes(treeIndex, rank)
		bh.linkNodes(treeIndex, rank)
	} else if noChildren < 2 && noChildren >= 0 {
		if rank == bh.getTree(treeIndex).RootRank()-1 {
			if rank > 0 {
				nodes := bh.getTree(treeIndex).DecRank()
				for _, n := range nodes {
					childNodes := n.DecRank()
					for _, cn := range childNodes {
						bh.handleViolation(cn)
						bh.insertNode(treeIndex, cn)
					}
					bh.insertNode(treeIndex, n)
				}
			}
		} else {
			bh.updateHighRank(treeIndex, rank+1)
			if rank == bh.getTree(treeIndex).RootRank()-2 {
				nodes := bh.delink(treeIndex)
				for _, n := range nodes {
					bh.insertNode(treeIndex, n)
				}
			}
		}
	} else if noChildren < 0 {
		panic("Negative noChildren")
	}

	if rank == bh.getTree(treeIndex).RootRank()-2 {
		bh.updateHighRank(treeIndex, rank+1)
	}
}

//////////////////////////////////////////////////////////////////////////////
///////////////////////////Violation Handling/////////////////////////////////
//////////////////////////////////////////////////////////////////////////////

func (bh *BrodalHeap) mbyAddViolation(child *node) error {
	if !child.isBad() {
		return nil
	}

	if child.rank >= bh.getTree(1).RootRank() {
		err := bh.addToV(child)
		if err != nil {
			return err
		}
	} else {
		err := bh.addToW(child)
		if err != nil {
			return err
		}
	}

	return nil
}

func (bh *BrodalHeap) mbyRemoveFromViolating(notBad *node) error {
	if notBad.isBad() {
		return nil
	}

	switch notBad.parentViolatingList {
	case bh.getTree(1).wList():
		err := bh.removeFromW(notBad)
		if err != nil {
			panic(err.Error())
		}
	case nil:
		notBad.violatingSelf = nil
	default:
		// if notBad.violatingSelf == nil {
		notBad.parentViolatingList.Remove(notBad.violatingSelf)
		// }
		notBad.parentViolatingList = nil
	}

	notBad.parentViolatingList = nil
	notBad.violatingSelf = nil
	// notBad.removeSelfFromViolating()
	// if notBad.parentViolatingList == bh.getTree(1).root.wList {
	// 	err := bh.removeFromW(notBad)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

func (bh *BrodalHeap) addToV(child *node) error {
	if child.violatingSelf != nil {
		return nil
		// errMessage := "Child with rank %d, value %f is already in a violating set"
		// return errors.New(fmt.Sprintf(errMessage, child.rank, child.value))
	}
	child.violatingSelf = bh.getTree(1).vList().PushBack(child)
	child.parentViolatingList = bh.getTree(1).vList()
	return nil
}

func (bh *BrodalHeap) createAlphaSpace() {
	if bh.getTree(2) == nil || bh.getTree(1).vList().Len() <= ALPHA*bh.getTree(1).RootRank() {
		return
	}

	if bh.getTree(2).RootRank() == 0 {
		panic("Tree2 is of rank 0")
	}

	// currT1Rank := bh.getTree(1).RootRank()
	if bh.getTree(2).RootRank() <= bh.getTree(1).RootRank()+1 {
		nodes := bh.getTree(2).DecRank()
		for _, n := range nodes {
			bh.insertNode(1, n)
		}
		bh.insertNode(1, bh.getTree(2).root)
		bh.tree2 = nil
		return
	}

	newChild := bh.cutOffNode(2, bh.getTree(2).childrenRank[bh.getTree(1).RootRank()+1])
	nodes := newChild.DecRank()
	for _, n := range nodes {
		bh.insertNode(1, n)
	}

	bh.insertNode(1, newChild)
}

func (bh *BrodalHeap) addToW(child *node) error {
	// actions := bh.t1s.GetWGuide().forceIncrease(child.rank, bh.t1s.GetWNums(child.rank)+1, 2)
	actions := bh.t1s.GetWGuide().forceIncrease(child.rank, &bh.t1s.numOfNodesInT1W, 2)
	bh.handleActions(1, actions, child)
	return nil
}

func (bh *BrodalHeap) removeFromW(child *node) error {
	return bh.t1s.removeFromW(child)
}

func (bh *BrodalHeap) reduceW(rank int) error {
	t2Children, others, err := bh.t1s.childrenWithParentInW(rank, bh.getTree(2).root)
	if err != nil {
		return err
	}

	if len(t2Children) > 4 {
		rmNodes := append(others, t2Children...)[:2]
		for _, child := range rmNodes {
			err := bh.removeFromW(child)
			if err != nil {
				return err
			}
			if child.parent == bh.getTree(2).root {
				bh.insertNode(1, bh.cutOffNode(2, child))
			} else {
				err := bh.rmViolatingNode(child, nil)
				if err != nil {
					return err
				}
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

func (bh *BrodalHeap) reduceViolation(x1 *node, x2 *node) error {
	if !x1.isBad() || !x2.isBad() {
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
		err := bh.rmViolatingNode(x1, x2)
		if err != nil {
			return err
		}
	}
	return nil
}

func (bh *BrodalHeap) rmViolatingNode(x1 *node, x2 *node) error {

	if x1.parent == nil || x1.parent == bh.getTree(1).root || x1.parent == bh.getTree(1).root {
		errMessage := "Node has invalid parents, node rank: %d, node value: %f"
		return errors.New(fmt.Sprintf(errMessage, x1.rank, x1.value))
	}

	if x2 == nil {
		lb, rb := x1.LeftBrother(), x1.RightBrother()
		if rb != nil && rb.rank == x1.rank {
			x2 = rb
		} else if lb != nil && lb.rank == x1.rank {
			x2 = lb
		} else {
			errMessage := "Node doesn't have a brother of the same rank, rank: %d"
			return errors.New(fmt.Sprintf(errMessage, x1.rank))
		}
	}

	if x1.parent != x2.parent {
		errMessage := "Nodes have different parents, parent1 rank: %d, parent2 rank: %d"
		return errors.New(fmt.Sprintf(errMessage, x1.parent.rank, x2.parent.rank))
	}

	parent, grandParent := x1.parent, x1.parent.parent
	replacement := bh.getTree(1).ChildOfRank(parent.rank)

	if parent.numOfChildren[x1.rank] == 2 {
		if parent.rank == x1.rank + 1 {
			if grandParent != bh.getTree(1).root {
				bh.cutOffNode(1, replacement)
				grandParent.pushBackChild(replacement, parent)
				// ovo ispod mi je sumnjivo -> za razmisliti
				bh.mbyAddViolation(replacement)
				grandParent.removeChild(parent)
			} else {
				bh.cutOffNode(1, parent)
			}

			parent.removeChild(x1)
			parent.removeChild(x2)
			bh.insertNode(1, parent)
		} else {
			parent.removeChild(x1)
			parent.removeChild(x2)
		}
		bh.mbyRemoveFromViolating(x2)
		bh.insertNode(1, x1)
		bh.insertNode(1, x2)
	} else {
		parent.removeChild(x1)
		bh.insertNode(1, x1)
	}

	bh.mbyRemoveFromViolating(x1)
	return nil
}

func (bh *BrodalHeap) handleViolation(child *node) {
	bh.violationNodes.PushBack(child)
}

func (bh *BrodalHeap) updateViolations() error {
	for violation := bh.violationNodes.Front(); violation != nil; violation = violation.Next() {
		err := bh.updateViolation(violation.Value.(*node))
		if err != nil {
			return err
		}
	}
	bh.violationNodes = bh.violationNodes.Init()
	return nil
}

func (bh *BrodalHeap) updateViolation(violation *node) error {
	if !violation.isBad() {
		if violation.violatingSelf != nil {
			bh.mbyRemoveFromViolating(violation)
			// violation.removeSelfFromViolating()
		}
		violation.parentViolatingList = nil
	} else {
		err := bh.mbyAddViolation(violation)
		if err != nil {
			return err
		}
	}
	return nil
}
