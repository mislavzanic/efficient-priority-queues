package Brodal

import (
	"fmt"
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
)

type action struct {
	index int
	op    operation
	value int
}

func newGuide(upperBound int) *guide {
	return &guide{
		upperBound: upperBound,
		boundArray: []pair{},
		blocks:     []**pair{},
	}
}

func (guide *guide) forceIncrease(index int, actualValue int, reduceValue int) []action {
	ops := []action{}

	if actualValue > guide.boundArray[index].fst {

		// slucaj kad je x_i == upperBound -> prvo fixup x_i onda increase
		if guide.boundArray[index].fst == guide.upperBound {
			guide.fixUp(&guide.boundArray[index], reduceValue, &ops)
		}

		guide.increase(index, &ops)

		// if actualValue - reduceValue + 1 <= guide.upperBound - 2 { return ops }

		if guide.boundArray[index].fst == guide.upperBound-1 && (*guide.blocks[index]) != nil {
			guide.fixUp(*guide.blocks[index], reduceValue, &ops)
		} else if guide.boundArray[index].fst == guide.upperBound {
			if (*guide.blocks[index]) == nil {
				guide.fixUp(&guide.boundArray[index], reduceValue, &ops)
			} else {
				guide.fixUp(*guide.blocks[index], reduceValue, &ops)
				guide.fixUp(&guide.boundArray[index], reduceValue, &ops)
			}
		}
	}

	return ops
}

func (guide *guide) fixUp(pair *pair, reduceValue int, ops *[]action) {
	guide.reduce(pair.snd, reduceValue, ops)
	(*guide.blocks[pair.snd]) = nil

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
			*ops = append(*ops, action{index, Increase, 1})
		}
	}
}

func (guide *guide) reduce(index int, reduceValue int, ops *[]action) {
	if ops != nil {
		*ops = append(*ops, action{index, Reduce, reduceValue})
	}

	if guide.boundArray[index].fst != guide.upperBound {
		panic("Ipak je moguce... ja sam u guide")
	}

	guide.boundArray[index].fst -= reduceValue
	// if guide.boundArray[index].fst < guide.upperBound - 2 {
	// 	guide.boundArray[index].fst = guide.upperBound - 2
	// }
	guide.increase(index+1, nil)
}

func (guide *guide) expand(rank int) {
	if len(guide.boundArray) < rank-1 {
		panic(fmt.Sprintf("Incorrect values %d, %d", len(guide.boundArray), rank-1))
	}
	if rank > len(guide.boundArray) {
		guide.boundArray = append(guide.boundArray, pair{fst: UPPER_BOUND - 2, snd: rank - 1})
		ptr := &guide.boundArray[rank-1]
		guide.blocks = append(guide.blocks, &ptr)
	}
}

func (guide *guide) ToString() string {
	str := ""
	for _, p := range guide.boundArray {
		str += strconv.Itoa(p.fst) + ","
	}
	return str
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
