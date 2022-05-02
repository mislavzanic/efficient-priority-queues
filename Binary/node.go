package Binary

import (
	"fmt"
)

type itemType = int;

type Node struct {
    leftChild *Node
    rightChild *Node
	val itemType
}

func (this *Node) insert (node *Node) {
	if this.val < node.val {
		if this.leftChild == nil {
			this.leftChild = node
		}

		if this.rightChild == nil {
			this.rightChild = node;
		}
	}
}
