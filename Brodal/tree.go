package Brodal

import (
	"container/list"
)

var UPPER_BOUND int = 7
var LOWER_BOUND int = -2

type tree struct {
	root          *node

	rankPointersW *list.List
	childrenRank  []*node

	upperBoundGuide *guide
	lowerBoundGuide *guide
	mainTreeGuideW  *guide
}

type reduceType byte
const (
	linkReduce reduceType = 0
	cutReduce             = 1
)

func (tree *tree) rank() uint {
	return tree.root.rank
}

func (tree *tree) addRootChild(child *node) {
	tree.root.addChild(child , tree.childrenRank[child.rank])
	if child.rank < tree.root.rank - 2 {
		tree.lowerBoundGuide.update(-int(tree.root.numOfChildren[child.rank]), child.rank)
	}
}

func (tree *tree) removeRootChild(child *node) {
	tree.root.removeChild(child)
	if child.rank < tree.root.rank - 2 {
		tree.upperBoundGuide.update(int(tree.root.numOfChildren[child.rank]), child.rank)
	}
}

func (tree *tree) insert(node *node, isMinTree bool) {
	tree.addRootChild(node)

	if node.rank < tree.root.rank - 2 {
		actions := tree.upperBoundGuide.forceIncrease(int(node.rank), int(tree.root.numOfChildren[node.rank]), 3)

		for _, act := range actions {
			tree.performeAction(node, act, linkReduce, isMinTree)
		}

	}

	tree.handleHighRank(tree.root.rank - 2, isMinTree)
	tree.handleHighRank(tree.root.rank - 1, isMinTree)

	if !isMinTree {
		// TODO
	}
}

func (tree *tree) cut(node *node, isMinTree bool) {
	tree.removeRootChild(node)

	if node.rank < tree.root.rank - 2 {
		if tree.childrenRank[node.rank].leftBrother().rank < node.rank + 1 {
			panic("Don't know what to do")
		}

		reduceVal := 2
		if tree.childrenRank[node.rank].leftBrother().numOfChildren[node.rank] == 3 {
			reduceVal = 3
		}

		actions := tree.lowerBoundGuide.forceIncrease(int(node.rank), int(tree.root.numOfChildren[node.rank]), int(reduceVal))
		for _, act := range actions {
			tree.performeAction(node, act, cutReduce, isMinTree)
		}

	}

	tree.handleHighRank(tree.root.rank - 2, isMinTree)
	tree.handleHighRank(tree.root.rank - 1, isMinTree)
}

func (tree *tree) incRank(node1 *node, node2 *node) {
	if tree.rank() > node1.rank || tree.rank() > node2.rank {
		panic("Tree ranks don't match")
	}

	tree.root.incRank(node1, node2)
	tree.childrenRank = append(tree.childrenRank, node2)
}

func (tree *tree) performeAction(node *node, action action, reduceType reduceType, isMinTree bool) {
	if reduceType == linkReduce {
		tree.link(uint(action.index))

		tree.lowerBoundGuide.update(-int(tree.root.numOfChildren[action.index]), uint(action.index))
		tree.lowerBoundGuide.update(-int(tree.root.numOfChildren[action.index + 1]), uint(action.index + 1))
	} else {

		removedTree := tree.childrenRank[action.index]
		tree.removeRootChild(removedTree)

		nodes, _ := removedTree.delink()
		for _, n := range nodes {
			tree.insert(n, isMinTree)
		}

		tree.insert(removedTree, isMinTree)
	}
}

func (tree *tree) link(rank uint) {
	nodeX := tree.childrenRank[rank]
	nodeY, nodeZ := nodeX.rightBrother(), nodeX.rightBrother().rightBrother()

	if nodeZ.rightBrother().rank == rank {
		tree.childrenRank[rank] = nodeZ.rightBrother()
	} else {
		tree.childrenRank[rank] = nil
	}

	minNode, nodeX, nodeY := getMinNode(nodeX, nodeY, nodeZ)

	minNode.link(nodeX, nodeY)
}

func (tree *tree) handleHighRank(rank uint, isMinTree bool) {
	if tree.root.numOfChildren[rank] > 7 {
		nodeSliceX, _ := tree.root.delink()
		nodeSliceY, _ := tree.root.delink()
		nodeSliceZ, _ := tree.root.delink()

		nodeSliceX[0].incRank(nodeSliceX[1], nodeSliceY[0])
		nodeSliceY[1].incRank(nodeSliceZ[0], nodeSliceZ[1])

		if rank == tree.root.rank - 1 {
			tree.incRank(nodeSliceX[0], nodeSliceY[0])
		} else {
			tree.insert(nodeSliceX[0], isMinTree)
			tree.insert(nodeSliceY[1], isMinTree)
		}
	} else if tree.root.numOfChildren[rank] < 2 {
		//...
	}
}

func (tree *tree) reduceViolaton(x1 *node, x2 *node) {
	if x1.isGood() || x2.isGood() {
		if x1.isGood() {
			x1.removeSelfFromViolating()
		}
		if x2.isGood() {
			x2.removeSelfFromViolating()
		}
	} else {
		if x1.parent != x2.parent {
			if x1.parent.value <= x2.parent.value {
				x1.swapBrothers(x2)
			} else {
				x2.swapBrothers(x1)
			}
		}

		if x1.parent.numOfChildren[x1.rank] == 2 {
			if x1.parent.rank == x1.rank + 1 {
				//... TODO -- after tree.cut() is finished
			} else {
				x1.parent.removeChild(x1)
				x1.parent.removeChild(x2)
				tree.insert(x1, true)
				tree.insert(x2, true)
			}
		} else {
			x1.parent.removeChild(x1)
			tree.insert(x1, true)
		}
	}
}

// ######################################## UTIL #######################################

func newTree(value float64, treeIndex uint) *tree {

	tree := &tree{
		root: newNode(value),
		rankPointersW: list.New(),
		childrenRank: nil,
		upperBoundGuide: newGuide(UPPER_BOUND),
		lowerBoundGuide: newGuide(LOWER_BOUND),
		mainTreeGuideW: nil,
	}

	if treeIndex == 1 {
		tree.mainTreeGuideW = newGuide(6)
	}

	return tree
}

func getMinTree(tree1 *tree, tree2 *tree) (*tree, *tree) {
	if tree1 == nil || tree2 == nil {
		panic("One of the trees is nil")
	}

	if tree1.root.value > tree2.root.value {
		return tree2, tree1
	} else {
		return tree1, tree2
	}
}

func getMaxTree(trees []*tree) (*tree, []*tree) {
	if len(trees) == 0 {
		panic("There are no trees")
	}

	maxTree := trees[0]
	newTrees := [](*tree){}

	for _, tree := range trees {
		if tree != nil {
			if maxTree.root.rank < tree.root.rank {
				newTrees = append(newTrees, maxTree)
				maxTree = tree
			} else {
				newTrees = append(newTrees, tree)
			}
		}
	}

	return maxTree, newTrees
}

func mbySwapTree(ptr1 *tree, ptr2 *tree, cond bool) (*tree, *tree) {
	if cond {
		temp := ptr1
		ptr1 = ptr2
		ptr2 = temp
	}
	return ptr1, ptr2
}
