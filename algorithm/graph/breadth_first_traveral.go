package graph

import "fmt"

// 宽度优先遍历

func BFS(node *Node) {
	var queue = make([]*Node, 0)
	var nodeMap = make(map[*Node]struct{}) // 可以把去重的 hash 替换成数组
	queue = append(queue, node)
	nodeMap[node] = struct{}{}
	for len(queue) > 0 {
		node = queue[0]
		queue = queue[1:]
		// 处理的地方
		fmt.Println(node.value)
		// 处理的地方
		for _, next := range node.nexts {
			if _, ok := nodeMap[next]; !ok {
				queue = append(queue, next)
				nodeMap[next] = struct{}{}
			}
		}
	}
}
