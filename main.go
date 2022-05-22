package main

import (
	// "fmt"
	// "math/rand"
	// "time"

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
	}

	for i := 10; i < 20; i++ {
		heap.Insert(Brodal.ValType(i))
	}

	for i := 10; i < 20; i++ {
		heap.Insert(Brodal.ValType(-i))
	}

	for i := 0; i < 30; i++ {
		heap.Insert(Brodal.ValType(-i))
		heap.Insert(Brodal.ValType(i))
	}

	for i := 0; i < 3; i++ {
		heap.Insert(Brodal.ValType(-i))
		println(i)
		heap.Insert(Brodal.ValType(i))
	}
}
