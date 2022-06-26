package Brodal

import (
	"container/list"
	"errors"
	"fmt"
)

type reduceType byte

const (
	linkReduce reduceType = 0
	cutReduce             = 1
)

const UPPER_BOUND int = 7
const LOWER_BOUND int = -2
const GUIDE_BOUND int = 6

type tree[T Number] struct {
	root            *node[T]
	id              uint
	childrenRank    []*node[T]
	upperBoundGuide *guide
	lowerBoundGuide *guide
	parentHeap      *BrodalHeap[T]
}

func newTree[T Number](value T, ident any, treeIndex uint, pH *BrodalHeap[T]) *tree[T] {
	return &tree[T]{
		root:            newNode(value, ident),
		id:              treeIndex,
		childrenRank:    []*node[T]{},
		upperBoundGuide: newGuide(UPPER_BOUND),
		lowerBoundGuide: newGuide(LOWER_BOUND),
		parentHeap: pH,
	}
}

func (this *tree[T]) RootRank() int {
	return this.root.rank
}

func (this *tree[T]) RootValue() T {
	return this.root.value
}

func (this *tree[T]) SetRootValue(newValue T) {
	this.root.SetValue(newValue)
}

func (this *tree[T]) NumOfRootChildren(rank int) int {
	if rank < 0 {
		return -1
	}
	if rank >= this.RootRank() {
		panic("rank >= this.RootRank()")
	}
	return this.root.numOfChildren[rank]
}

func (this *tree[T]) ChildOfRank(rank int) *node[T] {
	if child, err := this.childOfRank(rank); err != nil {
		panic(fmt.Sprint(err.Error()))
	} else {
		return child
	}
}

func (this *tree[T]) childOfRank(rank int) (*node[T], error) {
	if rank < this.RootRank() {
		return this.childrenRank[rank], nil

	} else {
		return nil, errors.New(fmt.Sprintf("Rank %d is greater than roots rank %d", rank, this.RootRank()))
	}
}

func (this *tree[T]) vList() *list.List {
	return this.root.vList
}

func (this *tree[T]) wList() *list.List {
	return this.root.wList
}

func (tree *tree[T]) Children() *list.List {
	return tree.root.children
}

func (tree *tree[T]) RmRfRoot() *list.List {
	children := tree.Children()
	for e := children.Front(); e != nil; e = e.Next() {
		e.Value.(*node[T]).parent = nil
	}
	tree.childrenRank = nil
	return children
}

func (tree *tree[T]) LeftChild() *node[T] {
	return tree.root.LeftChild()
}

func (this *tree[T]) insertNode(child *node[T]) (bool, error) {

	if len(this.childrenRank) < child.rank {
		panic("preveliko je")
	}

	rankInc := this.MbyIncRank(child.rank == this.RootRank())

	if child.rank > this.RootRank() {
		panic("fakkk")
	}

	if _, err := this.root.pushBackChild(child, this.childrenRank[child.rank]); err == nil {
		if this.childrenRank[child.rank] == nil {
			this.childrenRank[child.rank] = child
		}
		return rankInc, nil
	} else {
		return false, err
	}
}

func (this *tree[T]) InsertNodes(children ...*node[T]) {
	for _, n := range children {
		if _, err := this.insertNode(n); err != nil {
			panic(fmt.Sprint(err.Error()))
		}
	}
}

func (tree *tree[T]) cutOffNode(child *node[T]) (*node[T], error) {

	if child == tree.childrenRank[child.rank] {
		tree.childrenRank[child.rank] = nil
		if tree.root.children.Len() > 1 {
			if tree.root.children.Back().Value.(*node[T]) != child {
				if rb, err := child.rightBrother(); err == nil {
					if rb.rank == child.rank {
						tree.childrenRank[child.rank] = rb
					}
				} else {
					return nil, err
				}
			}
		}
	}

	_, err := tree.root.removeChild(child)
	if err != nil {
		panic(err.Error())
	}
	tree.parentHeap.mbyRemoveFromViolating(child)

	for f := tree.parentHeap.getTree(1).wList().Front(); f != nil; f = f.Next() {
		if f.Value.(*node[T]) == child {
			panic("tusam")
		}
	}
	return child, err
}

func (this *tree[T]) RemoveChildren(rank int, num int) []*node[T] {
	nodes := []*node[T]{}
	for i := 0; i < num; i++ {
		if n, err := this.cutOffNode(this.childrenRank[rank]); err == nil {
			nodes = append(nodes, n)
		} else {
			panic(fmt.Sprint(err.Error()))
		}
	}
	return nodes
}

func (this *tree[T]) RemoveChild(rank int) *node[T] {
	return this.RemoveChildren(rank, 1)[0]
}

func (tree *tree[T]) Delink() []*node[T] {
	if nodes, err := tree.root.delink(); err == nil {
		tree.childrenRank[tree.RootRank()-1] = tree.LeftChild()
		return nodes
	} else {
		panic(fmt.Sprint(err.Error()))
	}
}

func (this *tree[T]) MbyIncRank(condition bool) bool {
	if condition {
		this.incRank()
	}
	return condition
}

func (this *tree[T]) incRank() {
	this.root.incRank()

	if len(this.childrenRank) > this.RootRank()-1 {
		this.childrenRank[this.RootRank()-1] = nil
	} else {
		this.childrenRank = append(this.childrenRank, nil)
	}

	if this.RootRank()-2 > 0 {
		this.upperBoundGuide.expand(this.RootRank()-2, this.NumOfRootChildren(this.RootRank()-3))
		this.lowerBoundGuide.expand(this.RootRank()-2, -this.NumOfRootChildren(this.RootRank()-3))
	}
}

func (this *tree[T]) DecRank() []*node[T] {
	if nodes, err := this.root.decRank(); err != nil {
		panic(err.Error())
	} else {
		this.upperBoundGuide.remove(this.RootRank() - 2)
		this.lowerBoundGuide.remove(this.RootRank() - 2)
		// this.childrenRank[len(this.childrenRank) - 1] = nil
		return nodes
	}
}

func (tree *tree[T]) AskGuide(rank int, numOfChildren int, insert bool) []action {
	lbReduceVal := 2
	if tree.childrenRank[rank+1].numOfChildren[rank] == 3 {
		lbReduceVal = 3
	}

	if insert {
		// act2 := tree.lowerBoundGuide.forceDecrease(rank, -numOfChildren-1, lbReduceVal)
		return tree.upperBoundGuide.forceIncrease(rank, &tree.root.numOfChildren, 3)
	}
	// act1 := tree.upperBoundGuide.forceDecrease(rank, numOfChildren-1, 3)

	return tree.lowerBoundGuide.forceIncrease(rank, &tree.root.numOfChildren, lbReduceVal)
}

func (tree *tree[T]) Link(rank int) *node[T] {

	nodes := tree.RemoveChildren(rank, 3)
	minNode, nodeX, nodeY := getMinNodeFrom3(nodes[0], nodes[1], nodes[2])

	if _, err := minNode.link(nodeX, nodeY); err != nil {
		panic(fmt.Sprint(err.Error()))
	}

	if _, err := tree.insertNode(minNode); err != nil {
		panic(fmt.Sprint(err.Error()))
	}

	return minNode
}
