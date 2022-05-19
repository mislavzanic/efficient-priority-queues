package Brodal

import (
	"container/list"
	// "fmt"
)

type reduceType byte

const (
	linkReduce reduceType = 0
	cutReduce             = 1
)

const UPPER_BOUND int = 7
const LOWER_BOUND int = -2

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

func (this *tree) numOfRootChildren(rank int) int {
	if rank < 0 { return -1 }
	if rank >= this.RootRank() {
		panic("rank >= this.RootRank()")
	}
	return this.root.numOfChildren[rank]
}

func (tree *tree) Children() *list.List {
	return tree.root.children
}

func (tree *tree) rmRfRoot() *list.List {
	children := tree.Children()
	for e := children.Front(); e != nil; e = e.Next() {
		e.Value.(*node).parent = nil
	}
	tree.childrenRank = nil
	return children
}

func (tree *tree) LeftmostSon() *node {
	return tree.root.leftSon()
}

func (this *tree) insertNode(child *node) bool {

	if len(this.childrenRank) < child.rank {
		panic("preveliko je")
	}

	this.mbyIncRank(child.rank == this.RootRank())

	isViolation := this.root.pushBackChild(child, this.childrenRank[child.rank])
	if this.childrenRank[child.rank] == nil {
		this.childrenRank[child.rank] = child
	}

	return isViolation
}

// func (tree *tree) addRootChild(child *node) bool {
// 	return tree.root.pushBackChild(child, tree.childrenRank[child.rank])
// }

// func (tree *tree) linkToRoot(child1 *node, child2 *node) {
// 	tree.root.link(child1, child2)

// 	if child1.parent != tree.root {
// 		panic("child1 roditelj nije root")
// 	}

// 	if len(tree.childrenRank) == int(child1.rank) {
// 		tree.childrenRank = append(tree.childrenRank, child1)
// 	} else {
// 		tree.childrenRank[child1.rank] = child1
// 	}

// 	if tree.RootRank() - 2 > 0 {
// 	}
// }

func (tree *tree) removeRootChild(child *node) *node {

	if child == tree.childrenRank[child.rank] {
		tree.childrenRank[child.rank] = nil

		if tree.root.children.Len() > 1 {

			if tree.root.children.Back().Value.(*node) != child {

				if child.rightBrother().rank == child.rank {
					tree.childrenRank[child.rank] = child.rightBrother()
				}

			}

		}

	}

	return tree.root.removeChild(child)
}

func (this *tree) removeChildrenWithRank(rank int, num int) []*node {
	nodes := []*node{}
	for i := 0; i < num; i++ {
		nodes = append(nodes, this.removeRootChild(this.childrenRank[rank]))
	}
	return nodes
}

func (tree *tree) delink() []*node {
	nodes := tree.root.delink()
	tree.childrenRank[tree.RootRank() - 1] = tree.LeftmostSon()
	return nodes
}

func (this *tree) mbyIncRank(condition bool) {
	if condition {
		this.incRank()
	}
}
func (this *tree) incRank() {
	this.root.incRank()

	if len(this.childrenRank) < this.RootRank() {
		this.childrenRank[this.RootRank() - 1] = nil
	} else {
		this.childrenRank = append(this.childrenRank, nil)
	}

	if this.RootRank() - 2 > 0 {
		this.upperBoundGuide.expand(this.RootRank() - 2, this.numOfRootChildren(this.RootRank() - 3))
		this.lowerBoundGuide.expand(this.RootRank() - 2, -this.numOfRootChildren(this.RootRank() - 3))
	}
}

func (this *tree) decRank() []*node {
	// napravit nes sa vodiljama
	return this.removeChildrenWithRank(this.RootRank() - 1, this.numOfRootChildren(this.RootRank() - 1))
}

// func (tree *tree) incRank(node1 *node, node2 *node) {
// 	if tree.RootRank() > node1.rank || tree.RootRank() > node2.rank {
// 		panic("Tree ranks don't match")
// 	}

// 	tree.root.link(node1, node2)
// 	if len(tree.childrenRank) == node1.rank {
// 		tree.childrenRank = append(tree.childrenRank, nil)
// 	}
// 	if node1.parent != tree.root {
// 		panic("node1 nema root roditelja")
// 	}
// 	tree.childrenRank[node1.rank] = node1
// 	if tree.RootRank() - 2 > 0 {
// 		tree.upperBoundGuide.expand(tree.RootRank() - 2, tree.root.numOfChildren[tree.RootRank() - 3])
// 		tree.lowerBoundGuide.expand(tree.RootRank() - 2, -tree.root.numOfChildren[tree.RootRank() - 3])
// 	}
// }

func (tree *tree) askGuide(rank int, numOfChildren int, increase bool) ([]action,[]action) {
	lbReduceVal := 2
	if tree.childrenRank[rank+1].numOfChildren[rank] == 3 {
		lbReduceVal = 3
	}

	if increase {
		act1 := tree.upperBoundGuide.forceIncrease(rank, numOfChildren+1, 3)
		act2 := tree.lowerBoundGuide.forceDecrease(rank, -numOfChildren-1, lbReduceVal)
		return act1, act2
	}
	act2 := tree.lowerBoundGuide.forceIncrease(rank, -numOfChildren+1, lbReduceVal)
	act1 := tree.upperBoundGuide.forceDecrease(rank, numOfChildren-1, 3)

	return act2, act1
}

func (tree *tree) linkRank(rank int) *node {

	if tree.numOfRootChildren(rank) != 7 {
		panic("manje od 7")
	}

	nodes := tree.removeChildrenWithRank(rank, 3)
	minNode, nodeX, nodeY := getMinNodeFrom3(nodes[0], nodes[1], nodes[2])

	minNode.link(nodeX, nodeY)
	tree.insertNode(minNode)

	return minNode
}

