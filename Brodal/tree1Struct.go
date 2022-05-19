package Brodal

type tree1Struct struct {
	tree1           *tree
	numOfNodesInT1W []int
	rankPointersT1W []*node
	t1GuideW        *guide
}

func newT1S(value valType) *tree1Struct {
	return &tree1Struct{
		tree1:           newTree(value, 1),
		numOfNodesInT1W: []int{},
		rankPointersT1W: []*node{},
		t1GuideW:        newGuide(6),
	}
}

func newEmptyT1S() *tree1Struct {
	return &tree1Struct{
		tree1:           nil,
		numOfNodesInT1W: []int{},
		rankPointersT1W: []*node{},
		t1GuideW:        newGuide(6),
	}
}

func (t1s *tree1Struct) getTree() (*tree) {
	return t1s.tree1
}
