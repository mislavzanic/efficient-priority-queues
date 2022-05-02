package main

import (
	"fmt"
)

type BrodalHeap struct {
	root *Node
}

type BrodalQueue struct {
    T1 *BrodalHeap
    T2 *BrodalHeap
}
