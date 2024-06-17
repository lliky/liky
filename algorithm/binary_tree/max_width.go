package binary_tree

func MaxWidth(head *TreeNode) int {
	if head == nil {
		return 0
	}
	res := -1
	q := make([]*TreeNode, 0)
	q = append(q, head)
	for len(q) > 0 {
		index := len(q)
		res = max(res, index)
		for index > 0 {
			cur := q[0]
			q = q[1:]
			if cur.Left != nil {
				q = append(q, cur.Left)
			}
			if cur.Right != nil {
				q = append(q, cur.Right)
			}
			index--
		}
	}
	return res
}
