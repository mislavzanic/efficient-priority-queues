package Brodal

import (
	"container/list"
	"fmt"
)

type node struct {
	value valType
	rank  int

	self          *list.Element
	violatingSelf *list.Element

	parent *node

	children      *list.List
	numOfChildren []int

	vList *list.List
	wList *list.List

	parentViolatingList *list.List
}

func newNode(value valType) *node {
	node := new(node)
	node.value = value
	node.rank = 0
	node.children = list.New()
	node.vList = list.New()
	node.wList = list.New()

	return node
}

func (parent *node) leftSon() *node {
	if parent.children == nil || parent.children.Len() == 0 {
		return nil
	}
	return parent.children.Front().Value.(*node)
}

func (this *node) leftBrother() *node {
	return this.self.Prev().Value.(*node)
}

func (this *node) rightBrother() *node {
	return this.self.Next().Value.(*node)
}

func (this *node) getMinFromW() *node {
	if this.wList.Len() == 0 { return nil }
	return getMinFromList(this.wList).Value.(*node)
}

func (this *node) getMinFromV() *node {
	if this.vList.Len() == 0 { return nil }
	return getMinFromList(this.vList).Value.(*node)
}

func (this *node) getMinFromChildren() *node {
	return getMinFromList(this.children).Value.(*node)
}

func (node *node) isGood() bool {
	return node.value >= node.parent.value
}

func (parent *node) removeChild(child *node) *node {
	parent.children.Remove(child.self)
	child.parent = nil
	parent.numOfChildren[child.rank]--

	parent.mbyUpdateRank()

	return child
}

func (parent *node) removeFirstChild() *node {
	child := parent.children.Front().Value.(*node)

	if parent.rank-1 != child.rank {
		panic("Incorrect ranks")
	}

	return parent.removeChild(child)
}

func (parent *node) pushBackChild(child *node, newBrother *node) bool {
	return parent.addBrother(child, newBrother, true)
}

func (parent *node) pushFrontChild(child *node, newBrother *node) bool {
	return parent.addBrother(child, newBrother, false)
}

func (parent *node) addBrother(child *node, newBrother *node, left bool) bool {
	if parent.rank == child.rank {
		panic("Increase rank of parent first")
	}

	if newBrother != nil {
		if newBrother.parent != parent {
			panic(fmt.Sprintf("newBrother parent is not the parent, %f, %f", parent.value, newBrother.parent.value))
		}
	}

	if child.parent != nil {
		child.parent.removeChild(child)
	}

	child.parent = parent
	if newBrother == nil {
		child.self = parent.children.PushFront(child)
	} else {
		if left {
			child.self = parent.children.InsertAfter(child, newBrother.self)
		} else {
			child.self = parent.children.InsertBefore(child, newBrother.self)
		}
	}

	parent.numOfChildren[child.rank]++
	return parent.value > child.value
}

func (this *node) swapBrothers(other *node) {
	brother := func() *node {
		if other.leftBrother().rank == other.rank {
			return other.leftBrother()
		} else {
			return other.rightBrother()
		}
	}()
	this.parent.addBrother(brother, this, true)
	other.parent.addBrother(this, other, true)
}

func (parent *node) mbyUpdateRank() {
	if parent.children.Len() == 0 {
		parent.rank = 0
	} else {
		parent.rank = parent.children.Front().Value.(*node).rank + 1
		if parent.rank > len(parent.numOfChildren) {
			parent.numOfChildren = append(parent.numOfChildren, 0)
		}
	}
}

func (parent *node) incRank() {
	parent.rank++
	if len(parent.numOfChildren) < int(parent.rank) {
		parent.numOfChildren = append(parent.numOfChildren, 0)
	}
}

func (node *node) link(xNode *node, yNode *node) {

	if node.rank != xNode.rank || node.rank != yNode.rank {
		panic(fmt.Sprintf("Node ranks don't match, %d, %d, %d", node.rank, xNode.rank, yNode.rank))
	}

	node.incRank()
	node.pushFrontChild(xNode, node.leftSon())
	node.pushBackChild(yNode, node.leftSon())
	node.mbyUpdateRank()

	if node.parent != nil {
		node.parent.mbyUpdateRank()
		node.parent.numOfChildren[node.rank]++
		node.parent.numOfChildren[node.rank-1]--
	}
}

func (parent *node) delink() []*node {
	node1 := parent.removeFirstChild()
	node2 := parent.removeFirstChild()
	if node1 == node2 {
		panic(fmt.Sprintf("node1 i node2 su jednaki; %d rang n1, %d rang roditelja", node1.rank, parent.rank))
	}
	if parent.rank > 0 {
		if parent.numOfChildren[parent.rank-1] == 1 {
			node3 := parent.removeFirstChild()
			return []*node{node1, node2, node3}
		}
	}
	return []*node{node1, node2}
}

func (this *node) removeSelfFromViolating() {
	if this.violatingSelf != nil {
		this.parentViolatingList.Remove(this.violatingSelf)
		this.violatingSelf = nil
		this.parentViolatingList = nil
	}
}
