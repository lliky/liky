package binary_tree

// 完全二叉树
// 除去最后一层外，所有层都被完全填满，并且最后一层中的所有节点都尽可能靠左。

// 条件：
// 1. 对于任意节点，有右孩子无左孩子，直接 false
// 2. 在 1 不违规的情况下，遇到孩子不双全后，后续所有节点都是叶子节点

func BCT(root *TreeNode) bool {
	if root == nil {
		return true
	}
	queue := make([]*TreeNode, 0)
	queue = append(queue, root)
	flag := false
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		if node.Right != nil && node.Left == nil {
			return false
		}

		if node.Left != nil {
			queue = append(queue, node.Left)
		}
		if node.Right != nil {
			queue = append(queue, node.Right)
		}

		if flag && (node.Right != nil || node.Left != nil) {
			return false
		}

		if node.Right == nil {
			flag = true
		}
	}
	return true
}
