package Brodal

import (
	"container/list"
	"errors"
	// "errors"
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

type tree struct {
	root            *node
	id              uint
	childrenRank    []*node
	upperBoundGuide *guide
	lowerBoundGuide *guide
}

func newTree(value valType, treeIndex uint) *tree {
	return &tree{
		root:            newNode(value),
		id:              treeIndex,
		childrenRank:    []*node{},
		upperBoundGuide: newGuide(UPPER_BOUND),
		lowerBoundGuide: newGuide(LOWER_BOUND),
	}
}

func (this *tree) RootRank() int {
	return this.root.rank
}

func (this *tree) RootValue() valType {
	return this.root.value
}

func (this *tree) NumOfRootChildren(rank int) int {
	if rank < 0 { return -1 }
	if rank >= this.RootRank() {
		panic("rank >= this.RootRank()")
	}
	return this.root.numOfChildren[rank]
}

func (this *tree) ChildOfRank(rank int) *node {
	if child, err := this.childOfRank(rank); err != nil {
		panic(fmt.Sprint(err.Error()))
	} else {
		return child
	}
}

func (this *tree) childOfRank(rank int) (*node, error) {
	if rank < this.RootRank() {
		return this.childrenRank[rank], nil
	} else {
		return nil, errors.New(fmt.Sprintf("Rank %d is greater than roots rank", rank, this.RootRank()))
	}
}

func (this *tree) vList() *list.List {
	return this.vList()
}

func (this *tree) wList() *list.List {
	return this.wList()
}

func (tree *tree) Children() *list.List {
	return tree.root.children
}

func (tree *tree) RmRfRoot() *list.List {
	children := tree.Children()
	for e := children.Front(); e != nil; e = e.Next() {
		e.Value.(*node).parent = nil
	}
	tree.childrenRank = nil
	return children
}

func (tree *tree) LeftChild() *node {
	return tree.root.LeftChild()
}

func (this *tree) insertNode(child *node) error {

	if len(this.childrenRank) < child.rank {
		panic("preveliko je")
	}

	this.MbyIncRank(child.rank == this.RootRank())

	if _, err := this.root.pushBackChild(child, this.childrenRank[child.rank]); err == nil {
		if this.childrenRank[child.rank] == nil {
			this.childrenRank[child.rank] = child
		}
		return nil
	} else {
		return err
	}
}

func (this *tree) InsertNodes(children ...*node) {
	for _, n := range children {
		if err := this.insertNode(n); err != nil {
			panic(fmt.Sprint(err.Error()))
		}
	}
}

func (tree *tree) cutOffNode(child *node) (*node, error) {

	if child == tree.childrenRank[child.rank] {
		tree.childrenRank[child.rank] = nil
		if tree.root.children.Len() > 1 {
			if tree.root.children.Back().Value.(*node) != child {
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

	return tree.root.removeChild(child)
}

func (this *tree) RemoveChildren(rank int, num int) []*node {
	nodes := []*node{}
	for i := 0; i < num; i++ {
		if n, err := this.cutOffNode(this.childrenRank[rank]); err == nil {
			nodes = append(nodes, n)
		} else {
			panic(fmt.Sprint(err.Error()))
		}
	}
	return nodes
}

func (this *tree) RemoveChild(rank int) *node {
	return this.RemoveChildren(rank, 1)[0]
}

func (tree *tree) Delink() []*node {
	if nodes, err := tree.root.delink(); err == nil {
		tree.childrenRank[tree.RootRank() - 1] = tree.LeftChild()
		return nodes
	} else {
		panic(fmt.Sprint(err.Error()))
	}
}

func (this *tree) MbyIncRank(condition bool) {
	if condition {
		this.incRank()
	}
}

func (this *tree) incRank() {
	this.root.incRank()

	if len(this.childrenRank) < this.RootRank() - 1 {
		this.childrenRank[this.RootRank() - 1] = nil
	} else {
		this.childrenRank = append(this.childrenRank, nil)
	}

	if this.RootRank() - 2 > 0 {
		this.upperBoundGuide.expand(this.RootRank() - 2, this.NumOfRootChildren(this.RootRank() - 3))
		this.lowerBoundGuide.expand(this.RootRank() - 2, -this.NumOfRootChildren(this.RootRank() - 3))
	}
}

func (this *tree) DecRank() []*node {
	if nodes, err := this.root.decRank(); err != nil {
		panic(err.Error())
	} else {
		this.upperBoundGuide.remove(this.RootRank() - 2)
		this.lowerBoundGuide.remove(this.RootRank() - 2)
		return nodes
	}
}

func (tree *tree) AskGuide(rank int, numOfChildren int, insert bool) ([]action,[]action) {
	lbReduceVal := 2
	if tree.childrenRank[rank+1].numOfChildren[rank] == 3 {
		lbReduceVal = 3
	}

	if insert {
		act1 := tree.upperBoundGuide.forceIncrease(rank, numOfChildren+1, 3)
		act2 := tree.lowerBoundGuide.forceDecrease(rank, -numOfChildren-1, lbReduceVal)
		return act1, act2
	}
	act2 := tree.lowerBoundGuide.forceIncrease(rank, -numOfChildren+1, lbReduceVal)
	act1 := tree.upperBoundGuide.forceDecrease(rank, numOfChildren-1, 3)

	return act2, act1
}

func (tree *tree) Link(rank int) *node {

	nodes := tree.RemoveChildren(rank, 3)
	minNode, nodeX, nodeY := getMinNodeFrom3(nodes[0], nodes[1], nodes[2])

	if _, err := minNode.link(nodeX, nodeY); err != nil {
		panic(fmt.Sprint(err.Error()))
	}

	if err := tree.insertNode(minNode); err != nil {
		panic(fmt.Sprint(err.Error()))

	}

	return minNode
}

