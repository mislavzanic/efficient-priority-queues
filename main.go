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
	println("bok")
}
