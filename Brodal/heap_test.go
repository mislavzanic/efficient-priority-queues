package Brodal

import (
	"testing"
)

func TestHeap(t *testing.T) {
	pq := NewHeap(1)
	pq.Insert(2)

	if pq.tree2 == nil {
		t.Error("Tree2 is nil")
	}

	if pq.tree2.root.value != 2 {
		t.Errorf("Tree2 root value is %f, not 2", pq.tree2.root.value)
	} else {
		t.Logf("Tree2 root value is %f", pq.tree2.root.value)
	}

	if pq.Min() != 1 {
		t.Errorf("PQ min value is %f, not 1", pq.Min())
	}

	if pq.tree1.RootRank() != 0 {
		t.Errorf("Tree1 rank is %d, not 0", pq.tree1.RootRank())
	}
	if pq.tree2.RootRank() != 0 {
		t.Errorf("Tree2 rank is %d, not 0", pq.tree2.RootRank())
	}

	pq.Insert(3)

	if pq.tree2 != nil {
		t.Error("Tree2 is not nil")
	}

	if pq.tree1.RootRank() != 1 {
		t.Errorf("Tree1 rank is %d, not 1", pq.tree1.RootRank())
	}
}
