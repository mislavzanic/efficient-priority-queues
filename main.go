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

type Number interface {
	int | int8 | int16 | int32 | int64 | float32 | float64
}

type PriorityQueue[T Number] interface {
	Meld(interface{})
	Min() T
	DeleteMin() T
	Insert(T)
	Empty() bool
}

func createRandomValueArray(n int64) []float64 {
	arr := []float64{}
	for i := 0; int64(i) < n; i++ {
		arr = append(arr, rand.Float64())
	}
	return arr
}

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


// func test_rand_values[T Number]() {
// 	heap := Brodal.NewHeap(T(rand.Float64()))

// 	for i := 0; i < 100000; i++ {
// 		val := rand.Float64()
// 		heap.Insert(T(val))
// 	}

// 	min := heap.Min()
// 	for !heap.Empty() {
// 		newMin := heap.DeleteMin()
// 		// println(newMin)
// 		if min > newMin {
// 			panic("nije dobar min")
// 		}
// 		min = newMin
// 	}

// }

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

	testInsertAndDeleteMin(boHeap, randArr)
	testInsertAndDeleteMin(fbHeap, randArr)
}
