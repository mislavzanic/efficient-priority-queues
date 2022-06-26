package Brodal

import (
	"container/list"
	"errors"
	"fmt"
	"math"
)

type Number interface {
	int | int8 | int16 | int32 | int64 | float32 | float64
}

type BrodalHeap[T Number] struct {
	t1s            *tree1Struct[T]
	tree2          *tree[T]
	violationNodes *list.List
}

const ALPHA int = 10

func NewEmptyHeap[T Number]() *BrodalHeap[T] {
	return &BrodalHeap[T]{
		t1s:            newEmptyT1S[T](),
		tree2:          nil,
		violationNodes: list.New(),
	}
}

func NewHeap[T Number](value T, ident any) *BrodalHeap[T] {
	bh := NewEmptyHeap[T]()
	bh.t1s = newT1S(value, ident, bh)
	return bh
}

func (bh *BrodalHeap[T]) getTree(index int) *tree[T] {
	if index == 1 {
		return bh.t1s.GetTree()
	}
	return bh.tree2
}

func (bh *BrodalHeap[T]) Empty() bool {
	return bh.getTree(1) == nil
}

func (bh *BrodalHeap[T]) Min() *node[T] {
	return bh.getTree(1).root
}

//////////////////////////////////////////////////////////////////////////////
/////////////////////////////////DeleteMin////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////

func (bh *BrodalHeap[T]) DeleteMin() *node[T] {
	min := bh.Min()

	if bh.getTree(2) != nil {
		if bh.getTree(2).RootRank() == 0 && bh.getTree(1).RootRank() == 0 {
			bh.t1s = newT1S(bh.getTree(2).RootValue(), bh.getTree(2).root.Ident(), bh)
			bh.tree2 = nil
			return min
		} else {
			bh.moveT2ToT1()
		}
	} else {
		if bh.getTree(1).RootRank() == 0 {
			bh.t1s = newEmptyT1S[T]()
			return min
		}
	}


	newMin, other := bh.getNewMin()
	trees := bh.getTree(1).RmRfRoot()
	newMinV, newMinW := newMin.vList, newMin.wList
	oldV, oldW := bh.getTree(1).vList(), bh.getTree(1).wList()

	bh.t1s = newT1S(newMin.value, newMin.ident, bh)
	bh.tree2 = nil

	for e := trees.Back(); e != nil; e = e.Prev() {
		if e.Value.(*node[T]) == newMin {
			continue
		}
		for e.Value.(*node[T]).rank > bh.getTree(1).RootRank() {
			nodes := e.Value.(*node[T]).DecRank()
			for _, n := range nodes {
				bh.insertNode(1, n)
			}
		}
		bh.insertNode(1, e.Value.(*node[T]))
	}


	for e := newMin.children.Front(); e != nil; {
		newE := e.Next()
		bh.insertNode(1, e.Value.(*node[T]))
		e = newE
	}


	for e := oldV.Front(); e != nil; e = e.Next() {
		bh.handleViolation(e.Value.(*node[T]))
	}

	for e := oldW.Front(); e != nil; e = e.Next() {
		bh.handleViolation(e.Value.(*node[T]))
	}

	for e := newMinV.Front(); e != nil; e = e.Next() {
		bh.handleViolation(e.Value.(*node[T]))
	}

	for e := newMinW.Front(); e != nil; e = e.Next() {
		bh.handleViolation(e.Value.(*node[T]))
	}

	bh.handleViolation(other)
	err := bh.updateViolations()
	if err != nil {
		panic(err.Error())
	}

	for e := bh.getTree(1).wList().Front(); e != nil; {
		if e.Next().Value.(*node[T]).rank == e.Value.(*node[T]).rank {
			bh.reduceViolation(e.Value.(*node[T]), e.Next().Value.(*node[T]))
		} else {
			e = e.Next()
		}
	}

	return min
}

func (bh *BrodalHeap[T]) DecreaseKey(child *node[T], newValue T) {
	if newValue < bh.Min().Value() {
		child.value = bh.Min().Value()
		bh.getTree(1).SetRootValue(newValue)
	} else {
		child.SetValue(newValue)
	}
	bh.handleViolation(child)
}

func (bh *BrodalHeap[T]) Delete(child *node[T]) {
	bh.DecreaseKey(child, T(math.Inf(-1)))
	bh.DeleteMin()
}

func (bh *BrodalHeap[T]) Insert(value T, ident any) {
	bh.Meld(NewHeap(value, ident))
}

//////////////////////////////////////////////////////////////////////////////
///////////////////////////////// Meld ///////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////

