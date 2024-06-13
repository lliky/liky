package graph

type Graph struct {
	nodes map[int]*Node // k 代表点的编号，value 代表实际的点
	edges map[*Edge]struct{}
}

func NewGraph() *Graph {
	return &Graph{
		nodes: make(map[int]*Node),
		edges: make(map[*Edge]struct{}),
	}
}
