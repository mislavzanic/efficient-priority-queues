package Brodal

import "math"

type BrodalHeap struct {
	tree1 *tree
	tree2 *tree

	numOfNodesInT1W []int
	rankPointersT1W []*node
	t1GuideW        *guide
}

const ALPHA int = 10

func NewHeap(value float64) *BrodalHeap {
	return &BrodalHeap{
		tree1: newTree(value, 1),
		tree2: nil,
	}
}

func (bh *BrodalHeap) Min() float64 {
	return bh.tree1.root.value
}

func (bh *BrodalHeap) DecreaseKey(currKey *node, value float64) {
	if value < bh.tree1.root.value {
		currKey.value = bh.tree1.root.value
		bh.tree1.root.value = currKey.value
	} else {
		currKey.value = value
		bh.mbyAddViolation(currKey)
	}
}

func (bh *BrodalHeap) Delete(node *node) {
	bh.DecreaseKey(node, math.Inf(-1))
	bh.DeleteMin()
}

func (bh *BrodalHeap) DeleteMin() {
	child := bh.tree2.Children().Front()

	for bh.tree2.Children().Len() != 0 {
		bh.tree2.removeRootChild(child.Value.(*node))
		bh.insertNodeIntoTree(bh.tree1, child.Value.(*node))
		child = bh.tree2.Children().Front()
	}

	bh.insertNodeIntoTree(bh.tree1, bh.tree2.root)
	bh.tree2 = nil

	minW := bh.tree1.root.getMinFromW()
	minV := bh.tree1.root.getMinFromV()
	minTree := bh.tree1.root.getMinFromChildren()
	newMin, _, _ := getMinNodeFrom3(minW, minV, minTree)

	mbySwap := bh.tree1.childrenRank[newMin.rank]
	indepTrees := bh.tree1.rmRfRoot()

	if newMin.parent != nil {
		mbySwap.parent = newMin.parent
		newMin.parent = nil

		indepTrees.Remove(mbySwap.self)
		mbySwap.self = newMin.self
		newMin.self = nil
	}

	oldV, oldW := bh.tree1.root.vList, bh.tree1.root.wList

	bh.tree1 = newTree(newMin.value, 1)

	for e := indepTrees.Back(); e != nil; e = e.Prev() {
		if e.Value.(*node).rank == 0 {
			bh.Insert(e.Value.(*node).value)
		} else {
			bh.insertNodeIntoTree(bh.tree1, e.Value.(*node))
		}
	}

	for e := newMin.children.Back(); e != nil; e = e.Prev() {
		bh.insertNodeIntoTree(bh.tree1, e.Value.(*node))
	}

	for e := oldV.Front(); e != nil; e = e.Next() {
		bh.mbyAddViolation(e.Value.(*node))
	}

	for e := oldW.Front(); e != nil; e = e.Next() {
		bh.mbyAddViolation(e.Value.(*node))
	}

	bh.mbyAddViolation(mbySwap)

	for e := bh.tree1.root.wList.Front(); e != nil; {
		if e.Next().Value.(*node).rank == e.Value.(*node).rank {
			bh.reduceViolation(e.Value.(*node), e.Next().Value.(*node))
		} else {
			e = e.Next()
		}
	}

	// nekak napraviti update guide-a za W
	// provjeriti jeli ovo bitno, il se radi vec gore automatski
}

func (bh *BrodalHeap) Insert(value float64) {
	otherHeap := NewHeap(value)
	bh.Meld(otherHeap)
}

func (bh *BrodalHeap) Meld(other *BrodalHeap) {
	trees := [](*tree){bh.tree1, bh.tree2, other.tree1, other.tree2}

	if bh.tree1.root.rank == 0 && bh.tree2 == nil {
		if other.tree1.root.rank == 0 && other.tree2 == nil {
			bh.tree1, other.tree1 = mbySwapTree(bh.tree1, other.tree1, bh.Min() < other.Min())
			bh.tree2 = other.tree2
		}
	} else {
		minTree, _ := getMinTree(bh.tree1, other.tree1)
		maxTree, others := getMaxTree(trees)

		if minTree.RootRank() == maxTree.RootRank() {
			maxTree = minTree
		}

		if maxTree.RootRank() == 0 {
			maxTree.linkToRoot(others[0].root, others[1].root)
			if len(others) == 3 {
				bh.insertNodeIntoTree(maxTree, others[2].root)
			}
		} else {
			for _, tree := range others {
				if tree != minTree {
					for maxTree.root.rank == tree.root.rank && maxTree != minTree {
						nodes := tree.delink()
						for _, n := range nodes {
							bh.insertNodeIntoTree(maxTree, n)
						}
					}
					bh.insertNodeIntoTree(maxTree, tree.root)
				}
			}
		}

		bh.tree1 = minTree
		if maxTree != minTree {
			bh.tree2 = maxTree
		} else {
			bh.tree2 = nil
		}
	}
}

func (bh *BrodalHeap) mbyRemoveFromViolating(node *node) {
	if node.parentViolatingList != nil && node.isGood() {
		bh.removeFromViolating(node)
	}
}

