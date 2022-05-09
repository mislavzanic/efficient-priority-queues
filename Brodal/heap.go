package Brodal

import "math"

type BrodalHeap struct {
	tree1 *tree
	tree2 *tree
}

type violationSetType byte

const (
	wSet  violationSetType = 0
	vSet                   = 1
	error                  = 2
)

const ALPHA int = 10

func newHeap(value float64) *BrodalHeap {
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
		bh.tree1.insert(child.Value.(*node))
		child = bh.tree2.children().Front()
	}

	bh.tree1.insert(bh.tree2.root)
	bh.tree2 = nil

	newMin := bh.tree1.childrenRank[bh.tree1.rootRank()]
}

func (bh *BrodalHeap) Insert(value float64) {
	otherHeap := newHeap(value)
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
						maxTree.insert(n)
					}
				}

				maxTree.insert(tree.root)
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

func (bh *BrodalHeap) addViolation(bad *node) violationSetType {
	return bh.tree1.addViolation(bad)
}

func (bh *BrodalHeap) updateViolatingSet(setType violationSetType, node *node) {
	if setType == error {
		panic("Wrong set type")
	}

	if setType == wSet {
		bh.updateWSet(node)
	} else {

	}
}

func (bh *BrodalHeap) updateWSet(bad *node) {

	actions := bh.tree1.mainTreeGuideW.forceIncrease(int(bad.rank), int(bh.tree1.numOfNodesInW[bad.rank]), 2)

	for _, action := range actions {
		bh.performeAction(action)
	}
}

// big TODO
func (bh *BrodalHeap) performeAction(act action) {
	numOfSonsOfT2 := 0
	notSonsOfT2 := []*node{}

	for e := bh.tree1.rankPointersW[act.index].violatingSelf; e.Value.(*node).rank != uint(act.index); e = e.Next() {
		if e.Value.(*node).parent == bh.tree2.root {
			numOfSonsOfT2++
		} else {
			notSonsOfT2 = append(notSonsOfT2, e.Value.(*node))
		}
	}

	if numOfSonsOfT2 > 4 {
		numOfRemoved := 0
		for _, rmNode := range notSonsOfT2 {

			bh.tree1.removeViolatingNode(rmNode)
			rmNode.removeSelfFromViolating()
			bh.tree1.Insert(rmNode)
			numOfRemoved++
		}

		for e := bh.tree1.rankPointersW[act.index].violatingSelf; e.Value.(*node).rank != uint(act.index) && numOfRemoved < 2; {
			bh.tree1.rankPointersW[act.index] = e.Next().Value.(*node)

			e.Value.(*node).removeSelfFromViolating()
			bh.tree2.cut(e.Value.(*node))

			bh.tree1.Insert(e.Value.(*node))
			numOfRemoved++
		}

		bh.tree1.numOfNodesInW[act.index] -= uint(numOfRemoved)

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
			bh.tree1.removeViolatingNode(notGood)
		}

	}
}
