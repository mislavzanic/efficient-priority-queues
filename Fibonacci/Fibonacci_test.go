package Fibonacci

import (
	// "fmt"
	"math/rand"
	"testing"
	"time"
)

func TestInsert(t *testing.T) {
	heap := NewFibHeap()
	heap.Insert(-1)
	heap.Insert(1)
	heap.Insert(1000)
	heap.Insert(-5)

	if heap.size != 4 {
		t.Error("Heap size is not 4")
	}

	if heap.min.val != -5 {
		t.Error("Heap min value is not -5")
	}
}

func TestExtractMin(t *testing.T) {
	heap := NewFibHeap()
	rand.Seed(time.Now().Unix())

	for i := 0; i < 1000; i++ {
		val := rand.Float64()
		heap.Insert(val)
	}

	if heap.size != 1000 {
		t.Error("Wrong heap size")
	}

	lastMin := heap.ExtractMin()
	for i := 0; i < 999; i++ {
		min := heap.ExtractMin()
		if lastMin.val > min.val {
			// t.Error(fmt.Sprint("lastMin with value of %f is bigger than current min with value of %f", lastMin.val, min.val))
			t.Error("Wrong min")
		}
		if heap.size != uint(999 - i - 1) {
			t.Error("Wrong heap size")
		}
		lastMin = min
	}
}

// func TestDelete(t *testing.T) {

// }
