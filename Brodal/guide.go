package Brodal

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

	if actualValue > guide.upperBound-2 {
		guide.increase(index, nil)

		if guide.boundArray[index].fst == guide.upperBound-1 && (*guide.blocks[index]) != nil {
			guide.fixUp(*guide.blocks[index], reduceValue, &ops)
		} else if guide.boundArray[index].fst == guide.upperBound {
			if (*guide.blocks[index]) != nil {
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
	guide.increase(pair.snd+1, nil)

	(*guide.blocks[pair.snd]) = nil

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

func (guide *guide) increase(index int, ops *[]action) {
	guide.boundArray[index].fst++
	if ops != nil {
		*ops = append(*ops, action{index, Increase, 1})
	}
}

func (guide *guide) reduce(index int, value int, ops *[]action) {
	if ops != nil {
		*ops = append(*ops, action{index, Reduce, value})
	}

	if guide.upperBound == UPPER_BOUND || guide.upperBound == LOWER_BOUND {
		guide.boundArray[index].fst -= value
		guide.boundArray[index+1].fst++
	}
}

func (guide *guide) expand() {
	guide.blocks = append(guide.blocks)
	guide.boundArray = append(guide.boundArray)
}

func (guide *guide) update(value int, rank uint) {
	if value >= guide.upperBound-2 {
		guide.boundArray[rank].fst = value
	}
}
