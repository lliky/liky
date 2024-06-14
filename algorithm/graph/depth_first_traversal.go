package graph

import "fmt"

//深度优先遍历 进栈的时候进行处理

func DFS(node *Node) {
	stack := make([]*Node, 0)
	nodeMap := make(map[*Node]struct{})
	stack = append(stack, node)
	nodeMap[node] = struct{}{}
	fmt.Println(node.value)
	for len(stack) > 0 {
		cur := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		for _, next := range cur.nexts {
			if _, ok := nodeMap[next]; !ok {
				stack = append(stack, cur) // 需要把当前的点压栈回去
				stack = append(stack, next)
				nodeMap[next] = struct{}{}
				fmt.Println(next.value)
				break
			}
		}
	}
}
