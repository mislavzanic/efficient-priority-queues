package Fibonacci

import (
	"container/list"
	"fmt"
	"math"
)

type Number interface {
	int | int8 | int16 | int32 | int64 | float32 | float64
}

type node[T Number] struct {
	self     *list.Element
	parent   *node[T]
	children *list.List

	index    uint
	val      T
	degree   uint
	mark     bool
}

func newNode[T Number](value T) *node[T] {
	node := new(node[T])
	node.children = list.New()
	node.val = value
	return node
}

type FibHeap[T Number] struct {
	min  *node[T]
	roots *list.List
	size     uint
}

func NewFibHeap[T Number]() *FibHeap[T] {
	return &FibHeap[T]{
		min:   nil,
		roots: list.New(),
		size:  0,
	}
}

func (heap *FibHeap[T]) maxDegree() int {
	return int(math.Floor(math.Log(float64(heap.size)) * 2))
}

func (heap *FibHeap[T]) Empty() bool {
	return heap.size == 0
}

func (heap *FibHeap[T]) Insert(value T) {
	node := newNode(value)
	node.self = heap.roots.PushBack(node)
	heap.size++

	if heap.min == nil || heap.min.val > value {
		heap.min = node
	}
}

func (heap *FibHeap[T]) Meld(fibHeapInterface interface{}) {
	other := fibHeapInterface.(*FibHeap[T])
	for e := other.roots.Front(); e != nil; e = e.Next() {
		if e.Value.(*node[T]) == other.min {
			if heap.min == nil || heap.min.val > other.min.val {
				heap.min = e.Value.(*node[T])
			}
		}
		e.Value.(*node[T]).self = heap.roots.PushBack(e.Value)
	}

	heap.size += other.size;
}

func (heap *FibHeap[T]) Min() T {
	return heap.min.val
}

func (heap *FibHeap[T]) DeleteMin() T {
	min := heap.min

	if min != nil {

		for e := min.children.Front(); e != nil; e = e.Next() {
			e.Value.(*node[T]).parent = nil;
			e.Value.(*node[T]).self = heap.roots.PushBack(e.Value)
		}

		heap.roots.Remove(heap.min.self)
		if heap.roots.Len() == 0 {
			heap.min = nil
		} else {
			heap.min = heap.roots.Front().Value.(*node[T])
			heap.consolidate()
		}
		heap.size--
	}
	return min.val
}

func (heap *FibHeap[T]) consolidate() {

	degrees := [](*list.Element){}
	for i := 0; i < heap.maxDegree(); i++ {
		degrees = append(degrees, nil)
	}

	for tree := heap.roots.Front(); tree != nil; {

		degree := tree.Value.(*node[T]).degree

		if degrees[degree] == nil {
			degrees[degree] = tree
			tree = tree.Next()
		} else if degrees[degree] == tree {
			tree = tree.Next()
		} else {
			for degrees[tree.Value.(*node[T]).degree] != nil {

				otherTree := degrees[tree.Value.(*node[T]).degree]
				degrees[tree.Value.(*node[T]).degree] = nil

				if otherTree.Value.(*node[T]).val < tree.Value.(*node[T]).val {
					heap.link(otherTree.Value.(*node[T]), tree.Value.(*node[T]))
					tree = otherTree
				} else {
					heap.link(tree.Value.(*node[T]), otherTree.Value.(*node[T]))
				}
			}
			degrees[tree.Value.(*node[T]).degree] = tree
		}
	}

	heap.resetMin()
}

func (heap *FibHeap[T]) link(node *node[T], other *node[T]) {
	heap.roots.Remove(other.self)
	other.mark = false
	other.self = node.children.PushBack(other)
	node.degree++
}

func (heap *FibHeap[T]) resetMin() {
	for root := heap.roots.Front(); root != nil; root = root.Next() {
		if heap.min.val > root.Value.(*node[T]).val {
			heap.min = root.Value.(*node[T])
		}
	}
}

func (heap *FibHeap[T]) DecreaseKey(node *node[T], val T) {
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

func (heap *FibHeap[T]) cut(child *node[T], parent *node[T]) {
	parent.children.Remove(child.self)
	parent.degree--
	child.parent = nil
	child.mark = false
	child.self = heap.roots.PushBack(child)
}

func (heap *FibHeap[T]) cascadingCut(node *node[T]) {
	if node.parent != nil {
		if !node.mark {
			node.mark = !node.mark
		} else {
			heap.cut(node, node.parent)
			heap.cascadingCut(node.parent)
		}
	}
}

func (heap *FibHeap[T]) Delete(node *node[T]) {
	heap.DecreaseKey(node, T(math.Inf(-1)))
	heap.DeleteMin()
}
