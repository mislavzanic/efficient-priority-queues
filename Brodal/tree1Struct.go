package Brodal

import (
	"errors"
	"fmt"
)

type tree1Struct struct {
	tree1           *tree
	numOfNodesInT1W []int
	rankPointersT1W []*node
	t1GuideW        *guide
	lastKnownRank   int
}

func newT1S(value ValType) *tree1Struct {
	return &tree1Struct{
		tree1:           newTree(value, 1),
		numOfNodesInT1W: []int{},
		rankPointersT1W: []*node{},
		t1GuideW:        newGuide(6),
		lastKnownRank:   0,
	}
}

func newEmptyT1S() *tree1Struct {
	return &tree1Struct{
		tree1:           nil,
		numOfNodesInT1W: []int{},
		rankPointersT1W: []*node{},
		t1GuideW:        newGuide(6),
		lastKnownRank:   0,
	}
}

func (this *tree1Struct) GetTree() (*tree) {
	if this.tree1.root.rank != this.lastKnownRank {
		this.lastKnownRank = this.tree1.root.rank
		this.Update()
	}
	return this.tree1
}

func (this *tree1Struct) GetWGuide() *guide {
	if this.tree1.root.rank != this.lastKnownRank {
		this.lastKnownRank = this.tree1.root.rank
		this.Update()
	}
	return this.t1GuideW
}

func (this *tree1Struct) GetWPointers(index int) *node {
	return this.rankPointersT1W[index]
}

func (this *tree1Struct) GetWNums(index int) int {
	return this.numOfNodesInT1W[index]
}

func (this *tree1Struct) Update() {
	for this.tree1.RootRank() > len(this.numOfNodesInT1W) {
		this.numOfNodesInT1W = append(this.numOfNodesInT1W, 0)
	}

	for this.tree1.RootRank() > len(this.rankPointersT1W) {
		this.rankPointersT1W = append(this.rankPointersT1W, nil)
	}

	this.t1GuideW.expand(this.tree1.RootRank(), 0)
	this.t1GuideW.remove(this.tree1.RootRank())
}

func (this *tree1Struct) insertNewW(child *node) error {
	if this.numOfNodesInT1W[child.rank] == 0 {
		this.rankPointersT1W[child.rank] = child
		child.violatingSelf = this.tree1.vList().PushFront(child)
	} else {
		if this.rankPointersT1W[child.rank].violatingSelf == nil {
			errorStr := "Child with rank %d, value %f has violating self == nil"
			return errors.New(fmt.Sprintf(errorStr, this.rankPointersT1W[child.rank].rank, this.rankPointersT1W[child.rank].value))
		}
		child.violatingSelf = this.tree1.vList().InsertAfter(child, this.rankPointersT1W[child.rank].violatingSelf)
	}
	this.numOfNodesInT1W[child.rank]++
	child.parentViolatingList = this.tree1.wList()
	return nil
}

func (this *tree1Struct) removeFromW(child *node) error {
	// moguci bug -> izbacivanje cvora iz W jer mu se rank povecava -> rank mu je >= ranka korijena pa se trigera error ispod
	if child.rank >= len(this.numOfNodesInT1W) || child.rank >= len(this.rankPointersT1W) {

		// errMessage := "Rank of a node and lengths of lists of W set don't match, child rank: %d, child value: %f"
		// return errors.New(fmt.Sprintf(errMessage, child.rank, child.value))

	} else if this.rankPointersT1W[child.rank] == child {
		this.rankPointersT1W[child.rank] = nil
		if  this.tree1.wList().Back() != child.violatingSelf && child.violatingSelf.Next().Value.(*node).rank == child.rank {
			this.rankPointersT1W[child.rank] = child.violatingSelf.Next().Value.(*node)
		}
	}
	this.tree1.wList().Remove(child.violatingSelf)
	this.numOfNodesInT1W[child.rank]--
	child.removeSelfFromViolating()
	return nil
}

func (this *tree1Struct) childrenWithParentInW(rank int, parent *node) ([]*node, []*node, error) {
	parentChildren, otherChildren := []*node{}, []*node{}
	if this.GetWNums(rank) != 6 {
		errMessage := "There are %d of W violations of rank %d, not 6"
		return nil, nil, errors.New(fmt.Sprintf(errMessage, this.GetWNums(rank), rank))
	}
	for e := this.GetWPointers(rank).violatingSelf; e.Value.(*node).rank != rank; e = e.Next() {

		if e.Value.(*node).rank != rank {
			errMessage := "Missmatching ranks: %d requested, %d got"
			return nil, nil, errors.New(fmt.Sprintf(errMessage, rank, e.Value.(*node).rank))
		}

		if e.Value.(*node).parent == parent {
			parentChildren = append(parentChildren, e.Value.(*node))
		} else {
			otherChildren = append(otherChildren, e.Value.(*node))
		}
	}

	return parentChildren, otherChildren, nil
}
