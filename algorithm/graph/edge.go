package graph

type Edge struct {
	weight int
	from   *Node
	to     *Node
}

func NewEdge(weight int, from, to *Node) *Edge {
	return &Edge{
		weight: weight,
		from:   from,
		to:     to,
	}
}
