package main

import (
	"fmt"
	"time"
)

func sortingPerformance[T Number](pq PriorityQueue[T], randArr []T) {
	start := time.Now()
	for _, n := range randArr {
		pq.Insert(n)
	}
	meldTime := time.Since(start)
	fmt.Printf("%s .. Meld vrijeme\n", meldTime)

	start = time.Now()
	min := pq.Min()
	for !pq.Empty() {
		newMin := pq.DeleteMin()
		if min > newMin {
			panic("nije dobar min")
		}
		min = newMin
	}
	deleteMinTime := time.Since(start)
	fmt.Printf("%s .. DeleteMin vrijeme\n", deleteMinTime)
	fmt.Printf("%s .. Sveukupno\n", meldTime + deleteMinTime)
}