func (bh *BrodalHeap) removeFromViolating(notBad *node) {
	if notBad.parentViolatingList == bh.tree1.root.wList {
		bh.numOfNodesInT1W[notBad.rank]--
		if bh.rankPointersT1W[notBad.rank] == notBad {
			bh.rankPointersT1W[notBad.rank] = notBad.violatingSelf.Next().Value.(*node)
		}
	}
	notBad.removeSelfFromViolating()
}

func (bh *BrodalHeap) mbyAddViolation(mbyBad *node) {
	if !mbyBad.isGood() && mbyBad.violatingSelf == nil {
		bh.addViolation(mbyBad)
	}
}

func (bh *BrodalHeap) addViolation(bad *node) {
	if bad.rank > bh.tree1.RootRank() {
		bad.violatingSelf = bh.tree1.root.vList.PushFront(bad)
		bh.updateVSet(bad)
	} else {
		bh.updateWSet(bad)
	}
}

func (bh *BrodalHeap) updateWSet(bad *node) {

	acts := bh.t1GuideW.forceIncrease(bad.rank, bh.numOfNodesInT1W[bad.rank] + 1, 2)

	for _, act := range acts {
		if act.op == Increase {
			if bh.t1GuideW.boundArray[bad.rank].fst != 0 {
				bad.violatingSelf = bh.tree1.root.wList.InsertAfter(bad, bh.rankPointersT1W[bad.rank].violatingSelf)
			} else {
				bad.violatingSelf = bh.tree1.root.wList.PushFront(bad)
				bh.rankPointersT1W[bad.rank] = bad
			}

			bh.numOfNodesInT1W[bad.rank]++
		} else {
			bh.reduceWViolations(act)
		}
	}
}

func (bh *BrodalHeap) updateVSet(bad *node) {
	if bh.tree1.root.vList.Len() > ALPHA * bh.tree1.RootRank() {
		if bh.tree2 == nil {
			panic("This can't happen")
		}
		if bh.tree2.RootRank() <= bh.tree1.RootRank() + 1 {
			for leftmostSon := bh.tree2.LeftmostSon(); leftmostSon.rank < bh.tree1.RootRank(); {
				bh.cutNodeFromTree(bh.tree2, leftmostSon)
				bh.insertNodeIntoTree(bh.tree1, leftmostSon)
			}
			bh.insertNodeIntoTree(bh.tree1, bh.tree2.root)
			bh.tree2 = nil
		} else {
			bh.derankAndInsertRootChild(bh.tree2, bh.tree1, bh.tree2.childrenRank[bh.tree1.RootRank() + 1])
		}
	}
}

func (bh *BrodalHeap) reduceWViolations(act action) {
	numOfSonsOfT2 := 0
	notSonsOfT2 := []*node{}

	for e := bh.rankPointersT1W[act.index].violatingSelf; e.Value.(*node).rank != act.index; e = e.Next() {
		if e.Value.(*node).parent == bh.tree2.root {
			numOfSonsOfT2++
		} else {
			notSonsOfT2 = append(notSonsOfT2, e.Value.(*node))
		}
	}

	if numOfSonsOfT2 > 4 {
		numOfRemoved := 0
		for _, rmNode := range notSonsOfT2 {

			bh.removeViolatingNode(rmNode, nil)
			numOfRemoved++
		}

		for e := bh.rankPointersT1W[act.index].violatingSelf; e.Value.(*node).rank != act.index && numOfRemoved < 2; {
			bh.rankPointersT1W[act.index] = e.Next().Value.(*node)

			e.Value.(*node).removeSelfFromViolating()
			bh.cutNodeFromTree(bh.tree2, e.Value.(*node))

			bh.insertNodeIntoTree(bh.tree1, e.Value.(*node))
			numOfRemoved++
		}
	} else {
		bh.reduceViolation(notSonsOfT2[0], notSonsOfT2[1])

		notGood := func () *node {
			if !notSonsOfT2[0].isGood() {
				return notSonsOfT2[0]
			} else if !notSonsOfT2[1].isGood() {
				return notSonsOfT2[1]
			} else { return nil }
		}()

		if  notGood != nil {
			bh.removeViolatingNode(notGood, nil)
		}
	}

	bh.numOfNodesInT1W[act.index] -= 2
}

func (bh *BrodalHeap) insertNodeIntoTree(tree *tree, node *node) {
	if node.rank < tree.RootRank() - 2 {
		bh.updateLowRank(tree, node, true)
	} else {
		tree.addRootChild(node)
		bh.mbyAddViolation(node)
	}
	bh.updateHighRank(tree, tree.root.rank-2)
	bh.updateHighRank(tree, tree.root.rank-1)
}

func (bh *BrodalHeap) cutNodeFromTree(tree *tree, node *node) {
	if node.rank < tree.RootRank() - 2 {
		bh.updateLowRank(tree, node, false)
	} else {
		bh.mbyRemoveFromViolating(node)
		tree.removeRootChild(node)
	}
	bh.updateHighRank(tree, tree.root.rank-2)
	bh.updateHighRank(tree, tree.root.rank-1)
}

