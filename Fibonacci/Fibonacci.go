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
	val      float64
	degree   uint
	mark     bool
}

func newNode(value float64) *node {
	node := new(node)
	node.children = list.New()
	node.val = value
	return node
}

type FibHeap struct {
	min  *node
	roots *list.List
	size     uint
}

func NewFibHeap() *FibHeap {
	return &FibHeap{
		min:   nil,
		roots: list.New(),
		size:  0,
	}
}

func (heap *FibHeap) maxDegree() int {
	return int(math.Floor(math.Log(float64(heap.size)) * 2))
}

func (heap *FibHeap) Empty() bool {
	return heap.size == 0
}

func (heap *FibHeap) Insert(value float64) {
	node := newNode(value)
	node.self = heap.roots.PushBack(node)
	heap.size++

	if heap.min == nil || heap.min.val > value {
		heap.min = node
	}
}

func (heap *FibHeap) Union(other *FibHeap) {
	for e := other.roots.Front(); e != nil; e = e.Next() {
		if e.Value.(*node) == other.min {
			if heap.min == nil || heap.min.val > other.min.val {
				heap.min = e.Value.(*node)
			}
		}
		e.Value.(*node).self = heap.roots.PushBack(e.Value)
	}

	heap.size += other.size;
}

func (heap *FibHeap) ExtractMin() *node {
	min := heap.min

	if min != nil {

		for e := min.children.Front(); e != nil; e = e.Next() {
			e.Value.(*node).parent = nil;
			e.Value.(*node).self = heap.roots.PushBack(e.Value)
		}

		heap.roots.Remove(heap.min.self)
		if heap.roots.Len() == 0 {
			heap.min = nil
		} else {
			heap.min = heap.roots.Front().Value.(*node)
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

	for tree := heap.roots.Front(); tree != nil; {

		degree := tree.Value.(*node).degree

		if degrees[degree] == nil {
			degrees[degree] = tree
			tree = tree.Next()
		} else if degrees[degree] == tree {
			tree = tree.Next()
		} else {
			for degrees[tree.Value.(*node).degree] != nil {

				otherTree := degrees[tree.Value.(*node).degree]
				degrees[tree.Value.(*node).degree] = nil

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
	heap.roots.Remove(other.self)
	other.mark = false
	other.self = node.children.PushBack(other)
	node.degree++
}

func (heap *FibHeap) resetMin() {
	for root := heap.roots.Front(); root != nil; root = root.Next() {
		if heap.min.val > root.Value.(*node).val {
			heap.min = root.Value.(*node)
		}
	}
}

func (heap *FibHeap) DecreaseKey(node *node, val float64) {
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
	child.self = heap.roots.PushBack(child)
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

func (heap *FibHeap) Delete(node *node) {
	heap.DecreaseKey(node, math.Inf(-1))
	heap.ExtractMin()
}
