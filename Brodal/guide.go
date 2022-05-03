package Brodal

type pair struct {
	fst uint
	snd uint
}

type guide struct {
	upperBound uint
	boundArray []pair
	blocks     []**pair
}

type operation byte

const (
	Increase operation = 0
	Reduce   operation = 1
)

type action struct {
	index uint
	op    operation
}

func newGuide(upperBound uint) *guide {
	return &guide{
		upperBound: upperBound,
		boundArray: []pair{},
		blocks:     []**pair{},
	}
}

func (guide *guide) increase(index uint) {
	ops := []action{}

	guide.boundArray[index].fst++
	ops = append(ops, action{index, Increase})

	if guide.boundArray[index].fst == guide.upperBound - 1 && (*guide.blocks[index]) != nil {
		guide.fixUp(*guide.blocks[index], &ops)
	} else if guide.boundArray[index].fst == guide.upperBound {
		if (*guide.blocks[index]) != nil {
			guide.fixUp(&guide.boundArray[index], &ops)
		} else {
			guide.fixUp(*guide.blocks[index], &ops)
			guide.fixUp(&guide.boundArray[index], &ops)
		}
	}
}

func (guide *guide) reduce(index uint, value int) {
}

func (guide *guide) fixUp(pair *pair, ops *[]action) {
	guide.boundArray[pair.snd + 1].fst++
	*ops = append(*ops, action{pair.snd + 1, Increase})

	(*guide.blocks[pair.snd]) = nil

	if guide.boundArray[pair.snd + 1].fst == guide.upperBound - 1 {
		if (*guide.blocks[pair.snd + 1]) != nil {
			guide.blocks[pair.snd] = guide.blocks[pair.snd + 1]
		}
	}

	if guide.boundArray[pair.snd + 1].fst == guide.upperBound {
		ptr := &guide.boundArray[pair.snd + 1]
		guide.blocks[pair.snd + 1] = &ptr
		guide.blocks[pair.snd] = &ptr
	}
}