func (bh *BrodalHeap) derankAndInsertRootChild(removeRoot *tree, insertRoot *tree, child *node) {
	rank := child.rank
	bh.mbyRemoveFromViolating(child)
	removeRoot.removeRootChild(child)

	for child.rank >= rank {
		nodes := child.delink()
		for _, n := range nodes {
			bh.insertNodeIntoTree(insertRoot, n)
		}
	}
	bh.insertNodeIntoTree(insertRoot, child)
}

func (bh *BrodalHeap) updateLowRank(tree *tree, node *node, insert bool) {
	response := tree.askGuide(node.rank, tree.root.numOfChildren[node.rank], insert)

	for _, act := range response {
		if insert {
			if act.op == Increase {
				tree.addRootChild(node)
				bh.mbyAddViolation(node)
			} else {
				tree.link(act.index)
			}
		} else {
			if act.op == Increase {
				bh.mbyRemoveFromViolating(node)
				tree.removeRootChild(node)
			} else {
				bh.derankAndInsertRootChild(tree, tree, tree.childrenRank[act.index + 1])
			}
		}
	}
}

func (bh *BrodalHeap) updateHighRank(tree *tree, rank int) {
	if tree.root.numOfChildren[rank] > 7 {
		if rank == tree.root.rank-1 {
			nodeSliceX := tree.root.delink()
			nodeSliceY := tree.root.delink()
			nodeSliceZ := tree.root.delink()

			minNode1, nodeX1, nodeY1 := getMinNodeFrom3(nodeSliceX[0], nodeSliceX[1], nodeSliceY[0])
			minNode2, nodeX2, nodeY2 := getMinNodeFrom3(nodeSliceY[1], nodeSliceZ[1], nodeSliceZ[0])

			minNode1.link(nodeX1, nodeY1)
			minNode2.link(nodeX2, nodeY2)

			bh.incTreeRank(tree, minNode1, minNode2)
		} else {
			tree1 := tree.removeRootChild(tree.childrenRank[rank])
			tree2 := tree.removeRootChild(tree.childrenRank[rank])
			tree3 := tree.removeRootChild(tree.childrenRank[rank])

			minTree, tree1, tree2 := getMinNodeFrom3(tree1, tree2, tree3)
			minTree.link(tree1, tree2)
			bh.insertNodeIntoTree(tree, minTree)
			bh.updateHighRank(tree, rank + 1)
		}
	} else if tree.root.numOfChildren[rank] < 2 {
		if rank == tree.RootRank() - 2 {
			takeChildrenFrom := tree.childrenRank[rank + 1]

			if takeChildrenFrom.numOfChildren[rank] <= 3 {
				bh.mbyRemoveFromViolating(takeChildrenFrom)
				tree.removeRootChild(takeChildrenFrom)
			}

			nodes := takeChildrenFrom.delink()
			for _, n := range nodes {
				bh.insertNodeIntoTree(tree, n)
			}

			bh.updateHighRank(tree, rank + 1)

		} else {
			bh.derankAndInsertRootChild(tree, tree, tree.childrenRank[rank])
		}
	}
}

func (bh *BrodalHeap) incTreeRank(tree *tree, node1 *node, node2 *node) {

	if tree.id == 1 {
		bh.rankPointersT1W = append(bh.rankPointersT1W, nil)
		bh.numOfNodesInT1W = append(bh.numOfNodesInT1W, 0)
	}

	tree.incRank(node1, node2)

	bh.mbyAddViolation(node1)
	bh.mbyAddViolation(node2)
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
		bh.removeViolatingNode(x1, x2)
	}
}

func (bh *BrodalHeap) removeViolatingNode(rmNode *node, otherBrother *node) {
	parent := rmNode.parent
	replacement := bh.tree1.childrenRank[parent.rank]
	grandParent := parent.parent
	if otherBrother == nil {
		otherBrother = func () *node {
			if rmNode.leftBrother().rank != rmNode.rank {
				return rmNode.rightBrother()
			} else {
				return rmNode.leftBrother()
			}
		}()
	}

	if parent.numOfChildren[rmNode.rank] == 2 {
		if parent.rank == rmNode.rank+1 {
			if grandParent != bh.tree1.root {
				bh.cutNodeFromTree(bh.tree1, replacement)
				grandParent.pushBackChild(replacement, parent)
				bh.mbyAddViolation(replacement)
				grandParent.removeChild(parent)
			} else {
				bh.cutNodeFromTree(bh.tree1, parent)
			}

			parent.removeChild(rmNode)
			parent.removeChild(otherBrother)
			bh.insertNodeIntoTree(bh.tree1, parent)
		} else {
			parent.removeChild(rmNode)
			parent.removeChild(otherBrother)
		}
		bh.mbyRemoveFromViolating(otherBrother)

		bh.insertNodeIntoTree(bh.tree1, otherBrother)
		bh.insertNodeIntoTree(bh.tree1, rmNode)
	} else {
		parent.removeChild(rmNode)
		bh.insertNodeIntoTree(bh.tree1, rmNode)
	}
	bh.mbyRemoveFromViolating(rmNode)
}
