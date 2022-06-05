package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/mislavzanic/heaps/Brodal"
)

func test_rand_values() {
	heap := Brodal.NewHeap(Brodal.ValType(rand.Float64()))

	for i := 0; i < 10000; i++ {
		val := rand.Float64()
		heap.Insert(Brodal.ValType(val))
	}

	min := heap.Min()
	for !heap.Empty() {
		newMin := heap.DeleteMin()
		// println(newMin)
		if min > newMin {
			panic("nije dobar min")
		}
		min = newMin
	}

}

func main() {
	rand.Seed(time.Now().Unix())
	for i := 0; i < 10; i++ {
		start := time.Now()
		test_rand_values()
		elapsed := time.Since(start)
		fmt.Printf("%s vrijeme...\n", elapsed)
	}
	// heap.Insert(100)
	// for i := 0; i > -100; i-- {
	// 	heap.Insert(Brodal.ValType(i))
	// 	heap.Insert(Brodal.ValType(-i))
	// }
}
