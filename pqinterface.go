package main

type Number interface {
	int | int8 | int16 | int32 | int64 | float32 | float64
}

type Node[T Number] interface {
	Value() T
	Ident() any
}

type PriorityQueue[T Number] interface {
	Meld(interface{})
	Min() *Node[T]
	DeleteMin() *Node[T]
	Insert(T, any)
	Empty() bool
}
