package graph

import "testing"

func TestDFS(t *testing.T) {
	matraix := [][]int{
		{1, 2, 1},
		{2, 1, 1},
		{1, 3, 4},
		{3, 1, 4},
		{2, 4, 1},
		{4, 2, 1},
	}
	graph := CreateGraph(matraix)
	node := graph.nodes[1]
	DFS(node)
}
