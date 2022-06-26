package main

import "math"

// "fmt"

func dijkstra(pq PriorityQueue[float64], pathMatrix [][]float64) []int {
	start, end := 0, len(pathMatrix) - 1
	dist, prev := []float64{}, []int{}

	for i := 0; i < end + 1; i++ {
		if i == start {
			dist = append(dist, 0)
		} else {
			dist = append(dist, math.Inf(1))
		}
		prev = append(prev, -1)
		pq.Insert(dist[i], i)
	}

	for !pq.Empty() {
		// du, u := pq.DeleteMin()
		// for v := 0; v < end + 1; v++ {
		// 	if v != u {
		// 		alt := dist[u] + pathMatrix[v][u]
		// 		if alt < dist[v] && dist[u] != math.Inf(1) {
		// 			dist[v] = alt
		// 			prev[v] = u
		// 			pq.DecreaseKey()
		// 		}
		// 	}
		// }
	}
	return nil
}
