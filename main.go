package main

import (
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/mislavzanic/heaps/Brodal"
	// "github.com/mislavzanic/heaps/Fibonacci"
)

func testSorting[T Number, I any](n int64) {
	randArr := createRandomValueArray[T](n)
	boHeap := Brodal.NewEmptyHeap[T]()
	// fbHeap := Fibonacci.NewFibHeap[T]()

	sortingPerformance[T](boHeap, randArr)
	// sortingPerformance[T, I](fbHeap, randArr)
}

func testDijkstra[T Number](n int64) {
	// randMatrix := createRandomMatrix[T](n)

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

	testSorting[float64, any](heapSize)
}
