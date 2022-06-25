package main

type Number interface {
	int | int8 | int16 | int32 | int64 | float32 | float64
}

type PriorityQueue[T Number] interface {
	Meld(interface{})
	Min() T
	DeleteMin() T
	Insert(T)
	Empty() bool
}
