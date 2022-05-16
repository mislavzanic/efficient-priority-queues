package Brodal

import (
	"container/list"
)

type node struct {
	value float64
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
	return getMinFromList(this.wList).Value.(*node)
}

func (this *node) getMinFromV() *node {
	return getMinFromList(this.vList).Value.(*node)
}

func (this *node) getMinFromChildren() *node {
	return getMinFromList(this.children).Value.(*node)
}

func (node *node) isGood() bool {
	return node.value > node.parent.value
}

func (parent *node) removeChild(child *node) int {
	parent.children.Remove(child.self)
	parent.numOfChildren[child.rank]--

	parent.mbyUpdateRank()

	return parent.numOfChildren[child.rank]
}

func (parent *node) removeFirstChild() (*node, int) {
	child := parent.children.Front().Value.(*node)

	if parent.rank - 1 != child.rank {
		panic("Incorrect ranks")
	}

	numOfChildren := parent.removeChild(child)
	return child, numOfChildren
}

func (parent *node) addFirstChildren(child1 *node, child2 *node) {
	child1.self = parent.children.PushBack(child1)
	child1.parent = parent

	if len(parent.numOfChildren) == int(child1.rank) {
		parent.numOfChildren = append(parent.numOfChildren, 1)
	} else {
		parent.numOfChildren[child1.rank] += 1
	}

	parent.mbyUpdateRank()
	// parent.rank++
	parent.addBrother(child2, child1, true)
}

func (parent *node) pushBackChild(child *node, newBrother *node) {
	parent.addBrother(child, newBrother, true)
}

func (parent *node) pushFrontChild(child *node, newBrother *node) {
	parent.addBrother(child, newBrother, false)
}

func (parent *node) addBrother(child *node, newBrother *node, left bool) {
	if parent.rank == child.rank {
		panic("Increase rank of parent first")
	}

	if child.parent != nil {
		child.parent.removeChild(child)
	}

	child.parent = parent
	if newBrother == nil {
		child.self = parent.children.PushBack(child)
	} else {
		if left {
			child.self = parent.children.InsertAfter(child, newBrother.self)
		} else {
			child.self = parent.children.InsertBefore(child, newBrother.self)
		}
	}

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
		panic("Node ranks don't match")
	}

	node.incRank()
	node.addBrother(xNode, node.leftSon(), false)
	node.addBrother(yNode, node.leftSon(), true)
	node.mbyUpdateRank()

	if node.parent != nil {
		node.parent.mbyUpdateRank()
		node.parent.numOfChildren[node.rank]++
		node.parent.numOfChildren[node.rank - 1]--
	}
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

func (this *node) removeSelfFromViolating() {
	this.parentViolatingList.Remove(this.violatingSelf)
	this.parentViolatingList = nil
}
