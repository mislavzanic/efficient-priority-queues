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

func (node *node) link(xNode *node, yNode *node) {

}
