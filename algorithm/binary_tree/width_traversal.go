package binary_tree

import "fmt"

func WidthTraversal(head *Node) {
	if head == nil {
		return
	}
	q := make([]*Node, 0)
	q = append(q, head)
	for len(q) > 0 {
		cur := q[0]
		q = q[1:]
		fmt.Printf("%v ", cur.Val)
		if cur.Left != nil {
			q = append(q, cur.Left)
		}
		if cur.Right != nil {
			q = append(q, cur.Right)
		}
	}
	fmt.Println()
}
