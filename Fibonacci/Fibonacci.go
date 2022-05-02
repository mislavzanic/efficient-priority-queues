package Fibonacci

type FibHeap struct {
	minNode *Node
	size    int
}

func NewFibHeap() *FibHeap {
	return &FibHeap{
		minNode: nil,
		size:    0,
	}
}
