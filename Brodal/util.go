package Brodal

import "container/list"


func getMinValTree(t1s1 *tree1Struct, t1s2 *tree1Struct) (*tree1Struct, *tree1Struct) {
	if t1s1 == nil || t1s2 == nil {
		panic("One of the trees is nil")
	}

	if t1s1.tree1.root.value > t1s2.tree1.root.value {
		return t1s2, t1s1
	} else {
		return t1s1, t1s2
	}
}

func getMaxTree(trees []*tree) (*tree, []*tree) {
	if len(trees) == 0 {
		panic("There are no trees")
	}

	maxTree := trees[0]
	maxTreeIndex := 0
	newTrees := [](*tree){}

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

func mbySwapTree(ptr1 *tree, ptr2 *tree, cond bool) (*tree, *tree) {
	if cond {
		temp := ptr1
		ptr1 = ptr2
		ptr2 = temp
	}
	return ptr1, ptr2
}

func getMinNode(xNode *node, yNode *node) (*node, *node) {
	if xNode == nil { return yNode, xNode }
	if yNode == nil { return xNode, yNode }

	if xNode.value < yNode.value {
		return xNode, yNode
	}
	return yNode, xNode
}

func getMinNodeFrom3(xNode *node, yNode *node, zNode *node) (*node, *node, *node) {
	minNode, otherNode := getMinNode(xNode, yNode)
	minNode, otherOtherNode := getMinNode(minNode, zNode)
	return minNode, otherNode, otherOtherNode
}

func mbySwapNode(ptr1 *node, ptr2 *node, cond bool) {
	if cond {
		temp := *ptr1
		*ptr1 = *ptr2
		*ptr2 = temp
	}
}
func getMinFromList(list *list.List) *list.Element {
	minVal := list.Front()
	for e := list.Front(); e != nil; e = e.Next() {
		if e.Value.(*node).value < minVal.Value.(*node).value {
			minVal = e
		}
	}

	return minVal
}
