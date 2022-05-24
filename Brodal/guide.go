package Brodal

import (
	"fmt"
	"math"
	"strconv"
)

type pair struct {
	fst int // value
	snd int // index
}

type guide struct {
	upperBound int
	boundArray []pair
	blocks     []**pair
}

type operation byte

const (
	Increase operation = 0
	Reduce   operation = 1
	Decrease operation = 2
)

type action struct {
	index int
	op    operation
	value int
	bound int
}

func newGuide(upperBound int) *guide {
	return &guide{
		upperBound: upperBound,
		boundArray: []pair{},
		blocks:     []**pair{},
	}
}

func (guide *guide) forceIncrease(index int, valArray *[]int, reduceValue int) []action {
	ops := []action{}

	actualValue := (*valArray)[index] + 1
	if guide.upperBound < 0 {
		actualValue = -(*valArray)[index] + 1
	}
	if actualValue > guide.boundArray[index].fst {

		if guide.boundArray[index].fst == guide.upperBound {
			guide.fixUp(&guide.boundArray[index], (*valArray)[index], reduceValue, &ops)
		}

		guide.increase(index, &ops)

		if guide.boundArray[index].fst == guide.upperBound-1 && (*guide.blocks[index]) != nil {
			guide.fixUp(*guide.blocks[index], (*valArray)[index], reduceValue, &ops)
		} else if guide.boundArray[index].fst == guide.upperBound {
			if (*guide.blocks[index]) == nil {
				guide.fixUp(&guide.boundArray[index], (*valArray)[index], reduceValue, &ops)
			} else {
				guide.fixUp(*guide.blocks[index], (*valArray)[index], reduceValue, &ops)
				guide.fixUp(&guide.boundArray[index], (*valArray)[index], reduceValue, &ops)
			}
		}
	} else {
		ops = append(ops, action{index, Increase, 1, guide.upperBound})
	}

	return ops
}

// func (guide *guide) forceDecrease(index int, actualValue int, reduceValue int) []action {
// 	ops := []action{}
// 	if actualValue > guide.upperBound - 2 {
// 		guide.decrease(index, nil)
// 		if guide.boundArray[index].fst == guide.upperBound-1 {
// 			(*guide.blocks[index]) = nil
// 		} else {
// 			if (*guide.blocks[index]) != nil {
// 				guide.fixUp(*guide.blocks[index], reduceValue, &ops)
// 			}
// 		}
// 	}
// 	return ops
// }

func (guide *guide) fixUp(pair *pair, actualValue int, reduceValue int, ops *[]action) {
	(*guide.blocks[pair.snd]) = nil

	if actualValue != guide.upperBound {
		guide.boundArray[pair.snd].fst = int(math.Max(float64(actualValue), float64(guide.upperBound - 2)))
		return
	}

	guide.reduce(pair.snd, reduceValue, ops)

	if pair.snd != len(guide.boundArray)-1 {

		if guide.boundArray[pair.snd+1].fst == guide.upperBound-1 {
			if (*guide.blocks[pair.snd+1]) != nil {
				guide.blocks[pair.snd] = guide.blocks[pair.snd+1]
			}
		} else if guide.boundArray[pair.snd+1].fst == guide.upperBound {
			ptr := &guide.boundArray[pair.snd+1]
			guide.blocks[pair.snd+1] = &ptr
			guide.blocks[pair.snd] = &ptr
		}
	}
}

func (guide *guide) increase(index int, ops *[]action) {
	if index < len(guide.boundArray) {
		guide.boundArray[index].fst++
		if ops != nil {
			*ops = append(*ops, action{index, Increase, 1, guide.upperBound})
		}
	}
}

func (guide *guide) decrease(index int, ops *[]action) {
	if index < len(guide.boundArray) {
		guide.boundArray[index].fst--
		if ops != nil {
			*ops = append(*ops, action{index, Decrease, 1, guide.upperBound})
		}
	}
}

func (guide *guide) reduce(index int, reduceValue int, ops *[]action) {
	if ops != nil {
		*ops = append(*ops, action{index, Reduce, reduceValue, guide.upperBound})
	}

	if guide.boundArray[index].fst != guide.upperBound {
		panic("Ipak je moguce... ja sam u guide")
	}

	guide.boundArray[index].fst -= reduceValue
	guide.increase(index+1, nil)
}

func (guide *guide) expand(rank int, num int) {
	if rank <= 0 { return }
	if len(guide.boundArray) < rank-1 {
		panic(fmt.Sprintf("Incorrect values %d, %d", len(guide.boundArray), rank-1))
	}
	if rank > len(guide.boundArray) {
		guide.boundArray = append(guide.boundArray, pair{fst: int(math.Max(float64(num), float64(guide.upperBound - 2))), snd: rank - 1})
		var ptr *pair = nil
		guide.blocks = append(guide.blocks, &ptr)
	}
}

func (this *guide) remove(rank int) {
	if rank < 0 { return }
	if len(this.boundArray) <= rank { return }
	(*this.blocks[rank]) = nil
	this.blocks = this.blocks[:len(this.blocks) - 1]
	this.boundArray = this.boundArray[:len(this.boundArray) - 1]
}

func (guide *guide) ToString() string {
	str := ""
	for _, p := range guide.boundArray {
		str += strconv.Itoa(p.fst) + ","
	}
	return str + "\n"
}

func (guide *guide) blockToString() string {
	str := guide.ToString()
	str2, str3 := "", ""
	for i, b := range guide.blocks {
		if *b == nil {
			str2 += "| "
			str3 += "- "
		} else {
			if (**b).snd == i {
				str2 += "| "
				str3 += strconv.Itoa((**b).snd) + " "
			} else {
				str2 += "`-"
				str3 += "  "
			}
		}
	}
	return str + "\n" + str2 + "\n" + str3
}
