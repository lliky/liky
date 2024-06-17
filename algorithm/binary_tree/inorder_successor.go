package binary_tree

import "fmt"

// 二叉树(一般树和搜索树)
// 1. p 无右子树，往上找左孩子
// 2. p 有右子树，右子树最左
// 3. 根节点没有，右子树

// InorderSuccessor1 BST + 递归
//  1. p.Val > root.Val 说明 p 的后继节点在 root 的右子树
//  2. p.Val = root.Val, 说明 p == root, 后继节点在 root 的右子树   // 1.2 可以合并
//  3. p.Val < root.Val, 说明 p 的后继节点在左子树或者 root 节点
//     i.如果左子树找到后继节点，直接返回
//     ii. 如果左子树没有找到后继结点，说明 p 没有右子树，那么 root 就是它的后继结点
func InorderSuccessor1(root, p *TreeNode) *TreeNode {
	if root == nil {
		return nil
	}
	if p.Val >= root.Val {
		return InorderSuccessor1(root.Right, p)
	}
	res := InorderSuccessor1(root.Left, p)
	if res != nil {
		return res
	}
	return root
}

// InorderSuccessor2 BST + 非递归
// 1. 如果 p 存在右子树，那么后继节点在右子树最左边
// 2. 如果 p 不存在左子树，那么后续节点，就是节点 p 往上到 root 路径中，第一个左儿子在路径上的点。
// 因为这个节点的左子树最大值就是 p。
func InorderSuccessor2(root, p *TreeNode) *TreeNode {
	if p.Right != nil { // 存在右子树
		p = p.Right
		for p.Left != nil {
			p = p.Left
		}
		return p
	}
	var res *TreeNode
	for root != p { // 相当于单链表找前一个
		if p.Val > root.Val {
			root = root.Right
		} else {
			res = p
			root = root.Left
		}
	}
	return res
}

// InorderSuccessor3 一般树 + 递归
// 通过中序遍历，把中序遍历的前一个节点保留下来
var res = make([]*TreeNode, 0)

func InorderSuccessor3(root, p *TreeNode) *TreeNode {
	processInorder(root)
	fmt.Println(res)
	for i := range res {
		if i+1 < len(res) && res[i] == p {
			return res[i+1]
		}
	}
	return nil
}

func processInorder(root *TreeNode) {
	if root == nil {
		return
	}
	processInorder(root.Left)
	res = append(res, root)
	processInorder(root.Right)
}

// InorderSuccessor4 一般树 + 非递归
func InorderSuccessor4(root, p *TreeNode) *TreeNode {
	s := make([]*TreeNode, 0)
	var pre *TreeNode
	for len(s) > 0 || root != nil {
		if root != nil {
			s = append(s, root)
			root = root.Left
		} else {
			root = s[len(s)-1]
			s = s[:len(s)-1]
			if p == pre {
				return root
			}
			pre = root
			root = root.Right
		}
	}
	return nil
}
