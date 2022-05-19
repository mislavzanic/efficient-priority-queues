// // package Brodal

// import (
// 	"math"
// 	"fmt"
// )

// type tree1Struct struct {
// 	tree1           *tree
// 	numOfNodesInT1W []int
// 	rankPointersT1W []*node
// 	t1GuideW        *guide
// }

// type BrodalHeap struct {
// 	t1s   *tree1Struct
// 	tree2 *tree
// }

// const ALPHA int = 10

// func EmptyHeap() *BrodalHeap {
// 	return &BrodalHeap{
// 		t1s:             nil,
// 		tree2:           nil,
// 	}
// }

// func newT1S(value float64) *tree1Struct {
// 	return &tree1Struct{
// 		tree1:           newTree(value, 1),
// 		numOfNodesInT1W: []int{},
// 		rankPointersT1W: []*node{},
// 		t1GuideW:        newGuide(6),
// 	}
// }

// func NewHeap(value float64) *BrodalHeap {
// 	return &BrodalHeap{
// 		t1s:   newT1S(value),
// 		tree2: nil,
// 	}
// }

// func (bh *BrodalHeap) reset(newMin float64) {
// 	bh.t1s = newT1S(newMin)
// 	bh.tree2 = nil
// }

// func (bh *BrodalHeap) Min() float64 {
// 	return bh.t1s.tree1.root.value
// }

// func (bh *BrodalHeap) DecreaseKey(currKey *node, value float64) {
// 	if value < bh.t1s.tree1.root.value {
// 		currKey.value = bh.t1s.tree1.root.value
// 		bh.t1s.tree1.root.value = currKey.value
// 	} else {
// 		currKey.value = value
// 		bh.mbyAddViolation(currKey)
// 	}
// }

// func (bh *BrodalHeap) Delete(node *node) {
// 	bh.DecreaseKey(node, math.Inf(-1))
// 	bh.DeleteMin()
// }

// func (bh *BrodalHeap) DeleteMin() {
// 	if bh.t1s == nil {
// 		return
// 	}

// 	if bh.tree2 != nil {
// 		if bh.tree2.RootRank() == 0 && bh.t1s.tree1.RootRank() == 0 {
// 			bh.reset(bh.tree2.root.value)
// 			return
// 		} else {
// 			bh.moveT2ToT1()
// 		}
// 	} else {
// 		if bh.t1s.tree1.RootRank() == 0 {
// 			bh.reset(0)
// 			bh.t1s.tree1 = nil
// 			return
// 		}
// 	}


// 	minW := bh.t1s.tree1.root.getMinFromW()
// 	minV := bh.t1s.tree1.root.getMinFromV()
// 	minTree := bh.t1s.tree1.root.getMinFromChildren()
// 	newMin, _, _ := getMinNodeFrom3(minW, minV, minTree)

// 	mbySwap := bh.t1s.tree1.childrenRank[newMin.rank]
// 	if newMin.parent != bh.t1s.tree1.root {
// 		mbySwap.parent = newMin.parent
// 		newMin.parent = bh.t1s.tree1.root

// 		temp := mbySwap.self
// 		mbySwap.self = newMin.self
// 		newMin.self = temp
// 		bh.t1s.tree1.childrenRank[newMin.rank] = newMin
// 	}

// 	bh.cutNodeFromTree(bh.t1s.tree1, newMin)

// 	indepTrees := bh.t1s.tree1.rmRfRoot()
// 	oldV, oldW := bh.t1s.tree1.root.vList, bh.t1s.tree1.root.wList

// 	bh.reset(newMin.value)

