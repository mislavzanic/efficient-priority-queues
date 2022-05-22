package main

import (
	"github.com/mislavzanic/heaps/Brodal"
)

func main() {
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
		println(heap.ToString())
	}

	for i := 10; i < 20; i++ {
		heap.Insert(Brodal.ValType(i))
		println(heap.ToString())
	}

	for i := 10; i < 20; i++ {
		// bug nakon ubacivanja -18 ... mby??
		heap.Insert(Brodal.ValType(-i))
		println(heap.ToString())
	}

	for i := 0; i < 30; i++ {
		// -4 error daje
		heap.Insert(Brodal.ValType(-i))
		println(heap.ToString())
		heap.Insert(Brodal.ValType(i))
		println(heap.ToString())
	}

	for i := 0; i < 3; i++ {
		heap.Insert(Brodal.ValType(-i))
		println(i)
		heap.Insert(Brodal.ValType(i))
	}
}
