package binary_tree

import "math"

// leetcode 98
// 二叉搜索数性质如下：
// 1. 节点的左子树只包含小于当前节点的数
// 2. 节点的右子树只包含大于当前节点的数
// 3. 所有左子树和右子树自身必须也是二叉搜索树
// 中序遍历，就是递增的

func BST1(root *Node) bool {
	var preVal = math.MinInt
	return processBST1(root, &preVal)
}

func processBST1(root *Node, preVal *int) bool {
	if root == nil {
		return true
	}
	res := processBST1(root.Left, preVal)
	if !res {
		return res
	}
	if *preVal >= root.Val {
		return false
	} else {
		*preVal = root.Val
	}
	return processBST1(root.Right, preVal)
}
