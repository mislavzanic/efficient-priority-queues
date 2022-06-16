package Brodal

import (
	"errors"
	"fmt"
)

type tree1Struct[T BoNumber] struct {
	tree1           *tree[T]
	numOfNodesInT1W []int
	rankPointersT1W []*node[T]
	t1GuideW        *guide
	lastKnownRank   int
}

func newT1S[T BoNumber](value T, parent *BrodalHeap[T]) *tree1Struct[T] {
	return &tree1Struct[T]{
		tree1:           newTree(value, 1, parent),
		numOfNodesInT1W: []int{},
		rankPointersT1W: []*node[T]{},
		t1GuideW:        newGuide(6),
		lastKnownRank:   0,
	}
}

func newEmptyT1S[T BoNumber]() *tree1Struct[T] {
	return &tree1Struct[T]{
		tree1:           nil,
		numOfNodesInT1W: []int{},
		rankPointersT1W: []*node[T]{},
		t1GuideW:        newGuide(6),
		lastKnownRank:   0,
	}
}

func (this *tree1Struct[T]) GetTree() (*tree[T]) {
	if this.tree1 != nil && this.tree1.root.rank != this.lastKnownRank {
		this.lastKnownRank = this.tree1.root.rank
		this.Update()
	}
	return this.tree1
}

func (this *tree1Struct[T]) GetWGuide() *guide {
	if this.tree1.root.rank != this.lastKnownRank {
		this.lastKnownRank = this.tree1.root.rank
		this.Update()
	}
	return this.t1GuideW
}

func (this *tree1Struct[T]) GetWPointers(index int) *node[T] {
	return this.rankPointersT1W[index]
}

func (this *tree1Struct[T]) GetWNums(index int) int {
	return this.numOfNodesInT1W[index]
}

func (this *tree1Struct[T]) Update() {
	for this.tree1.RootRank() > len(this.numOfNodesInT1W) {
		this.numOfNodesInT1W = append(this.numOfNodesInT1W, 0)
	}

	for this.tree1.RootRank() > len(this.rankPointersT1W) {
		this.rankPointersT1W = append(this.rankPointersT1W, nil)
	}

	this.t1GuideW.expand(this.tree1.RootRank(), 0)
	this.t1GuideW.remove(this.tree1.RootRank())
}

func (this *tree1Struct[T]) insertNewW(child *node[T]) error {
	if child.parentViolatingList == this.tree1.wList() {
		return nil
	} else if child.parentViolatingList != nil {
		child.removeSelfFromViolating()
	}

	if this.numOfNodesInT1W[child.rank] == 0 {
		child.violatingSelf = this.tree1.wList().PushFront(child)
		this.rankPointersT1W[child.rank] = child
	} else {
		if this.rankPointersT1W[child.rank].violatingSelf == nil {
			errorStr := "Child with rank %d, value %f has violating self == nil"
			return errors.New(fmt.Sprintf(errorStr, this.rankPointersT1W[child.rank].rank, this.rankPointersT1W[child.rank].value))
		}
		child.violatingSelf = this.tree1.wList().InsertAfter(child, this.rankPointersT1W[child.rank].violatingSelf)
	}
	this.numOfNodesInT1W[child.rank]++
	child.parentViolatingList = this.tree1.wList()
	return nil
}

func (this *tree1Struct[T]) removeFromW(child *node[T]) error {
	if child.rank >= len(this.numOfNodesInT1W) || child.rank >= len(this.rankPointersT1W) {

		errMessage := "Rank of a node and lengths of lists of W set don't match, child rank: %d, child value: %f"
		return errors.New(fmt.Sprintf(errMessage, child.rank, child.value))

	} else if this.rankPointersT1W[child.rank] == child {
		this.rankPointersT1W[child.rank] = nil
		if  this.tree1.wList().Back() != child.violatingSelf && child.violatingSelf.Next().Value.(*node[T]).rank == child.rank {
			this.rankPointersT1W[child.rank] = child.violatingSelf.Next().Value.(*node[T])
		}
	}
	this.tree1.root.wList.Remove(child.violatingSelf)
	this.numOfNodesInT1W[child.rank]--
	child.removeSelfFromViolating()
	return nil
}

func (this *tree1Struct[T]) childrenWithParentInW(rank int, parent *node[T]) ([]*node[T], []*node[T], error) {
	parentChildren, otherChildren := []*node[T]{}, []*node[T]{}
	if this.GetWNums(rank) != 6 {
		errMessage := "There are %d of W violations of rank %d, not 6"
		return nil, nil, errors.New(fmt.Sprintf(errMessage, this.GetWNums(rank), rank))
	}
	for e := this.GetWPointers(rank).violatingSelf; e != nil && e.Value.(*node[T]).rank == rank; e = e.Next() {

		if e.Value.(*node[T]).rank != rank {
			errMessage := "Missmatching ranks: %d requested, %d got"
			return nil, nil, errors.New(fmt.Sprintf(errMessage, rank, e.Value.(*node[T]).rank))
		}

		if e.Value.(*node[T]).parent == parent {
			parentChildren = append(parentChildren, e.Value.(*node[T]))
		} else {
			otherChildren = append(otherChildren, e.Value.(*node[T]))
		}
	}

	if len(parentChildren) + len(otherChildren) != 6 {
		f := this.GetWPointers(rank).violatingSelf
		for e := f; e != nil && e.Value.(*node[T]).rank != rank; e = e.Next() {
			println(e.Value.(*node[T]).rank)
		}
		panic("tusam")
	}
	return parentChildren, otherChildren, nil
}
