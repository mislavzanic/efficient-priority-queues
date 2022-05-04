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

	vList           *list.List
	wlist           *list.List

	nextInViolating *node
	prevInViolating *node
}

func newNode(value float64) *node {
	node := new(node)
	node.value = value
	node.rank = 0
	node.children = list.New()
	node.vList = list.New()
	node.wlist = list.New()

	return node
}

func (node *node) leftSon() *list.Element {
	return node.children.Front()
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
}

func (parent *node) addChild(child *node) {

	if child.parent != nil {
		child.parent.removeChild(child)
	}

	child.parent = parent
	child.self = parent.children.PushFront(child)
	// TODO dodavanje lijevog i desnog brata...
}

func (node *node) link(xNode *node, yNode *node) {
	node.addChild(xNode)
	node.addChild(yNode)
	node.rank++
}
