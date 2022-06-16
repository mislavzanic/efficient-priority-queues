package Brodal

import "container/list"


func getMinValTree[T BoNumber](t1s1 *tree1Struct[T], t1s2 *tree1Struct[T]) (*tree1Struct[T], *tree1Struct[T]) {
	if t1s1 == nil || t1s2 == nil {
		panic("One of the trees is nil")
	}

	if t1s1.tree1.root.value > t1s2.tree1.root.value {
		return t1s2, t1s1
	} else {
		return t1s1, t1s2
	}
}

func getMaxTree[T BoNumber](trees ...*tree[T]) (*tree[T], []*tree[T]) {
	if len(trees) == 0 {
		panic("There are no trees")
	}

	maxTree := trees[0]
	maxTreeIndex := 0
	newTrees := [](*tree[T]){}

	for i, tree := range trees {
		if tree != nil {
			if maxTree.root.rank < tree.root.rank {
				newTrees = append(newTrees, maxTree)
				maxTree = tree
				maxTreeIndex = i
			} else {
				newTrees = append(newTrees, tree)
			}
		}
	}

	return maxTree, append(newTrees[:maxTreeIndex],  newTrees[maxTreeIndex+1:]...)
}

func mbySwapTree[T BoNumber](ptr1 *tree[T], ptr2 *tree[T], cond bool) (*tree[T], *tree[T]) {
	if cond {
		temp := ptr1
		ptr1 = ptr2
		ptr2 = temp
	}
	return ptr1, ptr2
}

func getMinNode[T BoNumber](xNode *node[T], yNode *node[T]) (*node[T], *node[T]) {
	if xNode == nil { return yNode, xNode }
	if yNode == nil { return xNode, yNode }

	if xNode.value < yNode.value {
		return xNode, yNode
	}
	return yNode, xNode
}

func getMinNodeFrom3[T BoNumber](xNode *node[T], yNode *node[T], zNode *node[T]) (*node[T], *node[T], *node[T]) {
	minNode, otherNode := getMinNode(xNode, yNode)
	minNode, otherOtherNode := getMinNode(minNode, zNode)
	return minNode, otherNode, otherOtherNode
}

func mbySwapNode[T BoNumber](ptr1 *node[T], ptr2 *node[T], cond bool) {
	if cond {
		temp := *ptr1
		*ptr1 = *ptr2
		*ptr2 = temp
	}
}
func getMinFromList[T BoNumber](list *list.List) *list.Element {
	minVal := list.Front()
	for e := list.Front(); e != nil; e = e.Next() {
		if e.Value.(*node[T]).value < minVal.Value.(*node[T]).value {
			minVal = e
		}
	}

	return minVal
}
