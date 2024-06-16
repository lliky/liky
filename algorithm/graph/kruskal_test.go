package graph

import (
	"fmt"
	"testing"
)

func TestKruskalKMT(t *testing.T) {
	matrix := [][]int{
		{1, 2, 10},
		{2, 1, 10},
		{1, 3, 1},
		{3, 1, 1},
		{1, 4, 7},
		{4, 1, 7},
		{2, 3, 2},
		{3, 2, 2},
		{3, 4, 3},
		{4, 3, 3},
		{2, 5, 5},
		{5, 2, 5},
		{3, 5, 4},
		{5, 3, 4},
		{4, 5, 6},
		{5, 4, 6},
	}
	graph := CreateGraph(matrix)
	fmt.Println(len(graph.edges))
	fmt.Println(len(graph.nodes))
	res := KruskalKMT(graph)
	fmt.Println(res)
	for _, v := range res {
		fmt.Println(v)
	}
}

func TestNewMySets(t *testing.T) {
	node1 := NewNode(1)
	node2 := NewNode(2)
	node3 := NewNode(3)
	node4 := NewNode(4)
	nodes := map[int]*Node{
		node1.value: node1,
		node2.value: node2,
		node3.value: node3,
		node4.value: node4,
	}
	fmt.Printf("nodes: %v\n", nodes)
	sets := NewMySets(nodes)
	fmt.Printf("%v\n", sets)
	sets.Union(node4, node1)
	fmt.Printf("%v\n", sets)
	fmt.Println(sets.IsSameSet(node1, node4))
	fmt.Println(sets.IsSameSet(node2, node1))
}

func TestNewEdge(t *testing.T) {
	edge1 := NewEdge(1, nil, nil)
	edge2 := NewEdge(5, nil, nil)
	edge3 := NewEdge(4, nil, nil)
	edge4 := NewEdge(2, nil, nil)
	edge5 := NewEdge(-1, nil, nil)
	edge6 := NewEdge(3, nil, nil)
	pq := NewPriorityQueue()
	pq.Add(edge1)
	pq.Add(edge2)
	pq.Add(edge4)
	pq.Add(edge3)
	fmt.Println(pq.Pop())
	pq.Add(edge6)
	pq.Add(edge5)
	for pq.Len() > 0 {
		fmt.Println(pq.Pop())
	}
}
