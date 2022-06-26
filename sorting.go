package main

import (
	"fmt"
	"time"
)

func test(i any) {

}

func sortingPerformance[T Number](pq PriorityQueue[T], randArr []T) {
	start := time.Now()
	for i, n := range randArr {
		pq.Insert(n, i)
	}
	meldTime := time.Since(start)
	fmt.Printf("%s .. Meld vrijeme\n", meldTime)

	start = time.Now()
	min := (*pq.Min()).Value()
	for !pq.Empty() {
		newMin := (*pq.DeleteMin()).Value()
		if min > newMin {
			panic("nije dobar min")
		}
		min = newMin
	}
	deleteMinTime := time.Since(start)
	fmt.Printf("%s .. DeleteMin vrijeme\n", deleteMinTime)
	fmt.Printf("%s .. Sveukupno\n", meldTime + deleteMinTime)
}
