package graph

// 拓扑排序

// 卡恩算法
// 队列 L 作为结果队列
// 1. 先找到入度为 0 的节点，放到队列 L
// 2. 找到和入度为 0 相连的节点去掉，入度减一
// 3. 在找新入度为 0 的节点，放到队列 L
// 4. 重复，直到找不到入度为 0 的节点
// 如果队列 L 和节点数相等，说明排序完成
// 如果队列 L 和节点数不等，说明图中有环，无法进行拓扑排序

func TopologicalSorting(graph *Graph) []*Node {
	// key: node
	// value: 剩余的入度
	inMap := make(map[*Node]int)
	// 入度为 0 的队列
	zeroInQueue := make([]*Node, 0)
	for _, node := range graph.nodes {
		inMap[node] = node.in
		if node.in == 0 {
			zeroInQueue = append(zeroInQueue, node)
		}
	}
	// 拓扑排序保存结果队列
	result := make([]*Node, 0)
	for len(zeroInQueue) > 0 {
		cur := zeroInQueue[0]
		zeroInQueue = zeroInQueue[1:]
		result = append(result, cur)
		for _, next := range cur.nexts {
			inMap[next]--
			if inMap[next] == 0 {
				zeroInQueue = append(zeroInQueue, next)
			}

		}
	}
	return result
}