// 	for e := indepTrees.Back(); e != nil; {
// 		if e.Value.(*node).rank == 0 {
// 			bh.Insert(e.Value.(*node).value)
// 		} else {
// 			if e.Value.(*node).rank == bh.t1s.tree1.RootRank() {
// 				bh.incTreeRank(bh.t1s.tree1, e.Value.(*node), e.Prev().Value.(*node))
// 				e = e.Prev()
// 			} else {
// 				if bh.t1s.tree1.root.numOfChildren[bh.t1s.tree1.RootRank()-1] < 2 {
// 					bh.t1s.tree1.addRootChild(e.Value.(*node))
// 				} else {
// 					bh.insertNodeIntoTree(bh.t1s.tree1, e.Value.(*node))
// 				}
// 			}
// 		}
// 		e = e.Prev()
// 	}

// 	for e := newMin.children.Back(); e != nil; e = e.Prev() {
// 		bh.insertNodeIntoTree(bh.t1s.tree1, e.Value.(*node))
// 	}

// 	for e := oldV.Front(); e != nil; e = e.Next() {
// 		bh.mbyAddViolation(e.Value.(*node))
// 	}

// 	for e := oldW.Front(); e != nil; e = e.Next() {
// 		bh.mbyAddViolation(e.Value.(*node))
// 	}

// 	bh.mbyAddViolation(mbySwap)

// 	for e := bh.t1s.tree1.root.wList.Front(); e != nil; {
// 		if e.Next().Value.(*node).rank == e.Value.(*node).rank {
// 			bh.reduceViolation(e.Value.(*node), e.Next().Value.(*node))
// 		} else {
// 			e = e.Next()
// 		}
// 	}
// }

// func (bh *BrodalHeap) Insert(value float64) {
// 	otherHeap := NewHeap(value)
// 	bh.Meld(otherHeap)
// }

// // treba major refactor -> prvo odrediti koji su novi bh.t1 i bh.t2 i onda ubacivati -> potrebno zbog v i w setova
// func (bh *BrodalHeap) Meld(other *BrodalHeap) {
// 	// trees := [](*tree){bh.tree1, bh.tree2, other.tree1, other.tree2}

// 	minTree, otherMin := getMinValTree(bh.t1s, other.t1s)
// 	if other.tree2 == nil && bh.tree2 == nil {
// 		if bh.t1s == nil {
// 			bh.t1s = other.t1s
// 		} else {
// 			bh.t1s = minTree
// 			if minTree.tree1.RootRank() <= otherMin.tree1.RootRank() {
// 				bh.tree2 = otherMin.tree1
// 				bh.tree2.id = 2
// 			} else {
// 				if otherMin.tree1 == bh.t1s.tree1 {
// 					panic("panic")
// 				}
// 				bh.insertNodeIntoTree(bh.t1s.tree1, otherMin.tree1.root)
// 			}
// 		}
// 	} else if bh.t1s == nil {
// 		if bh.tree2 != nil {
// 			panic("Ovo se nebi trebalo dogoditi")
// 		}

// 		bh.t1s = minTree
// 		if other.tree2 != nil {
// 			bh.tree2 = other.tree2
// 		}
// 	} else {
// 		maxTree, others := getMaxTree([]*tree{otherMin.tree1, bh.tree2, other.tree2})
// 		if minTree.tree1.RootRank() == maxTree.RootRank() {
// 			others = append(others, maxTree)
// 			maxTree = minTree.tree1
// 		}

// 		bh.t1s = minTree
// 		if maxTree != minTree.tree1 {
// 			bh.tree2 = maxTree
// 			bh.tree2.id = 2
// 		} else {
// 			bh.tree2 = nil
// 		}


// 		if maxTree.RootRank() == 0 {
// 			bh.incTreeRank(maxTree, others[0].root, others[1].root)
// 			if len(others) == 3 {
// 				bh.insertNodeIntoTree(maxTree, others[2].root)
// 			}
// 		} else {
// 			for _, tree := range others {
// 				if tree != minTree.tree1 {
// 					for maxTree.root.rank == tree.root.rank {
// 						nodes := tree.delink()
// 						for _, n := range nodes {
// 							bh.insertNodeIntoTree(maxTree, n)
// 						}
// 					}
// 					bh.insertNodeIntoTree(maxTree, tree.root)
// 				}
// 			}
// 		}
// 	}
// 	if bh.tree2 != nil {
// 		for e := bh.tree2.root.children.Front(); e != nil; e = e.Next() {
// 			if e.Value.(*node).parent != bh.tree2.root {
// 				panic("tree2 tusam")
// 				if e.Value.(*node).parent == bh.t1s.tree1.root {
// 					panic("t1s")
// 				}
// 			}
// 		}
// 	}
// }

