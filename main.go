package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/mislavzanic/heaps/Brodal"
)

func main() {
	heap := Brodal.NewHeap(1)
	fmt.Println(heap.Min())
	rand.Seed(time.Now().Unix())
	for i := 0; i < 1000; i++ {
		val := rand.Float64()
		heap.Insert(val)
	}
	lastMin := heap.Min()
	heap.DeleteMin()
	for i := 0; i < 999; i++ {
		min := heap.Min()
		heap.DeleteMin()
		if lastMin > min {
			panic("Wrong min")
		}
		lastMin = min
	}
}
