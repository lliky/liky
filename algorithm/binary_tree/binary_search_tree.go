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

type BSTData struct {
	isBST bool
	max   int
	min   int
}

func BST2(root *Node) bool {
	return processBST2(root).isBST
}

func processBST2(root *Node) *BSTData {
	if root == nil {
		return nil
	}
	left := processBST2(root.Left)
	right := processBST2(root.Right)

	isBst := true
	if left != nil && (!left.isBST || left.max >= root.Val) {
		isBst = false
	}
	if right != nil && (!right.isBST || right.min <= root.Val) {
		isBst = false
	}
	maxV := root.Val
	minV := root.Val
	if left != nil {
		maxV = max(maxV, left.max)
		minV = min(minV, left.min)
	}
	if right != nil {
		maxV = max(maxV, right.max)
		minV = min(minV, right.min)
	}

	return &BSTData{
		isBST: isBst,
		max:   maxV,
		min:   minV,
	}
}

func BST3(root *Node) bool {
	s := make([]*Node, 0)
	var preVal = -1 << 63
	for len(s) > 0 || root != nil {
		if root != nil {
			s = append(s, root)
			root = root.Left
		} else {
			root = s[len(s)-1]
			s = s[:len(s)-1]
			if preVal >= root.Val {
				return false
			}
			preVal = root.Val
			root = root.Right
		}
	}
	return true
}
