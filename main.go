package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/mislavzanic/heaps/Brodal"
	"github.com/mislavzanic/heaps/Fibonacci"
)

func testInsertAndDeleteMin[T float64](pq PriorityQueue[float64], randArr []float64) {
	start := time.Now()
	for _, n := range randArr {
		pq.Insert(n)
	}
	elapsed := time.Since(start)
	fmt.Printf("%s Meld vrijeme...\n", elapsed)

	start = time.Now()
	min := pq.Min()
	for !pq.Empty() {
		newMin := pq.DeleteMin()
		if min > newMin {
			panic("nije dobar min")
		}
		min = newMin
	}
	elapsed = time.Since(start)
	fmt.Printf("%s DeleteMin vrijeme...\n", elapsed)
}


func main() {
	rand.Seed(time.Now().Unix())
	var heapSize int64
	if len(os.Args) > 1 {
		n, err := strconv.ParseInt(os.Args[1], 10, 64)
		if err != nil {
			panic(err.Error())
		}
		heapSize = n
	}

	randArr := createRandomValueArray(heapSize)
	boHeap := Brodal.NewEmptyHeap[float64]()
	fbHeap := Fibonacci.NewFibHeap[float64]()

	sortingPerformance[float64](boHeap, randArr)
	sortingPerformance[float64](fbHeap, randArr)
}
