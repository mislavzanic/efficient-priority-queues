package Brodal

import (
	"container/list"
)

type node struct {
    value           float64
	rank            uint

	self            *list.Element
	violatingSelf   *list.Element

	parent          *node

	children        *list.List
	childrenRanks   []uint

	vList           *list.List
	wList           *list.List

	parentViolatingList *list.List
}

func mbySwapNode(ptr1 *node, ptr2 *node, cond bool) {
	if cond {
		temp := *ptr1
		*ptr1 = *ptr2
		*ptr2 = temp
	}
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

func getMinNode(xNode *node, yNode *node, zNode *node) (*node, *node, *node) {
	minNode := zNode

	if xNode.value < minNode.value {
		minNode = xNode
		xNode = zNode
	}

	if yNode.value < minNode.value {
		temp := minNode
		minNode = yNode
		yNode = temp
	}

	return minNode, xNode, yNode
}

func (parent *node) removeChild(child *node) {
	parent.children.Remove(child.self)
	parent.childrenRanks[child.rank]--

	if parent.childrenRanks[child.rank] == 0 {
		if child.rank == parent.rank - 1 {
			parent.rank = parent.leftSon().rank + 1
		}
	}
}

func (parent *node) addChild(child *node, newRightBrother *node) {

	if child.parent != nil {
		child.parent.removeChild(child)
	}

	child.parent = parent

	child.self = parent.children.InsertBefore(child, newRightBrother.self)

	parent.childrenRanks[child.rank]++
}

func (this *node) swapBrothers(other *node) {
	brother := func () *node { if other.leftBrother().rank == other.rank { return other.leftBrother() } else { return other.rightBrother() }}()
	this.parent.addChild(brother, this)
	other.parent.addChild(this, other)
}

func (node *node) incRank(subNode1 *node, subNode2 *node) {

	if node.rank != subNode1.rank || node.rank != subNode2.rank {
		panic("Node ranks don't match")
	}

	node.childrenRanks = append(node.childrenRanks, 0)

	node.addChild(subNode1, node.leftSon())
	node.addChild(subNode2, node.leftSon())
}

func (node *node) link(xNode *node, yNode *node) {

	if xNode.rank + 1 != node.rank || yNode.rank + 1 != node.rank {
		panic("Only allowed to link nodes of rank r(x) - 1")
	}

	node.addChild(xNode, node.leftSon())
	node.addChild(yNode, node.leftSon())
	node.rank++
	node.parent.childrenRanks[node.rank]++
}


func (parent *node) delink() *node {
	child := parent.children.Front().Value.(*node)
	parent.removeChild(child)
	if parent.childrenRanks[parent.rank - 1] == 0 {
		parent.rank = parent.children.Front().Value.(*node).rank + 1
	}
	return child
}

func (this *node) removeSelfFromViolating() {
	if this.parentViolatingList != nil {
		this.parentViolatingList.Remove(this.violatingSelf)
		this.parentViolatingList = nil
	}
}
