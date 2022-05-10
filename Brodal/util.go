package Brodal


func getMinTree(tree1 *tree, tree2 *tree) (*tree, *tree) {
	if tree1 == nil || tree2 == nil {
		panic("One of the trees is nil")
	}

	if tree1.root.value > tree2.root.value {
		return tree2, tree1
	} else {
		return tree1, tree2
	}
}

func getMaxTree(trees []*tree) (*tree, []*tree) {
	if len(trees) == 0 {
		panic("There are no trees")
	}

	maxTree := trees[0]
	newTrees := [](*tree){}

	for _, tree := range trees {
		if tree != nil {
			if maxTree.root.rank < tree.root.rank {
				newTrees = append(newTrees, maxTree)
				maxTree = tree
			} else {
				newTrees = append(newTrees, tree)
			}
		}
	}

	return maxTree, newTrees
}

func mbySwapTree(ptr1 *tree, ptr2 *tree, cond bool) (*tree, *tree) {
	if cond {
		temp := ptr1
		ptr1 = ptr2
		ptr2 = temp
	}
	return ptr1, ptr2
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

func mbySwapNode(ptr1 *node, ptr2 *node, cond bool) {
	if cond {
		temp := *ptr1
		*ptr1 = *ptr2
		*ptr2 = temp
	}
}
