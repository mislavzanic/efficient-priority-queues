package Brodal

import (
	"testing"
	"math/rand"
	"time"
)

func checkMin(pq *BrodalHeap, expectedVal float64, t *testing.T) {
	if pq.Min() != expectedVal {
		t.Errorf("FAIL: PQ min value is %f, not %f", pq.Min(), expectedVal)
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

	if pq.t1s.tree1.RootRank() != 0 {
		t.Errorf("Tree1 rank is %d, not 0", pq.t1s.tree1.RootRank())
	}
	if pq.tree2.RootRank() != 0 {
		t.Errorf("Tree2 rank is %d, not 0", pq.tree2.RootRank())
	}

	pq.Insert(3)

	t.Logf("Tree1 left son has value %f", pq.t1s.tree1.LeftmostSon().value)

	checkIfTreeIsNil(pq.tree2, 2, t)

	if pq.t1s.tree1.RootRank() != 1 {
		t.Errorf("Tree1 rank is %d, not 1", pq.t1s.tree1.RootRank())
	}

	checkNumOfChildren(pq.t1s.tree1, 0, 2, t)

	pq.Insert(3)
	pq.Insert(100)
	pq.Insert(4)

	checkNumOfChildren(pq.t1s.tree1, 0, 5, t)

	pq.Insert(-1)

	if pq.t1s.tree1.RootRank() != 0 {
		t.Errorf("Tree1 rank is %d, not 0", pq.t1s.tree1.RootRank())
	}
	pq.Insert(10000)

	if pq.t1s.tree1.RootRank() != 0 {
		t.Errorf("Tree1 rank is %d, not 0", pq.t1s.tree1.RootRank())
	}
	pq.Insert(3)
	if pq.t1s.tree1.RootRank() != 0 {
		t.Errorf("Tree1 rank is %d, not 0", pq.t1s.tree1.RootRank())
	}
	pq.Insert(3)

	checkNumOfChildren(pq.tree2, 1, 2, t)


	if pq.t1s.tree1.RootRank() != 0 {
		t.Errorf("Tree1 rank is %d, not 0", pq.t1s.tree1.RootRank())
	}

	t.Log("DeleteMin...")
	pq.DeleteMin()

	checkMin(pq, 1, t)
	checkNumOfChildren(pq.t1s.tree1, 1, 2, t)
	t.Log("DeleteMin...")

	pq.DeleteMin()
	checkMin(pq, 2, t)

	pq.DeleteMin()
	checkMin(pq, 3, t)

	pq.DeleteMin()
	checkMin(pq, 3, t)

	pq.Insert(5)
	pq.Insert(-1)

	pq.DeleteMin()
	checkMin(pq, 3, t)

	pq.DeleteMin()
	checkMin(pq, 3, t)

	pq.DeleteMin()
	checkMin(pq, 4, t)

	pq.DeleteMin()
	checkMin(pq, 5, t)

	pq.Insert(6)
	pq.Insert(10)
	pq.Insert(1)
	checkMin(pq, 1, t)

	pq.DeleteMin()
	checkMin(pq, 5, t)

	pq.DeleteMin()
	checkMin(pq, 6, t)
}

func TestTest(t *testing.T) {
	heap := NewHeap(1)
	heap.Insert(3)
	heap.Insert(2)
}

func TestRandomInsertions(t *testing.T) {
	heap := NewHeap(1)
	rand.Seed(time.Now().Unix())

	for i := 0; i < 10000; i++ {
		val := rand.Float64()
		println(i, val)
		heap.Insert(val)
		// heap.Insert(float64(i))
	}
	// lastMin := heap.Min()
	// heap.DeleteMin()
	// for i := 0; i < 999; i++ {
	// 	min := heap.Min()
	// 	heap.DeleteMin()
	// 	if lastMin > min {
	// 		t.Error("Wrong min")
	// 	}
	// 	lastMin = min
	// }
}
