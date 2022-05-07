package Brodal

import "math"

type BrodalHeap struct {
    tree1 *tree
    tree2 *tree
}

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
			if currKey.rank > bh.tree1.root.rank {
				bh.tree1.root.vList.PushBack(currKey)
			} else {
				bh.tree1.root.wList.PushBack(currKey)
			}
			// bh.updateViolating()
			// TODO take care of a violation
		}
	}
}

func (bh *BrodalHeap) Delete(node *node) {
	bh.DecreaseKey(node, math.Inf(-1))
	bh.DeleteMin()
}

func (bh *BrodalHeap) DeleteMin() {

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
					nodes := maxTree.delink()
					for _, n := range nodes {
						maxTree.insert(n, maxTree == minTree)
					}
				}

				maxTree.insert(tree.root, minTree == maxTree)
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
