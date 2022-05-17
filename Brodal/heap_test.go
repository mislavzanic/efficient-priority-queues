package Brodal

import (
	"testing"
)

func checkMin(pq *BrodalHeap, expectedVal float64, t *testing.T) {
	if pq.Min() != expectedVal {
		t.Errorf("PQ min value is %f, not %f", pq.Min(), expectedVal)
	} else {
		t.Logf("PQ min value is %f", pq.Min())
	}
}

func checkNumOfChildren(tree *tree, rank int, expectedVal int, t *testing.T) {
	if tree.root.numOfChildren[rank] != expectedVal {
		t.Errorf("Num of children with rank %d is %d, not %d", rank, tree.root.numOfChildren[rank], expectedVal)
	} else {
		t.Logf("Num of children with rank %d is %d", rank, tree.root.numOfChildren[rank])
	}
}

func checkIfTreeIsNil(tree *tree, index int, t *testing.T) {
	if tree != nil {
		t.Errorf("Tree with index %d and root value %f is not nil", index, tree.root.value)
	} else {
		t.Logf("Tree with index %d is nil", index)
	}
}

func checkIfTreeIsNotNil(tree *tree, index int, t *testing.T) {
	if tree == nil {
		t.Errorf("Tree with index %d is nil", index)
	} else {
		t.Logf("Tree with index %d is not nil", index)
	}
}

func TestHeap(t *testing.T) {
	pq := NewHeap(1)
	pq.Insert(2)

	checkIfTreeIsNotNil(pq.tree2, 2, t)

	if pq.tree2.root.value != 2 {
		t.Errorf("Tree2 root value is %f, not 2", pq.tree2.root.value)
	} else {
		t.Logf("Tree2 root value is %f", pq.tree2.root.value)
	}

	checkMin(pq, 1, t)

	if pq.tree1.RootRank() != 0 {
		t.Errorf("Tree1 rank is %d, not 0", pq.tree1.RootRank())
	}
	if pq.tree2.RootRank() != 0 {
		t.Errorf("Tree2 rank is %d, not 0", pq.tree2.RootRank())
	}

	pq.Insert(3)

	checkIfTreeIsNil(pq.tree2, 2, t)

	if pq.tree1.RootRank() != 1 {
		t.Errorf("Tree1 rank is %d, not 1", pq.tree1.RootRank())
	}

	checkNumOfChildren(pq.tree1, 0, 2, t)

	pq.Insert(3)
	pq.Insert(3)
	pq.Insert(3)

	checkNumOfChildren(pq.tree1, 0, 5, t)

	pq.Insert(3)
	pq.Insert(3)

	checkNumOfChildren(pq.tree1, 0, 7, t)

	pq.Insert(3)

	if pq.tree1.RootRank() != 2 {
		t.Errorf("Tree1 rank is %d, not 2", pq.tree1.RootRank())
	}

	checkNumOfChildren(pq.tree1, 0, 2, t)
	checkNumOfChildren(pq.tree1, 1, 2, t)

	pq.DeleteMin()

	checkMin(pq, 2, t)
}