// func (bh *BrodalHeap) mbyRemoveFromViolating(node *node) {
// 	if node.parentViolatingList != nil && node.isGood() {
// 		bh.removeFromViolating(node)
// 	}
// }

// func (bh *BrodalHeap) removeFromViolating(notBad *node) {
// 	if notBad.parentViolatingList == bh.t1s.tree1.root.wList {
// 		// println("rootrank", bh.t1s.tree1.root.rank, notBad.rank)
// 		bh.t1s.numOfNodesInT1W[notBad.rank]--
// 		if bh.t1s.rankPointersT1W[notBad.rank] == notBad {
// 			bh.t1s.rankPointersT1W[notBad.rank] = nil
// 			if notBad.violatingSelf != bh.t1s.tree1.root.wList.Back() {
// 				if notBad.violatingSelf.Next().Value.(*node).rank == notBad.rank {
// 					bh.t1s.rankPointersT1W[notBad.rank] = notBad.violatingSelf.Next().Value.(*node)
// 				}
// 			}
// 		}
// 	}
// 	notBad.removeSelfFromViolating()
// }

// func (bh *BrodalHeap) mbyAddViolation(mbyBad *node) {
// 	if mbyBad.parent == nil {
// 		return
// 	}
// 	if !mbyBad.isGood() && mbyBad.violatingSelf == nil {
// 		bh.addViolation(mbyBad)
// 	}
// }

// func (bh *BrodalHeap) addViolation(bad *node) {
// 	if bad.rank >= bh.t1s.tree1.RootRank() {
// 		bad.violatingSelf = bh.t1s.tree1.root.vList.PushFront(bad)
// 		bad.parentViolatingList = bh.t1s.tree1.root.vList
// 		bh.updateVSet(bad)
// 	} else {
// 		bh.updateWSet(bad)
// 	}
// }

// func (bh *BrodalHeap) updateWSet(bad *node) {

// 	acts := bh.t1s.t1GuideW.forceIncrease(bad.rank, bh.t1s.numOfNodesInT1W[bad.rank]+1, 2)

// 	for _, act := range acts {
// 		if act.op == Increase {
// 			if len(bh.t1s.rankPointersT1W) <= bad.rank {
// 				for len(bh.t1s.rankPointersT1W) <= bad.rank {
// 					bh.t1s.rankPointersT1W = append(bh.t1s.rankPointersT1W, nil)
// 				}
// 				bad.violatingSelf = bh.t1s.tree1.root.wList.PushFront(bad)
// 				bh.t1s.rankPointersT1W[bad.rank] = bad
// 			} else if bh.t1s.rankPointersT1W[bad.rank] == nil {
// 				bh.t1s.rankPointersT1W[bad.rank] = bad
// 				bad.violatingSelf = bh.t1s.tree1.root.wList.PushFront(bad)
// 			} else {
// 				bad.violatingSelf = bh.t1s.tree1.root.wList.InsertAfter(bad, bh.t1s.rankPointersT1W[bad.rank].violatingSelf)
// 			}

// 			for len(bh.t1s.numOfNodesInT1W) <= bad.rank {
// 				bh.t1s.numOfNodesInT1W = append(bh.t1s.numOfNodesInT1W, 0)
// 			}

// 			bh.t1s.numOfNodesInT1W[bad.rank]++
// 			bad.parentViolatingList = bh.t1s.tree1.root.wList
// 			println("println", bh.t1s.numOfNodesInT1W[bad.rank], bad.rank)
// 		} else {
// 			bh.reduceWViolations(act)
// 		}
// 	}
// }

