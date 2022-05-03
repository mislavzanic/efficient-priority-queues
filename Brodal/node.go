package main

type node struct {
    value int64
	rank int
	leftBrother     *node
	rightBrother    *node
	parent          *node
	leftmostSon     *node
	firstV          *node
	firstW          *node
	nextInViolating *node
	prevInViolating *node
}
