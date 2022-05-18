package Brodal

import "container/list"

type tree struct {
	root            *node
	id              uint
	childrenRank    []*node
	upperBoundGuide *guide
	lowerBoundGuide *guide
}

type reduceType byte

const (
	linkReduce reduceType = 0
	cutReduce             = 1
)

const UPPER_BOUND int = 7
const LOWER_BOUND int = -2

func (tree *tree) RootRank() int {
	return tree.root.rank
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

func (tree *tree) addRootChild(child *node) {
	tree.root.pushBackChild(child, tree.childrenRank[child.rank])
}

func (tree *tree) linkToRoot(child1 *node, child2 *node) {
	if len(tree.childrenRank) == int(child1.rank) {
		tree.childrenRank = append(tree.childrenRank, child1)
	} else {
		tree.childrenRank[child1.rank] = child1
	}

	tree.root.link(child1, child2)

	tree.upperBoundGuide.expand(tree.RootRank() - 2)
	tree.lowerBoundGuide.expand(tree.RootRank() - 2)
}

func (tree *tree) removeRootChild(child *node) *node {

	if child == tree.childrenRank[child.rank] {
		if tree.root.children.Len() > 1 {
			if tree.root.children.Back().Value.(*node) != child {
				if child.rightBrother().rank == child.rank {
					tree.childrenRank[child.rank] = child.rightBrother()
				}
			}
		} else {
			tree.childrenRank[child.rank] = nil
		}
	}

	tree.root.removeChild(child)

	return child
}

func (tree *tree) delink() []*node {
	nodes := tree.root.delink()
	tree.childrenRank[tree.RootRank() - 1] = tree.LeftmostSon()
	return nodes
}

func (tree *tree) incRank(node1 *node, node2 *node) {
	if tree.RootRank() > node1.rank || tree.RootRank() > node2.rank {
		panic("Tree ranks don't match")
	}

	tree.root.link(node1, node2)
	if len(tree.childrenRank) == node1.rank {
		tree.childrenRank = append(tree.childrenRank, nil)
	}
	tree.childrenRank[node1.rank] = node1
	tree.upperBoundGuide.expand(tree.RootRank() - 2)
	tree.lowerBoundGuide.expand(tree.RootRank() - 2)
}

func (tree *tree) askGuide(rank int, numOfChildren int, increase bool) []action {
	if increase {
		return tree.upperBoundGuide.forceIncrease(rank, numOfChildren+1, 3)
	}

	reduceVal := 2
	if tree.childrenRank[rank+1].numOfChildren[rank] == 3 {
		reduceVal = 3
	}
	return tree.lowerBoundGuide.forceIncrease(rank, -numOfChildren+1, reduceVal)
}

func (tree *tree) link(rank int) {
	nodeX := tree.childrenRank[rank]
	nodeY, nodeZ := nodeX.rightBrother(), nodeX.rightBrother().rightBrother()

	if nodeZ.rightBrother().rank == rank {
		tree.childrenRank[rank] = nodeZ.rightBrother()
	} else {
		tree.childrenRank[rank] = nil
	}

	minNode, nodeX, nodeY := getMinNodeFrom3(nodeX, nodeY, nodeZ)

	minNode.link(nodeX, nodeY)

	if minNode.rank == len(tree.childrenRank) {
		tree.childrenRank = append(tree.childrenRank, minNode)
	}

	tree.upperBoundGuide.expand(tree.RootRank() - 2)
	tree.lowerBoundGuide.expand(tree.RootRank() - 2)
}

func newTree(value float64, treeIndex uint) *tree {

	tree := &tree{
		root:            newNode(value),
		id:              treeIndex,
		childrenRank:    []*node{},
		upperBoundGuide: newGuide(UPPER_BOUND),
		lowerBoundGuide: newGuide(LOWER_BOUND),
	}

	return tree
}