func (bh *BrodalHeap[T]) Meld(pqInterface interface{}) {
	otherHeap := pqInterface.(*BrodalHeap[T])
	if bh.Empty() {
		bh.t1s = otherHeap.t1s
		bh.tree2 = otherHeap.tree2
		bh.getTree(1).parentHeap = bh
		if bh.getTree(2) != nil {
			bh.getTree(2).parentHeap = bh
		}
	} else {
		minT1S, otherT1S := getMinValTree(bh.t1s, otherHeap.t1s)
		bh.t1s = minT1S
		bh.getTree(1).parentHeap = bh
		if bh.getTree(2) == nil && otherHeap.getTree(2) == nil {
			if minT1S.GetTree().RootRank() <= otherT1S.GetTree().RootRank() {
				bh.tree2 = otherT1S.GetTree()
				bh.getTree(2).parentHeap = bh
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

func (bh *BrodalHeap[T]) insertNode(treeIndex int, child *node[T]) {
	bh.heapAction(treeIndex, child, true)
}

func (bh *BrodalHeap[T]) cutOffNode(treeIndex int, child *node[T]) *node[T] {
	bh.heapAction(treeIndex, child, false)
	return child
}

func (bh *BrodalHeap[T]) heapAction(treeIndex int, child *node[T], insert bool) {
	if child.rank < bh.getTree(treeIndex).RootRank()-2 {
		bh.lowRankAction(treeIndex, child, insert)
	} else {
		bh.highRankAction(treeIndex, child, insert)
	}

	bh.updateHighRank(treeIndex, bh.getTree(treeIndex).RootRank()-2)
	bh.updateHighRank(treeIndex, bh.getTree(treeIndex).RootRank()-1)
}

func (bh *BrodalHeap[T]) lowRankAction(treeIndex int, child *node[T], insert bool) {
	actionsInc := bh.getTree(treeIndex).AskGuide(child.rank, bh.getTree(treeIndex).NumOfRootChildren(child.rank)+1, insert)
	// bh.handleActions(treeIndex, append(actionsInc, actionsDec...), child)
	bh.handleActions(treeIndex, actionsInc, child)
}

func (bh *BrodalHeap[T]) highRankAction(treeIndex int, child *node[T], insert bool) {
	if insert {
		if _, err := bh.getTree(treeIndex).insertNode(child); err != nil {
			panic(err.Error())
		}
	} else {
		if _, err := bh.getTree(treeIndex).cutOffNode(child); err != nil {
			panic(err.Error())
		}
	}
	bh.handleViolation(child)
}

func (bh *BrodalHeap[T]) handleActions(treeIndex int, actions []action, child *node[T]) {
	for _, act := range actions {
		switch op := act.op; op {
		case Reduce:
			bh.doReduceAction(treeIndex, act.bound, act.index)
		default:
			bh.doDefaultAction(treeIndex, act.bound, child)
		}
	}
}

func (bh *BrodalHeap[T]) doReduceAction(treeIndex int, guideBound int, rank int) {
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

func (bh *BrodalHeap[T]) doDefaultAction(treeIndex int, guideBound int, child *node[T]) {
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

func (bh *BrodalHeap[T]) linkNodes(treeIndex int, rank int) {
	minNode := bh.getTree(treeIndex).Link(rank)
	bh.handleViolation(minNode)
	bh.handleViolation(minNode.LeftChild())
	bh.handleViolation(minNode.LeftChild().RightBrother())
}

func (bh *BrodalHeap[T]) delinkNodes(treeIndex int, rank int) (*node[T], []*node[T]) {
	child := bh.getTree(treeIndex).RemoveChild(rank)
	bh.handleViolation(child)
	nodes, err := child.delink()
	if err != nil {
		panic(err.Error())
	}

	return child, nodes
}

func (bh *BrodalHeap[T]) delink(treeIndex int) []*node [T]{
	nodes := bh.getTree(treeIndex).Delink()
	for _, n := range nodes {
		bh.handleViolation(n)
	}
	return nodes
}

func (bh *BrodalHeap[T]) updateHighRank(treeIndex int, rank int) {
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

func (bh *BrodalHeap[T]) mbyAddViolation(child *node[T]) error {
	if !child.isBad() {
		return nil
	}

	if child.rank >= bh.getTree(1).RootRank() {
		err := bh.addToV(child)
		if err != nil {
			return err
		}
	} else {
		if child.parentViolatingList != nil && child.violatingSelf != nil {
			if child.parentViolatingList == bh.getTree(1).wList() {
				return nil
			}
		}
		err := bh.addToW(child)
		if err != nil {
			return err
		}
	}

	return nil
}

func (bh *BrodalHeap[T]) mbyRemoveFromViolating(notBad *node[T]) error {
	if notBad.isBad() {
		return nil
	}

	if notBad.violatingSelf == nil && notBad.parentViolatingList != nil {
		panic("Violating self is nil")
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
		notBad.parentViolatingList.Remove(notBad.violatingSelf)
		notBad.parentViolatingList = nil
	}

	notBad.parentViolatingList =nil
	notBad.violatingSelf = nil

	return nil
}

func (bh *BrodalHeap[T]) addToV(child *node[T]) error {
	if child.violatingSelf != nil {
		switch child.parentViolatingList {
		case bh.getTree(1).wList():
			bh.removeFromW(child)
		case nil:
			child.violatingSelf = nil
		default:
			child.parentViolatingList.Remove(child.violatingSelf)
		}
		child.removeSelfFromViolating()
	}

	child.violatingSelf = bh.getTree(1).vList().PushBack(child)
	child.parentViolatingList = bh.getTree(1).vList()
	return nil
}

func (bh *BrodalHeap[T]) createAlphaSpace() {
	if bh.getTree(2) == nil || bh.getTree(1).vList().Len() <= ALPHA*bh.getTree(1).RootRank() {
		return
	}

	if bh.getTree(2).RootRank() == 0 {
		panic("Tree2 is of rank 0")
	}

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

func (bh *BrodalHeap[T]) addToW(child *node[T]) error {
	actions := bh.t1s.GetWGuide().forceIncrease(child.rank, &bh.t1s.numOfNodesInT1W, 2)
	bh.handleActions(1, actions, child)
	return nil
}

func (bh *BrodalHeap[T]) removeFromW(child *node[T]) error {
	err := bh.t1s.removeFromW(child)
	return err
}

func (bh *BrodalHeap[T]) reduceW(rank int) error {
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
		// bh.reduceViolation(others[0], others[1])
		bh.rmViolatingNode(others[0], nil)
		bh.rmViolatingNode(others[1], nil)
	}
	return nil
}

func (bh *BrodalHeap[T]) insertNewW(child *node[T]) error {
	return bh.t1s.insertNewW(child)
}

func (bh *BrodalHeap[T]) reduceViolation(x1 *node[T], x2 *node[T]) error {
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

func (bh *BrodalHeap[T]) rmViolatingNode(x1 *node[T], x2 *node[T]) error {

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

func (bh *BrodalHeap[T]) handleViolation(child *node[T]) {
	bh.violationNodes.PushBack(child)
}

func (bh *BrodalHeap[T]) updateViolations() error {
	for violation := bh.violationNodes.Front(); violation != nil; violation = violation.Next() {
		err := bh.updateViolation(violation.Value.(*node[T]))
		if err != nil {
			return err
		}
	}
	bh.violationNodes = bh.violationNodes.Init()
	return nil
}

func (bh *BrodalHeap[T]) updateViolation(violation *node[T]) error {
	if !violation.isBad() {
		if violation.violatingSelf != nil {
			err := bh.mbyRemoveFromViolating(violation)
			if err != nil {
				panic(err.Error())
			}
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

func (bh *BrodalHeap[T]) moveT2ToT1() {
	child := bh.getTree(2).Children().Back()
	for bh.getTree(2).Children().Len() > 0 {
		_, err := bh.getTree(2).cutOffNode(child.Value.(*node[T]))
		if err != nil {
			panic(err.Error())
		}
		bh.insertNode(1, child.Value.(*node[T]))
		child = bh.getTree(2).Children().Back()
	}
	bh.insertNode(1, bh.getTree(2).root)
	bh.tree2 = nil
}

func (bh *BrodalHeap[T]) getNewMin() (*node[T], *node[T]) {
	minW := bh.getTree(1).root.getMinFromW()
	minV := bh.getTree(1).root.getMinFromV()
	minSubTree := bh.getTree(1).root.getMinFromChildren()
	newMin, _, _ := getMinNodeFrom3(minW, minV, minSubTree)

	mbySwap := bh.t1s.tree1.childrenRank[newMin.rank]

	if newMin.parent != bh.getTree(1).root {
		mbySwap.parent = newMin.parent
		newMin.parent = bh.getTree(1).root

		temp := mbySwap.self
		mbySwap.self = newMin.self
		newMin.self = temp
		bh.getTree(1).childrenRank[newMin.rank] = newMin
	}

	return newMin, mbySwap
}
