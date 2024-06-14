package graph

import (
	"fmt"
	"testing"
)

func TestTopologicalSorting(t *testing.T) {
	matrix := [][]int{
		{1, 5, 0},
		{1, 2, 0},
		{2, 3, 0},
		{2, 4, 0},
		{5, 4, 0},
	}
	graph := CreateGraph(matrix)
	fmt.Println(len(graph.nodes))
	res := TopologicalSorting(graph)
	for _, v := range res {
		t.Logf("the node: %d\n", v.value)
	}
}
