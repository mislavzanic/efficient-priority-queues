package main

import (
	"math/rand"
	"time"

	"github.com/mislavzanic/heaps/Brodal"
)

func test_rand_values() {
	heap := Brodal.NewEmptyHeap()
	rand.Seed(time.Now().Unix())
	for i := 0; i < 1000000; i++ {
		val := rand.Float64()
		heap.Insert(Brodal.ValType(val))
	}
}

func main() {
	rand.Seed(time.Now().Unix())
	heap := Brodal.NewHeap(1)
	heap.Insert(2)
	heap.Insert(3)
	heap.Insert(-2)
	heap.Insert(-1)

	for i := 5; i < 10; i++ {
		heap.Insert(Brodal.ValType(i))
	}

	for i := 5; i < 10; i++ {
		heap.Insert(Brodal.ValType(-i))
	}

	heap.DeleteMin()
	heap.DeleteMin()
	heap.DeleteMin()

	for i := 5; i < 10; i++ {
		heap.Insert(Brodal.ValType(i))
		heap.Insert(Brodal.ValType(-i))
	}

	heap.DeleteMin()
	heap.DeleteMin()
	heap.DeleteMin()

	for i := 0; i < 300; i++ {
		val := rand.Float64()
		heap.Insert(Brodal.ValType(val))
	}

	min := heap.Min()
	for !heap.Empty() {
		newMin := heap.DeleteMin()
		if min > newMin {
			panic("nije dobar min")
		}
		min = newMin
	}
}