// func (bh *BrodalHeap) updateVSet(bad *node) {
// 	if bh.t1s.tree1.root.vList.Len() > ALPHA*bh.t1s.tree1.RootRank() {
// 		if bh.tree2 == nil {
// 			panic("This can't happen")
// 		}
// 		if bh.tree2.RootRank() <= bh.t1s.tree1.RootRank()+1 {
// 			for leftmostSon := bh.tree2.LeftmostSon(); bh.t1s.tree1.RootRank() <= bh.tree2.RootRank(); {
// 				bh.mbyRemoveFromViolating(leftmostSon)
// 				bh.tree2.removeRootChild(leftmostSon)
// 				if leftmostSon.rank == bh.t1s.tree1.RootRank() {
// 					nextLeftSon := bh.tree2.LeftmostSon()
// 					bh.cutNodeFromTree(bh.tree2, nextLeftSon)
// 					bh.incTreeRank(bh.t1s.tree1, leftmostSon, nextLeftSon)
// 				} else {
// 					if leftmostSon.rank >= bh.t1s.tree1.RootRank() {
// 						panic(fmt.Sprintf("Nije dobro, rang lijevog dijeteta je %d, a korijena je %d", leftmostSon.rank, bh.t1s.tree1.RootRank()))
// 					}
// 					if bh.tree2.childrenRank[leftmostSon.rank] == leftmostSon {
// 						panic("panic")
// 					}
// 					bh.insertNodeIntoTree(bh.t1s.tree1, leftmostSon)
// 				}
// 				leftmostSon = bh.tree2.LeftmostSon()
// 			}
// 			bh.insertNodeIntoTree(bh.t1s.tree1, bh.tree2.root)
// 			bh.tree2 = nil
// 		} else {
// 			// micanje djeteta -> guide se ne azurira
// 			r := bh.t1s.tree1.RootRank()+1
// 			c := bh.tree2.childrenRank[bh.t1s.tree1.RootRank()+1]
// 			bh.derankAndInsertRootChild(bh.tree2, bh.t1s.tree1, bh.tree2.childrenRank[bh.t1s.tree1.RootRank()+1])
// 			if bh.tree2.childrenRank[r] == c {
// 				println(c.rank)
// 				panic("panic")
// 			}
// 		}
// 	}
// }

// func (bh *BrodalHeap) reduceWViolations(act action) {
// 	numOfSonsOfT2 := 0
// 	notSonsOfT2 := []*node{}

// 	println("act", act.index, "len", len(bh.t1s.rankPointersT1W), bh.t1s.numOfNodesInT1W[act.index])

// 	for e := bh.t1s.rankPointersT1W[act.index].violatingSelf; e != nil; e = e.Next() {
// 		println("rang w", e.Value.(*node).rank)
// 	}
// 	for e := bh.t1s.rankPointersT1W[act.index].violatingSelf; e.Value.(*node).rank != act.index; e = e.Next() {
// 		println("jedansin")
// 		if e.Value.(*node).parent == bh.tree2.root {
// 			numOfSonsOfT2++
// 		} else {
// 			notSonsOfT2 = append(notSonsOfT2, e.Value.(*node))
// 		}
// 	}

// 	if numOfSonsOfT2 > 4 {
// 		numOfRemoved := 0
// 		for _, rmNode := range notSonsOfT2 {

// 			bh.removeViolatingNode(rmNode, nil)
// 			numOfRemoved++
// 		}

// 		for e := bh.t1s.rankPointersT1W[act.index].violatingSelf; e.Value.(*node).rank != act.index && numOfRemoved < 2; {
// 			bh.t1s.rankPointersT1W[act.index] = e.Next().Value.(*node)

// 			e.Value.(*node).removeSelfFromViolating()
// 			bh.cutNodeFromTree(bh.tree2, e.Value.(*node))

