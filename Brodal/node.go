package Brodal

import (
	"container/list"
)

type node struct {
	value float64
	rank  uint

	self          *list.Element
	violatingSelf *list.Element

	parent *node

	children      *list.List
	numOfChildren []uint

	vList *list.List
	wList *list.List

	parentViolatingList *list.List
}

func newNode(value float64) *node {
	node := new(node)
	node.value = value
	node.rank = 0
	node.children = list.New()
	node.vList = list.New()
	node.wList = list.New()

	return node
}

func (parent *node) leftSon() *node {
	return parent.children.Front().Value.(*node)
}

func (this *node) leftBrother() *node {
	return this.self.Prev().Value.(*node)
}

func (this *node) rightBrother() *node {
	return this.self.Next().Value.(*node)
}

func (node *node) isGood() bool {
	return node.value > node.parent.value
}

func (parent *node) removeChild(child *node) uint {
	parent.children.Remove(child.self)
	parent.numOfChildren[child.rank]--

	parent.mbyUpdateRank()

	return parent.numOfChildren[child.rank]
}

func (parent *node) removeFirstChild() (*node, uint) {
	child := parent.children.Front().Value.(*node)
	numOfChildren := parent.removeChild(child)
	return child, numOfChildren
}

func (parent *node) addChild(child *node, newRightBrother *node) {
	if parent.rank == child.rank {
		panic("Increase rank of parent first")
	}

	if child.parent != nil {
		child.parent.removeChild(child)
	}


	child.parent = parent
	child.self = parent.children.InsertBefore(child, newRightBrother.self)

	parent.numOfChildren[child.rank]++
}

func (this *node) swapBrothers(other *node) {
	brother := func() *node {
		if other.leftBrother().rank == other.rank {
			return other.leftBrother()
		} else {
			return other.rightBrother()
		}
	}()
	this.parent.addChild(brother, this)
	other.parent.addChild(this, other)
}

func (parent *node) mbyUpdateRank() {
	parent.rank = parent.children.Front().Value.(*node).rank + 1
}

func (parent *node) incRank() {
	if len(parent.numOfChildren) != int(parent.rank) {
		parent.numOfChildren = append(parent.numOfChildren, 0)
	}
	parent.rank++
}

func (node *node) link(xNode *node, yNode *node) {

	if node.rank != xNode.rank || node.rank != yNode.rank {
		panic("Node ranks don't match")
	}

	node.incRank()
	node.addChild(xNode, node.leftSon())
	node.addChild(yNode, node.leftSon())
	node.mbyUpdateRank()

	node.parent.numOfChildren[node.rank]++
}

func (parent *node) delink() []*node {
	node1, _ := parent.removeFirstChild()
	node2, _ := parent.removeFirstChild()
	if parent.numOfChildren[parent.rank-1] == 1 {
		node3, _ := parent.removeFirstChild()
		return []*node{node1, node2, node3}
	}
	return []*node{node1, node2}
}

func (this *node) mbyRemoveSelfFromViolating() {
	if this.parentViolatingList != nil {
		this.parentViolatingList.Remove(this.violatingSelf)
		this.parentViolatingList = nil
	}
}
