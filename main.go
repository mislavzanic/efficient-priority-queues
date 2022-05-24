package main

import (
	"math/rand"
	"time"

	"github.com/mislavzanic/heaps/Brodal"
)

func main() {
	heap := Brodal.NewEmptyHeap()
	rand.Seed(time.Now().Unix())
	for i := 0; i < 100000000; i++ {
		// println(i)
		val := rand.Float64()
		heap.Insert(Brodal.ValType(val))
	}
}
