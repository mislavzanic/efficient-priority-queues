package Fibonacci

import (
	"container/list"
	"fmt"
	"math"
)

type Number interface {
	int | int8 | int16 | int32 | int64 | float32 | float64
}

type node[T Number, I any] struct {
	self     *list.Element
	parent   *node[T, I]
	children *list.List

	index    uint
	val      T
	ident    I
	degree   uint
	mark     bool
}

func (this *node[T, I]) Value() T {
	return this.val
}

func (this *node[T, I]) Ident() I {
	return this.ident
}

func newNode[T Number, I any](value T, ident I) *node[T, I] {
	node := new(node[T, I])
	node.children = list.New()
	node.val = value
	node.ident = ident
	return node
}

type FibHeap[T Number, I any] struct {
	min  *node[T, I]
	roots *list.List
	size     uint
}

func NewFibHeap[T Number, I any]() *FibHeap[T, I] {
	return &FibHeap[T, I]{
		min:   nil,
		roots: list.New(),
		size:  0,
	}
}

func (heap *FibHeap[T, I]) maxDegree() int {
	return int(math.Floor(math.Log(float64(heap.size)) * 2))
}

func (heap *FibHeap[T, I]) Empty() bool {
	return heap.size == 0
}

func (heap *FibHeap[T, I]) Insert(value T, ident I) {
	node := newNode(value, ident)
	node.self = heap.roots.PushBack(node)
	heap.size++

	if heap.min == nil || heap.min.val > value {
		heap.min = node
	}
}

func (heap *FibHeap[T, I]) Meld(fibHeapInterface interface{}) {
	other := fibHeapInterface.(*FibHeap[T, I])
	for e := other.roots.Front(); e != nil; e = e.Next() {
		if e.Value.(*node[T, I]) == other.min {
			if heap.min == nil || heap.min.val > other.min.val {
				heap.min = e.Value.(*node[T, I])
			}
		}
		e.Value.(*node[T, I]).self = heap.roots.PushBack(e.Value)
	}

	heap.size += other.size;
}

func (heap *FibHeap[T, I]) Min() T {
	return heap.min.val
}

func (heap *FibHeap[T, I]) DeleteMin() T {
	min := heap.min

	if min != nil {

		for e := min.children.Front(); e != nil; e = e.Next() {
			e.Value.(*node[T, I]).parent = nil;
			e.Value.(*node[T, I]).self = heap.roots.PushBack(e.Value)
		}

		heap.roots.Remove(heap.min.self)
		if heap.roots.Len() == 0 {
			heap.min = nil
		} else {
			heap.min = heap.roots.Front().Value.(*node[T, I])
			heap.consolidate()
		}
		heap.size--
	}
	return min.val
}

func (heap *FibHeap[T, I]) consolidate() {

	degrees := [](*list.Element){}
	for i := 0; i < heap.maxDegree(); i++ {
		degrees = append(degrees, nil)
	}

	for tree := heap.roots.Front(); tree != nil; {

		degree := tree.Value.(*node[T, I]).degree

		if degrees[degree] == nil {
			degrees[degree] = tree
			tree = tree.Next()
		} else if degrees[degree] == tree {
			tree = tree.Next()
		} else {
			for degrees[tree.Value.(*node[T, I]).degree] != nil {

				otherTree := degrees[tree.Value.(*node[T, I]).degree]
				degrees[tree.Value.(*node[T, I]).degree] = nil

				if otherTree.Value.(*node[T, I]).val < tree.Value.(*node[T, I]).val {
					heap.link(otherTree.Value.(*node[T, I]), tree.Value.(*node[T, I]))
					tree = otherTree
				} else {
					heap.link(tree.Value.(*node[T, I]), otherTree.Value.(*node[T, I]))
				}
			}
			degrees[tree.Value.(*node[T, I]).degree] = tree
		}
	}

	heap.resetMin()
}

func (heap *FibHeap[T, I]) link(node *node[T, I], other *node[T, I]) {
	heap.roots.Remove(other.self)
	other.mark = false
	other.self = node.children.PushBack(other)
	node.degree++
}

func (heap *FibHeap[T, I]) resetMin() {
	for root := heap.roots.Front(); root != nil; root = root.Next() {
		if heap.min.val > root.Value.(*node[T, I]).val {
			heap.min = root.Value.(*node[T, I])
		}
	}
}

func (heap *FibHeap[T, I]) DecreaseKey(node *node[T, I], val T) {
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

func (heap *FibHeap[T, I]) cut(child *node[T, I], parent *node[T, I]) {
	parent.children.Remove(child.self)
	parent.degree--
	child.parent = nil
	child.mark = false
	child.self = heap.roots.PushBack(child)
}

func (heap *FibHeap[T, I]) cascadingCut(node *node[T, I]) {
	if node.parent != nil {
		if !node.mark {
			node.mark = !node.mark
		} else {
			heap.cut(node, node.parent)
			heap.cascadingCut(node.parent)
		}
	}
}

func (heap *FibHeap[T, I]) Delete(node *node[T, I]) {
	heap.DecreaseKey(node, T(math.Inf(-1)))
	heap.DeleteMin()
}
