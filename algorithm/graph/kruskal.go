package graph

import "fmt"

func KruskalKMT(graph *Graph) []*Edge {
	mySets := NewMySets(graph.nodes)
	pq := NewPriorityQueue()
	for edge, _ := range graph.edges {
		pq.Add(edge)
	}
	result := make([]*Edge, 0)
	i := 0
	for pq.Len() > 0 {
		edge := pq.Pop()
		i++
		fmt.Printf("i========%d, %v\n", i, mySets)
		if !mySets.IsSameSet(edge.from, edge.to) {
			fmt.Printf("aaaaaaaaa: from %v, to: %v\n", edge.from.value, edge.to.value)
			result = append(result, edge)
			mySets.Union(edge.from, edge.to)
		}
	}
	fmt.Println(mySets)
	return result
}

type PriorityQueue struct {
	edges []*Edge
}

func NewPriorityQueue() *PriorityQueue {
	return &PriorityQueue{
		edges: make([]*Edge, 0),
	}
}

func (pq *PriorityQueue) Add(edge *Edge) {
	pq.edges = append(pq.edges, edge)
	pq.heapInsert(pq.Len() - 1)
}

func (pq *PriorityQueue) Pop() *Edge {
	edge := pq.edges[0]
	pq.edges = pq.edges[1:]
	for i := pq.Len() / 2; i >= 0; i-- {
		pq.heapify(i, pq.Len())
	}
	return edge
}

func (pq *PriorityQueue) Len() int {
	return len(pq.edges)
}

func (pq *PriorityQueue) heapInsert(index int) {
	for pq.edges[index].weight < pq.edges[(index-1)/2].weight {
		pq.edges[index], pq.edges[(index-1)/2] = pq.edges[(index-1)/2], pq.edges[index]
		index = (index - 1) / 2
	}
}

func (pq *PriorityQueue) heapify(index, heapSize int) {
	left := index*2 + 1
	for left < heapSize {
		smallest := left
		if left+1 < heapSize && pq.edges[left].weight > pq.edges[left+1].weight {
			smallest = left + 1
		}
		if pq.edges[smallest].weight > pq.edges[index].weight {
			break
		}
		pq.edges[smallest], pq.edges[index] = pq.edges[index], pq.edges[smallest]

		index = smallest
		left = index*2 + 1
	}
}

type MySets map[*Node]map[int]struct{}

func NewMySets(nodes map[int]*Node) MySets {
	setMap := make(map[*Node]map[int]struct{})
	for _, v := range nodes {
		set := map[int]struct{}{
			v.value: {},
		}
		setMap[v] = set
	}
	return setMap
}

func (mySets MySets) IsSameSet(from, to *Node) bool {
	fromSet := mySets[from]
	_, ok := fromSet[to.value]
	return ok
}

// 把 to 加入到 from 中, 并且把 to 指向 from 集合

func (mySets MySets) Union(from, to *Node) {
	var src, dst map[int]struct{}
	if from.value > to.value {
		src = mySets[from]
		dst = mySets[to]
	} else {
		src = mySets[to]
		dst = mySets[from]
	}
	for k, v := range src {
		dst[k] = v
		mySets[to] = dst
		mySets[from] = dst
	}
}
