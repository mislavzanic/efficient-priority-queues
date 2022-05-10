package Brodal

import "math"

type BrodalHeap struct {
	tree1 *tree
	tree2 *tree

	numOfNodesInT1W []uint
	rankPointersT1W []*node
	t1GuideW        *guide
}

type violationSetType byte

const (
	wSet  violationSetType = 0
	vSet                   = 1
	notViolation           = 2
)

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
		if !currKey.isGood() {
			bh.updateViolatingSet(bh.addViolation(currKey), currKey)
			// TODO take care of a violation
		}
	}
}

func (bh *BrodalHeap) Delete(node *node) {
	bh.DecreaseKey(node, math.Inf(-1))
	bh.DeleteMin()
}

func (bh *BrodalHeap) DeleteMin() {
	child := bh.tree2.children().Front()

	for bh.tree2.children().Len() != 0 {
		bh.tree2.removeRootChild(child.Value.(*node))
		bh.tree1.Insert(child.Value.(*node))
		child = bh.tree2.children().Front()
	}

	bh.tree1.Insert(bh.tree2.root)
	bh.tree2 = nil

	newMin := bh.tree1.childrenRank[bh.tree1.rootRank()]
	//TODO
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

		for _, tree := range others {
			if tree != minTree && tree != maxTree {

				if maxTree.root.rank == tree.root.rank && maxTree != minTree {
					nodes, _ := tree.delinkFromRoot()
					for _, n := range nodes {
						bh.insertNodeIntoTree(maxTree, n)
					}
				}

				bh.insertNodeIntoTree(maxTree, tree.root)
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
	if node.parentViolatingList != nil { bh.removeFromViolating(node) }
}

func (bh *BrodalHeap) removeFromViolating(node *node) {
	node.removeSelfFromViolating()
}

func (bh *BrodalHeap) mbyAddViolation(mbyBad *node) violationSetType {
	if !mbyBad.isGood() { return bh.addViolation(mbyBad) }
	return notViolation
}

func (bh *BrodalHeap) addViolation(bad *node) violationSetType {
	if bad.rank > bh.tree1.RootRank() {
		bad.violatingSelf = bh.tree1.root.vList.PushFront(bad)
		return vSet
	}
	// else

	if bh.t1GuideW.boundArray[bad.rank].fst != 0 {
		bad.violatingSelf = bh.tree1.root.wList.InsertAfter(bad, bh.rankPointersT1W[bad.rank].violatingSelf)
	} else {
		bad.violatingSelf = bh.tree1.root.wList.PushFront(bad)
		bh.rankPointersT1W[bad.rank] = bad
	}

	bh.numOfNodesInT1W[bad.rank]++
	return wSet
}

func (bh *BrodalHeap) updateViolatingSet(setType violationSetType, node *node) {
	if setType == notViolation {
		panic("Wrong set type")
	}

	if setType == wSet {
		bh.updateWSet(node)
	} else {
		// TODO update vset
		bh.updateVSet(node)
	}
}

func (bh *BrodalHeap) updateWSet(bad *node) {

	acts := bh.t1GuideW.forceIncrease(int(bad.rank), int(bh.numOfNodesInT1W[bad.rank]), 2)

	for _, act := range acts {
		bh.performeAct(act)
	}
}

func (bh *BrodalHeap) performeAct(act action) {
	numOfSonsOfT2 := 0
	notSonsOfT2 := []*node{}

	for e := bh.rankPointersT1W[act.index].violatingSelf; e.Value.(*node).rank != uint(act.index); e = e.Next() {
		if e.Value.(*node).parent == bh.tree2.root {
			numOfSonsOfT2++
		} else {
			notSonsOfT2 = append(notSonsOfT2, e.Value.(*node))
		}
	}

	if numOfSonsOfT2 > 4 {
		numOfRemoved := 0
		for _, rmNode := range notSonsOfT2 {

			bh.tree1.removeViolatingNode(rmNode, nil)
			numOfRemoved++
		}

		for e := bh.rankPointersT1W[act.index].violatingSelf; e.Value.(*node).rank != uint(act.index) && numOfRemoved < 2; {
			bh.rankPointersT1W[act.index] = e.Next().Value.(*node)

			e.Value.(*node).removeSelfFromViolating()
			bh.tree2.cut(e.Value.(*node))

			bh.tree1.Insert(e.Value.(*node))
			numOfRemoved++
		}
	} else {
		// tree reduce violations
		bh.tree1.reduceViolaton(notSonsOfT2[0], notSonsOfT2[1])

		notGood := func () *node {
			if !notSonsOfT2[0].isGood() {
				return notSonsOfT2[0]
			} else if !notSonsOfT2[1].isGood() {
				return notSonsOfT2[1]
			} else { return nil }
		}()

		if  notGood != nil {
			bh.tree1.removeViolatingNode(notGood, nil)
		}
	}

	bh.numOfNodesInT1W[act.index] -= 2
}

func (bh *BrodalHeap) updateVSet(bad *node) {

}

func (bh *BrodalHeap) insertNodeIntoTree(tree *tree, node *node) {
	tree.addRootChild(node)
	bh.mbyAddViolation(node)

	bh.updateLowRank(tree, node, true)

	bh.updateHighRank(tree, tree.root.rank-2)
	bh.updateHighRank(tree, tree.root.rank-1)
}

func (bh *BrodalHeap) cutNodeFromTree(tree *tree, node *node) {
	bh.mbyRemoveFromViolating(node)
	tree.removeRootChild(node)

	bh.updateLowRank(tree, node, false)

	bh.updateHighRank(tree, tree.root.rank-2)
	bh.updateHighRank(tree, tree.root.rank-1)
}

func (bh *BrodalHeap) updateLowRank(tree *tree, node *node, insert bool) {
	if node.rank < tree.root.rank - 2 {
		response := tree.askGuide(int(node.rank), int(tree.root.numOfChildren[node.rank]), insert)

		for _, act := range response {
			if insert {
				tree.link(uint(act.index))
			} else {
				removeThis := tree.childrenRank[act.index]
				tree.removeRootChild(removeThis)

				nodes := removeThis.delink()
				for _, n := range nodes {
					bh.insertNodeIntoTree(tree, n)
				}
				bh.insertNodeIntoTree(tree, removeThis)
			}
		}
	}
}

func (bh *BrodalHeap) updateHighRank(tree *tree, rank uint) {
	if tree.root.numOfChildren[rank] > 7 {
		nodeSliceX := tree.root.delink()
		nodeSliceY := tree.root.delink()
		nodeSliceZ := tree.root.delink()

		minNode1, nodeX1, nodeY1 := getMinNode(nodeSliceX[0], nodeSliceX[1], nodeSliceY[0])
		minNode2, nodeX2, nodeY2 := getMinNode(nodeSliceY[1], nodeSliceZ[1], nodeSliceZ[0])

		minNode1.link(nodeX1, nodeY1)
		minNode2.link(nodeX2, nodeY2)

		if rank == tree.root.rank-1 {
			bh.incTreeRank(tree, minNode1, minNode2)
		} else {
			bh.insertNodeIntoTree(tree, nodeSliceX[0])
			bh.insertNodeIntoTree(tree, nodeSliceY[1])
		}
	} else if tree.root.numOfChildren[rank] < 2 {
		//...
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
		if x1.isGood() {
			x1.mbyRemoveSelfFromViolating()
			bh.numOfNodesInT1W[x1.rank]--
		}
		if x2.isGood() {
			x2.mbyRemoveSelfFromViolating()
			bh.numOfNodesInT1W[x2.rank]--
		}
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
				grandParent.addChild(replacement, parent)
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
		otherBrother.removeSelfFromViolating()

		bh.insertNodeIntoTree(bh.tree1, otherBrother)
		bh.insertNodeIntoTree(bh.tree1, rmNode)
	} else {
		parent.removeChild(rmNode)
		bh.insertNodeIntoTree(bh.tree1, rmNode)
	}
	rmNode.removeSelfFromViolating()
}
