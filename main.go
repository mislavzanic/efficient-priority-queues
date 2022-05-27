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
		println(heap.Size())
	}

	for i := 5; i < 10; i++ {
		heap.Insert(Brodal.ValType(-i))
		println(heap.Size())
	}



	heap.DeleteMin()
	println(heap.Size())
	heap.DeleteMin()
	println(heap.Size())
	heap.DeleteMin()

	for i := 5; i < 10; i++ {
		heap.Insert(Brodal.ValType(i))
		heap.Insert(Brodal.ValType(-i))
	}

	heap.DeleteMin()
	heap.DeleteMin()
	heap.DeleteMin()

	for i := 0; i < 30000; i++ {
		val := rand.Float64()
		heap.Insert(Brodal.ValType(val))
	}

	size := heap.Size()
	min := heap.Min()
	for !heap.Empty() {
		newMin := heap.DeleteMin()
		newSize := heap.Size()
		if newSize + 1 != size {
			println("velicine:", newSize, size)
			panic("error")
		}
		size = newSize
		// println(newMin)
		if min > newMin {
			panic("nije dobar min")
		}
		min = newMin
	}
}
