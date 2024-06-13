package graph

type Node struct {
	value int
	in    int
	out   int
	nexts []*Node
	edges []*Edge
}

func NewNode(value int) *Node {
	return &Node{
		value: value,
		nexts: make([]*Node, 0), // 由该点直接发散出去邻居有哪些点
		edges: make([]*Edge, 0),
	}
}
