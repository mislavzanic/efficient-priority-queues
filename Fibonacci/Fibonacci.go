package Fibonacci

import (
	"container/list"
	"fmt"
	"math"
)

type node struct {
	self     *list.Element
	parent   *node
	children *list.List

	index    uint
	val      int
	degree   uint
	mark     bool
}

func newNode(value int) *node {
	node := new(node)
	node.children = list.New()
	node.val = value
	return node
}

type FibHeap struct {
	min  *node
	rootList *list.List
	size     uint
}

func NewFibHeap() *FibHeap {
	return &FibHeap{
		min:  nil,
		size: 0,
	}
}

func (heap *FibHeap) maxDegree() int {
	return int(math.Floor(math.Log(float64(heap.size)) * 2))
}

func (heap *FibHeap) Empty() bool {
	return heap.size == 0
}

func (heap *FibHeap) insert(value int) {
	node := newNode(value)
	node.self = heap.rootList.PushBack(node)
	heap.size++

	if heap.min == nil || heap.min.val > value {
		heap.min = node
	}
}

func (heap *FibHeap) Union(other *FibHeap) {
	for e := other.rootList.Front(); e != nil; e = e.Next() {
		if e.Value.(*node) == other.min {
			if heap.min == nil || heap.min.val > other.min.val {
				heap.min = e.Value.(*node)
			}
		}
		e.Value.(*node).self = heap.rootList.PushBack(e.Value)
	}

	heap.size += other.size;
}

func (heap *FibHeap) ExtractMin() *node {
	min := heap.min

	if min != nil {

		for e := min.children.Front(); e != nil; e = e.Next() {
			e.Value.(*node).parent = nil;
			e.Value.(*node).self = heap.rootList.PushBack(e.Value)
		}

		heap.rootList.Remove(heap.min.self)
		if heap.rootList.Len() == 0 {
			heap.min = nil
		} else {
			heap.min = heap.rootList.Front().Value.(*node)
			heap.consolidate()
		}
		heap.size--
	}
	return min
}

func (heap *FibHeap) consolidate() {

	degrees := [](*list.Element){}
	for i := 0; i < heap.maxDegree(); i++ {
		degrees = append(degrees, nil)
	}

	for tree := heap.rootList.Front(); tree != nil; {

		degree := tree.Value.(*node).degree

		if degrees[degree] == nil {
			degrees[degree] = tree
			tree = tree.Next()
		} else if degrees[degree] == tree {
			tree = tree.Next()
		} else {
			for degrees[tree.Value.(*node).degree] != nil {

				otherTree := degrees[degree]
				degrees[degree] = nil

				if otherTree.Value.(*node).val < tree.Value.(*node).val {
					heap.link(otherTree.Value.(*node), tree.Value.(*node))
					tree = otherTree
				} else {
					heap.link(tree.Value.(*node), otherTree.Value.(*node))
				}
			}
			degrees[tree.Value.(*node).degree] = tree
		}
	}

	heap.resetMin()
}

func (heap *FibHeap) link(node *node, other *node) {
	heap.rootList.Remove(other.self)
	other.mark = false
	other.self = node.children.PushBack(other)
	node.degree++
}

func (heap *FibHeap) resetMin() {
	for root := heap.rootList.Front(); root != nil; root = root.Next() {
		if heap.min.val > root.Value.(*node).val {
			heap.min = root.Value.(*node)
		}
	}
}

func (heap *FibHeap) DecreaseKey(node *node, val int) {
	if node.val < val {
		panic(fmt.Sprintf("Existing key is smaller than the new key"))
	}

	// node.val = val
	// parent := node.parent
	if node.parent != nil && node.val < node.parent.val {
		heap.cut(node, node.parent)
		heap.cascadingCut(node.parent)
	}

	if node.val < heap.min.val {
		heap.min = node
	}
}

func (heap *FibHeap) cut(child *node, parent *node) {
	parent.children.Remove(child.self)
	parent.degree--
	child.parent = nil
	child.mark = false
	child.self = heap.rootList.PushBack(child)
}

func (heap *FibHeap) cascadingCut(node *node) {
	if node.parent != nil {
		if !node.mark {
			node.mark = !node.mark
		} else {
			heap.cut(node, node.parent)
			heap.cascadingCut(node.parent)
		}
	}
}
