package Brodal

import (
	"container/list"
)

type node struct {
    value           float64
	rank            uint

	self            *list.Element

	leftBrother     *node
	rightBrother    *node
	parent          *node

	children        *list.List
	childrenRanks   []uint

	vList           *list.List
	wList           *list.List

	nextInViolating *node
	prevInViolating *node
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

func (node *node) leftSon() *list.Element {
	return node.children.Front()
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
	child.leftBrother.rightBrother = child.rightBrother
	child.rightBrother.leftBrother = child.leftBrother
	parent.children.Remove(child.self)
	parent.childrenRanks[child.rank]--
}

func (parent *node) addChild(child *node) {

	if child.parent != nil {
		child.parent.removeChild(child)
	}

	child.parent = parent

	child.leftBrother = parent.children.Front().Value.(*node).leftBrother
	child.rightBrother = parent.children.Front().Value.(*node)

	parent.children.Front().Value.(*node).leftBrother.rightBrother = child
	parent.children.Front().Value.(*node).leftBrother = child

	child.self = parent.children.PushFront(child)

	parent.childrenRanks[child.rank]++
}

func (node *node) incRank(subNode1 *node, subNode2 *node) {

	if node.rank != subNode1.rank || node.rank != subNode2.rank {
		panic("Node ranks don't match")
	}

	node.childrenRanks = append(node.childrenRanks, 0)

	node.addChild(subNode1)
	node.addChild(subNode2)
}

func (node *node) link(xNode *node, yNode *node) {

	if xNode.rank + 1 != node.rank || yNode.rank + 1 != node.rank {
		panic("Only allowed to link nodes of rank r(x) - 1")
	}

	node.addChild(xNode)
	node.addChild(yNode)
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