// 			bh.insertNodeIntoTree(bh.t1s.tree1, e.Value.(*node))
// 			numOfRemoved++
// 		}
// 	} else {
// 		println("novibug", numOfSonsOfT2)
// 		bh.reduceViolation(notSonsOfT2[0], notSonsOfT2[1])

// 		notGood := func() *node {
// 			if !notSonsOfT2[0].isGood() {
// 				return notSonsOfT2[0]
// 			} else if !notSonsOfT2[1].isGood() {
// 				return notSonsOfT2[1]
// 			} else {
// 				return nil
// 			}
// 		}()

// 		if notGood != nil {
// 			bh.removeViolatingNode(notGood, nil)
// 		}
// 	}

// 	bh.t1s.numOfNodesInT1W[act.index] -= 2
// }

// func (bh *BrodalHeap) insertNodeIntoTree(tree *tree, node *node) {
// 	if node.rank < tree.RootRank()-2 {
// 		bh.updateLowRank(tree, node, true) // tu se nes sjebe
// 	} else {
// 		tree.addRootChild(node)
// 		bh.mbyAddViolation(node)
// 	}
// 	bh.updateHighRank(tree, tree.root.rank-2)
// 	bh.updateHighRank(tree, tree.root.rank-1)
// }

// func (bh *BrodalHeap) cutNodeFromTree(tree *tree, node *node) {
// 	if node.rank < tree.RootRank()-2 {
// 		println("tusamtusam")
// 		bh.updateLowRank(tree, node, false)
// 		if node.parent != nil {
// 			panic("nije nil")
// 		}
// 		println("tusamtusam")
// 	} else {
// 		bh.mbyRemoveFromViolating(node)
// 		tree.removeRootChild(node)
// 	}
// 	bh.updateHighRank(tree, tree.root.rank-2)
// 	bh.updateHighRank(tree, tree.root.rank-1)
// }

// func (bh *BrodalHeap) derankAndInsertRootChild(removeRoot *tree, insertRoot *tree, child *node) {
// 	rank := child.rank
// 	bh.cutNodeFromTree(removeRoot, child)

// 	if child == removeRoot.childrenRank[child.rank] {
// 		println("U paniku", child.rank, removeRoot.RootRank(), removeRoot.id)
// 		panic("panic2 nastavak")
// 	}
// 	if child.rank > 0 {
// 		for child.rank >= rank {
// 			nodes := child.delink()
// 			if nodes[0].rank == insertRoot.RootRank() {
// 				bh.incTreeRank(insertRoot, nodes[0], nodes[1])
// 				nodes = nodes[2:]
// 			}
// 			for _, n := range nodes {
// 				bh.insertNodeIntoTree(insertRoot, n)
// 			}
// 		}
// 		bh.insertNodeIntoTree(insertRoot, child)
// 	}
// }

// func (bh *BrodalHeap) updateLowRank(tree *tree, node *node, insert bool) {

// 	responseInc, responseDec := tree.askGuide(node.rank, tree.root.numOfChildren[node.rank], insert)

// 	println("duljina", len(responseInc))
// 	for _, act := range responseInc {
// 		println("act", act.op, insert)
// 		if insert {
// 			if act.op == Increase {
// 				tree.addRootChild(node)
// 				bh.mbyAddViolation(node)
// 			} else {
// 				if tree.root.numOfChildren[act.index] != 7 {
// 					panic("Nije dobro, ne smije biti manje od 7")
// 				}
// 				tree.link(act.index)
// 			}
// 		} else {
// 			println("tusamtusam")
// 			if act.op == Increase {
// 				println("tusam")
// 				bh.mbyRemoveFromViolating(node)
// 				tree.removeRootChild(node)
// 			} else {
// 				bh.derankAndInsertRootChild(tree, tree, tree.childrenRank[act.index+1])
// 			}
// 		}
// 	}

