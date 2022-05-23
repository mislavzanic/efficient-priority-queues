package Brodal

import (
	"container/list"
	"errors"
	"fmt"
)

type node struct {
	value ValType
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

func newNode(value ValType) *node {
	node := new(node)
	node.value = value
	node.rank = 0
	node.children = list.New()
	node.vList = list.New()
	node.wList = list.New()

	return node
}

func (this *node) ToString() string {
	str := "\t\tValue:"
	str += fmt.Sprint(this.value)
	str += "\n\t\tRank:"
	str += fmt.Sprint(this.rank)
	if this.parent != nil {
		str += "\n\t\tParent:"
		str += fmt.Sprint(this.parent.value)
	}
		str += "\n\t\tChildren:\n"
	for e := this.children.Front(); e != nil; e = e.Next() {
		str += e.Value.(*node).ToString()
	}
	return str + "\n"
}

func (this *node) Leaf() bool {
	return this.children == nil || this.children.Len() == 0
}

func (parent *node) LeftChild() *node {
	if n, err := parent.leftChild(); err == nil {
		return n
	} else {
		panic(fmt.Sprint("%w", err))
	}
}

func (parent *node) leftChild() (*node,error) {
	if parent.Leaf() {
		return nil, errors.New("Node doesn't have any children")
	}
	return parent.children.Front().Value.(*node), nil
}

func (this *node) LeftBrother() *node {
	if n, err := this.leftBrother(); err == nil {
		return n
	} else {
		return nil
		// panic(fmt.Sprint("%w", err))
	}
}

func (this *node) leftBrother() (*node,error) {
	if this.self.Prev() == nil {
		return nil, errors.New(fmt.Sprintf("Node with rank = %d, value = %f, doesn't have a left brother", this.rank, this.value))
	}
	return this.self.Prev().Value.(*node), nil
}

func (this *node) RightBrother() *node {
	if n, err := this.rightBrother(); err == nil {
		return n
	} else {
		return nil
		// panic(fmt.Sprint("%w", err))
	}
}

func (this *node) rightBrother() (*node, error) {
	if this.self.Next() == nil {
		return nil, errors.New(fmt.Sprintf("Node with rank = %d, value = %f, doesn't have a left brother", this.rank, this.value))
	}
	return this.self.Next().Value.(*node), nil
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
	if node.parent == nil {
		return true
	}
	return node.value >= node.parent.value
}

func (parent *node) removeChild(child *node) (*node, error) {
	if value := parent.children.Remove(child.self); value != nil {
		child.parent = nil
		parent.numOfChildren[child.rank]--
		parent.mbyUpdateRank()
		return child, nil
	}

	errorMsg := "Child with rank %d, value %f was not in child list of node with rank %d, value %f"
	return nil, errors.New(fmt.Sprintf(errorMsg, child.rank, child.value, parent.rank, parent.value))
}

func (parent *node) removeFirstChild() (*node, error) {
	child := parent.children.Front().Value.(*node)

	if parent.rank-1 != child.rank {
		return nil, errors.New(fmt.Sprintf("Parent rank %d and child rank %d don't match", parent.rank, child.rank))
	}

	return parent.removeChild(child)
}

func (parent *node) pushBackChild(child *node, newBrother *node) (bool, error) {
	return parent.addBrother(child, newBrother, true)
}

func (parent *node) pushFrontChild(child *node, newBrother *node) (bool, error) {
	return parent.addBrother(child, newBrother, false)
}

func (parent *node) addBrother(child *node, newBrother *node, left bool) (bool, error) {
	if parent.rank == child.rank {
		return false, errors.New("Increase rank of parent first")
	}

	if newBrother != nil {
		if newBrother.parent != parent {
			return false, errors.New(fmt.Sprintf("newBrother parent is not the parent, %f", parent.value))
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
	return parent.value > child.value, nil
}

func (this *node) swapBrothers(other *node) error {
	var brother *node = nil

	if err1 := other.setIfNoErrors(brother, other.leftBrother); err1 != nil {
		if err2 := other.setIfNoErrors(brother, other.rightBrother); err2 != nil {
			return errors.New(err1.Error() + err2.Error())
		}
	}

	if _, err := this.parent.addBrother(brother, this, true); err != nil {
		return err
	}

	if _, err := other.parent.addBrother(this, other, true); err != nil {
		return err
	}

	return nil
}

func (parent *node) mbyUpdateRank() {
	if parent.children.Len() == 0 {
		parent.rank = 0
	} else {
		parent.rank = parent.children.Front().Value.(*node).rank + 1
		for parent.rank > len(parent.numOfChildren) {
			parent.numOfChildren = append(parent.numOfChildren, 0)
		}
	}
}

func (parent *node) incRank() {
	parent.rank++
	// update ako je cvor u nekoj W listi --> potrebne zasebne klase za V i W !!!!
	for len(parent.numOfChildren) < int(parent.rank) {
		parent.numOfChildren = append(parent.numOfChildren, 0)
	}
}

func (this *node) decRank() ([]*node, error){

	if this.rank == 0 {
		return nil, errors.New(fmt.Sprintf("Node with value %f is of rank 0", this.value))
	}

	nodes, currRank := []*node{}, this.rank
	for this.rank == currRank {
		if child, err := this.removeChild(this.LeftChild()); err != nil {
			return nil, err
		} else {
			nodes = append(nodes, child)
		}
	}

	return nodes, nil
}

func (this *node) DecRank() []*node {
	if nodes, err := this.decRank(); err != nil {
		panic(err.Error())
	} else {
		return nodes
	}
}

func (node *node) link(xNode *node, yNode *node) (bool, error) {

	if node.rank != xNode.rank || node.rank != yNode.rank {
		return false ,errors.New(fmt.Sprintf("Node ranks don't match, %d, %d, %d", node.rank, xNode.rank, yNode.rank))
	}

	violating := false

	node.incRank()
	if node.Leaf() {
		if v, err := node.pushFrontChild(xNode, nil); err == nil {
			violating = v || violating
		} else {
			return false, err
		}
	} else {
		if lc, err := node.leftChild(); err == nil {
			if v, err := node.pushFrontChild(xNode, lc); err == nil {
				violating = v || violating
			} else {
				return false, err
			}
		} else {
			return false, err
		}
	}

	if v, err := node.pushBackChild(yNode, xNode); err == nil {
		violating = v || violating
	} else {
		return false, err
	}

	if node.parent != nil {
		node.parent.mbyUpdateRank()
		node.parent.numOfChildren[node.rank]++
		node.parent.numOfChildren[node.rank-1]--
	}

	return violating, nil

}

func (parent *node) delink() ([]*node, error) {
	if node1, err := parent.removeFirstChild(); err == nil {
		if node2, err := parent.removeFirstChild(); err == nil {
			if node1 == node2 {
				return nil, errors.New(fmt.Sprintf("node1 i node2 su jednaki; %d rang n1, %d rang roditelja", node1.rank, parent.rank))
			}

			if parent.numOfChildren[node1.rank] == 1 {
				if node3, err := parent.removeFirstChild(); err == nil {
					return []*node{node1, node2, node3}, nil
				} else {
					return nil, err
				}
			}

			return []*node{node1, node2}, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func (this *node) removeSelfFromViolating() {
	if this.parentViolatingList != nil {
		this.parentViolatingList.Remove(this.violatingSelf)
		this.parentViolatingList = nil
	}
	if this.violatingSelf == nil { return }
	this.violatingSelf = nil
}

func (this *node) setIfNoErrors(setThis *node, method func() (*node, error)) error {
	if n, err := method(); err == nil {
		setThis = n
		return err
	} else {
		return err
	}
}
