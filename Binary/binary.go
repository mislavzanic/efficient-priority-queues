package Binary

import "fmt"

type BinHeap struct {
	root *Node
}

func (heap *BinHeap) min() itemType {
	return heap.root.val
}

func (heap *BinHeap) insert(node *Node) {
	heap.root.insert(node);
}