// 	for _, act := range responseDec {
// 		if insert {
// 			if act.op == Increase {
// 				// bh.mbyRemoveFromViolating(node)
// 				// tree.removeRootChild(node)
// 			} else {
// 				bh.derankAndInsertRootChild(tree, tree, tree.childrenRank[act.index+1])
// 			}
// 		} else {
// 			if act.op == Increase {
// 				// tree.addRootChild(node)
// 				// bh.mbyAddViolation(node)
// 			} else {
// 				tree.link(act.index)
// 			}
// 		}
// 	}
// }

// func (bh *BrodalHeap) updateHighRank(tree *tree, rank int) {
// 	if rank < 0 {
// 		return
// 	}
// 	if tree.root.numOfChildren[rank] > 7 {
// 		if rank == tree.root.rank-1 {
// 			nodeSliceX := tree.delink()
// 			nodeSliceY := tree.delink()
// 			nodeSliceZ := tree.delink()

// 			minNode1, nodeX1, nodeY1 := getMinNodeFrom3(nodeSliceX[0], nodeSliceX[1], nodeSliceY[0])
// 			minNode2, nodeX2, nodeY2 := getMinNodeFrom3(nodeSliceY[1], nodeSliceZ[1], nodeSliceZ[0])

// 			if minNode1 == nodeX1 || minNode1 == nodeY1 {
// 				panic("jednaki su")
// 			}

// 			bh.removeFromViolating(minNode1)
// 			bh.removeFromViolating(minNode2)
// 			bh.removeFromViolating(nodeX1)
// 			bh.removeFromViolating(nodeX2)
// 			bh.removeFromViolating(nodeY1)
// 			bh.removeFromViolating(nodeY2)

// 			minNode1.link(nodeX1, nodeY1)
// 			minNode2.link(nodeX2, nodeY2)

// 			bh.incTreeRank(tree, minNode1, minNode2)
// 		} else {
// 			tree1 := tree.removeRootChild(tree.childrenRank[rank])
// 			tree2 := tree.removeRootChild(tree.childrenRank[rank])
// 			tree3 := tree.removeRootChild(tree.childrenRank[rank])

// 			bh.removeFromViolating(tree1)
// 			bh.removeFromViolating(tree2)
// 			bh.removeFromViolating(tree3)

// 			minTree, tree1, tree2 := getMinNodeFrom3(tree1, tree2, tree3)
// 			minTree.link(tree1, tree2)
// 			bh.insertNodeIntoTree(tree, minTree)
// 			bh.updateHighRank(tree, rank+1)
// 		}
// 	} else if tree.root.numOfChildren[rank] < 2 {
// 		if rank == tree.RootRank()-2 {
// 			takeChildrenFrom := tree.childrenRank[rank+1]

// 			if takeChildrenFrom.numOfChildren[rank] <= 3 {
// 				bh.mbyRemoveFromViolating(takeChildrenFrom)
// 				tree.removeRootChild(takeChildrenFrom)
// 			}

// 			nodes := takeChildrenFrom.delink()
// 			for _, n := range nodes {
// 				bh.insertNodeIntoTree(tree, n)
// 			}

// 			bh.updateHighRank(tree, rank+1)

// 		} else {
// 			// bh.derankAndInsertRootChild(tree, tree, tree.childrenRank[rank])
// 			derankChild := tree.childrenRank[rank]
// 			bh.removeFromViolating(tree.childrenRank[rank])
// 			tree.removeRootChild(derankChild)

// 			if derankChild.rank > 0 {
// 				for derankChild.rank >= rank {
// 					nodeSlice := derankChild.delink()
// 					if nodeSlice[0].rank == tree.RootRank() {
// 						bh.incTreeRank(tree, nodeSlice[0], nodeSlice[1])
// 						nodeSlice = nodeSlice[2:]
// 					}
// 					for _, n := range nodeSlice {
// 						bh.insertNodeIntoTree(tree, n)
// 					}
// 				}
// 				bh.insertNodeIntoTree(tree, derankChild)
// 			} else {
// 				panic("ne znam vise")
// 			}
// 		}
// 	}
// }

// func (bh *BrodalHeap) incTreeRank(tree *tree, node1 *node, node2 *node) {

