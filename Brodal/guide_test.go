package Brodal

import (
	"fmt"
	"testing"
)

func initGuide() guide {
	newGuide := newGuide(2)
	newGuide.boundArray = []pair{
		pair{fst:0,snd:0}, pair{fst:1, snd:1}, pair{fst:1, snd:2}, pair{fst:1, snd:3}, pair{fst:2,snd:4},
		pair{fst:0, snd:5}, pair{fst:1, snd:6}, pair{fst:1, snd:7}, pair{fst:2, snd:8},
	}
	ptr1 := &newGuide.boundArray[8]
	ptr2 := &newGuide.boundArray[4]
	newGuide.blocks = []**pair{
		&ptr2, &ptr2, &ptr2, &ptr2, &ptr2,
		&ptr1, &ptr1, &ptr1, &ptr1,
	}

	return *newGuide
}

func TestGuide(t *testing.T) {

	guide := initGuide()
	guide.boundArray = []pair{
		pair{fst:0,snd:0}, pair{fst:1, snd:1}, pair{fst:1, snd:2}, pair{fst:1, snd:3}, pair{fst:2,snd:4},
		pair{fst:0, snd:5}, pair{fst:1, snd:6}, pair{fst:1, snd:7}, pair{fst:2, snd:8},
	}
	ptr1 := &guide.boundArray[8]
	ptr2 := &guide.boundArray[4]
	guide.blocks = []**pair{
		&ptr2, &ptr2, &ptr2, &ptr2, &ptr2,
		&ptr1, &ptr1, &ptr1, &ptr1,
	}

	if guide.ToString() != "0,1,1,1,2,0,1,1,2," {
		t.Error(fmt.Sprintf("Guide repr: %s", guide.ToString()))
	}

	guide.forceIncrease(3, 2, 2)
	if guide.ToString() != "0,1,1,0,1,1,1,1,2," {
		t.Error(fmt.Sprintf("Guide repr: %s", guide.ToString()))
	}

	guide = initGuide()

	guide.forceIncrease(2, 2, 2)
	if guide.ToString() != "0,1,0,2,0,1,1,1,2," {
		t.Error(fmt.Sprintf("Guide repr: %s", guide.ToString()))
	}

	guide = initGuide()
	guide.forceIncrease(4, 3, 2)
	if guide.ToString() != "0,1,1,1,1,1,1,1,0," {
		t.Error(fmt.Sprintf("Guide repr: %s", guide.ToString()))
	}

	guide.forceIncrease(5, 3, 2)
	if guide.ToString() != "0,1,1,1,1,0,2,1,0," {
		t.Error(fmt.Sprintf("Guide repr: %s", guide.ToString()))
	}

	guide.forceIncrease(4, 3, 2)
	if guide.ToString() != "0,1,1,1,0,1,2,1,0," {
		t.Error(fmt.Sprintf("Guide repr: %s", guide.ToString()))
	}
	fmt.Println(guide.blockToString())
}
