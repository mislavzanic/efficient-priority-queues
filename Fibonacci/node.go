package Fibonacci

type Node struct {
    val   int64

	degree int
	size   int
	mark  bool

	next  *Node
	prev  *Node
	child *Node
	parent *Node
}

func NewNode(val int64) *Node {
	node := &Node{
		val:    val,

		mark:   false,
		size:   1,
		degree: 0,

		next:   nil,
		prev:   nil,
		child:  nil,
		parent: nil,
	}

	node.next = node
	node.prev = node
	return node
}

func (this *Node) GetVal() int64 {
	return this.val;
}

func (this *Node) GetNext() *Node {
	return this.next
}

func (this *Node) GetPrev() *Node {
	return this.prev
}

func (this *Node) InsertNext(node *Node) {
	if this.next == nil {
		this.next = node
		node.prev = this
	} else {
		prev_next := this.next
		node.prev = this;
		node.next = prev_next

		prev_next.prev = node
		this.next = node;
	}
}

func (this *Node) InsertPrev(node *Node) {
	if this.prev == nil {
		this.prev = node
		node.next = this
	} else {
		this.prev.InsertNext(node)
	}
}