// 	if tree.id == 1 {
// 		bh.t1s.rankPointersT1W = append(bh.t1s.rankPointersT1W, nil)
// 		bh.t1s.numOfNodesInT1W = append(bh.t1s.numOfNodesInT1W, 0)
// 		bh.t1s.t1GuideW.expand(tree.RootRank() + 1, 0)
// 	}

// 	bh.removeFromViolating(node1)
// 	bh.removeFromViolating(node2)

// 	tree.incRank(node1, node2)

// 	bh.mbyAddViolation(node1)
// 	bh.mbyAddViolation(node2)
// }

// func (bh *BrodalHeap) reduceViolation(x1 *node, x2 *node) {
// 	if x1.isGood() || x2.isGood() {
// 		bh.mbyRemoveFromViolating(x1)
// 		bh.mbyRemoveFromViolating(x2)
// 	} else {
// 		if x1.parent != x2.parent {
// 			if x1.parent.value <= x2.parent.value {
// 				x1.swapBrothers(x2)
// 			} else {
// 				x2.swapBrothers(x1)
// 			}
// 		}
// 		bh.removeViolatingNode(x1, x2)
// 	}
// }

// func (bh *BrodalHeap) removeViolatingNode(rmNode *node, otherBrother *node) {
// 	println("removeviolaingnode", rmNode.rank, otherBrother.rank, bh.t1s.tree1.root.rank)
// 	parent := rmNode.parent
// 	replacement := bh.t1s.tree1.childrenRank[parent.rank]
// 	grandParent := parent.parent
// 	if otherBrother == nil {
// 		otherBrother = func() *node {
// 			if rmNode.leftBrother().rank != rmNode.rank {
// 				return rmNode.rightBrother()
// 			} else {
// 				return rmNode.leftBrother()
// 			}
// 		}()
// 	}

// 	if parent.numOfChildren[rmNode.rank] == 2 {
// 		if parent.rank == rmNode.rank+1 {
// 			if grandParent != bh.t1s.tree1.root {
// 				bh.cutNodeFromTree(bh.t1s.tree1, replacement)
// 				grandParent.pushBackChild(replacement, parent)
// 				bh.mbyAddViolation(replacement)
// 				grandParent.removeChild(parent)
// 			} else {
// 				bh.cutNodeFromTree(bh.t1s.tree1, parent)
// 			}

// 			parent.removeChild(rmNode)
// 			parent.removeChild(otherBrother)
// 			bh.insertNodeIntoTree(bh.t1s.tree1, parent)
// 		} else {
// 			parent.removeChild(rmNode)
// 			parent.removeChild(otherBrother)
// 		}
// 		bh.mbyRemoveFromViolating(otherBrother)

// 		bh.insertNodeIntoTree(bh.t1s.tree1, otherBrother)
// 		bh.insertNodeIntoTree(bh.t1s.tree1, rmNode)
// 	} else {
// 		parent.removeChild(rmNode)
// 		bh.insertNodeIntoTree(bh.t1s.tree1, rmNode)
// 	}
// 	bh.mbyRemoveFromViolating(rmNode)
// }

// func (bh *BrodalHeap) moveT2ToT1() {
// 	child := bh.tree2.Children().Back()
// 	for bh.tree2.Children().Len() != 0 {
// 		bh.tree2.removeRootChild(child.Value.(*node))

// 		if child.Value.(*node).rank == bh.t1s.tree1.RootRank() {
// 			nextChild := bh.tree2.Children().Back().Value.(*node)
// 			bh.tree2.removeRootChild(nextChild)
// 			bh.incTreeRank(bh.t1s.tree1, child.Value.(*node), nextChild)
// 		} else {
// 			bh.insertNodeIntoTree(bh.t1s.tree1, child.Value.(*node))
// 		}
// 		child = bh.tree2.Children().Back()
// 	}
// 	bh.insertNodeIntoTree(bh.t1s.tree1, bh.tree2.root)
// 	bh.tree2 = nil
// }
