package main

import (
	"fmt"
)

type Node struct {
    value int64
	rank int
	leftBrother     *Node
	rightBrother    *Node
	parent          *Node
	leftmostSon     *Node
	firstV          *Node
	firstW          *Node
	nextInViolating *Node
	prevInViolating *Node
}
