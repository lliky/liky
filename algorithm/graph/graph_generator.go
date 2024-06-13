package graph

// matrix 所有的边
// N * 3 的矩阵
// [from 节点上面的值，to 节点上面的值, weight]

func CreateGraph(matrix [][]int) *Graph {
	graph := NewGraph()
	for i := 0; i < len(matrix); i++ {
		from := matrix[i][0]
		to := matrix[i][1]
		weight := matrix[i][2]
		if _, ok := graph.nodes[from]; !ok {
			graph.nodes[from] = NewNode(from)
		}
		if _, ok := graph.nodes[to]; !ok {
			graph.nodes[to] = NewNode(to)
		}
		fromNode := graph.nodes[from]
		toNode := graph.nodes[to]
		newEdge := NewEdge(weight, fromNode, toNode)
		fromNode.nexts = append(fromNode.nexts, toNode)
		fromNode.out++
		toNode.in++
		fromNode.edges = append(fromNode.edges, newEdge)
		graph.edges[newEdge] = struct{}{}
	}

	return graph
}
