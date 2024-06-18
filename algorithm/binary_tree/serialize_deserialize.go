package binary_tree

import (
	"strconv"
	"strings"
)

// 序列化与反序列化
// leetcode 297

// end flag : "#"
// split: "_"

const (
	flag  = "#"
	split = "_"
)

// 前序遍历

func PreSerialize(root *TreeNode) string {
	sb := strings.Builder{}
	var dfs func(*TreeNode)
	dfs = func(root *TreeNode) {
		if root == nil {
			sb.WriteString(flag + split)
			return
		}
		sb.WriteString(strconv.Itoa(root.Val) + split)
		dfs(root.Left)
		dfs(root.Right)
	}
	dfs(root)
	return sb.String()
}

func PreDeserialize(data string) *TreeNode {
	sS := strings.Split(data, split)
	// 后面的字符串是 ""
	sS = sS[:len(sS)-1]
	var dfs func() *TreeNode
	dfs = func() *TreeNode {
		x := sS[0]
		sS = sS[1:]
		if x == flag {
			return nil
		}
		val, _ := strconv.Atoi(x)
		left := dfs()
		right := dfs()
		return &TreeNode{
			Val:   val,
			Left:  left,
			Right: right,
		}
	}
	return dfs()
}

// 后序遍历序列化
// 左子树，右子树，根节点

func PostSerialize(root *TreeNode) string {
	sb := strings.Builder{}
	var dfs func(root *TreeNode)
	dfs = func(root *TreeNode) {
		if root == nil {
			sb.WriteString(flag + split)
			return
		}
		dfs(root.Left)
		dfs(root.Right)
		sb.WriteString(strconv.Itoa(root.Val) + split)
	}
	dfs(root)
	return sb.String()
}

// 后序遍历最后一个元素是根节点，我们从最一位开始读取
// 左右头，那么新建的时候，头右左

func PosDeserialize(str string) *TreeNode {
	datas := strings.Split(str, split)
	datas = datas[:len(datas)-1] // 这样分割模式，最后一位是空字符串，特殊处理下
	var dfs func() *TreeNode
	dfs = func() *TreeNode {
		x := datas[len(datas)-1] // 提前取出来
		datas = datas[:len(datas)-1]
		if x == flag {
			return nil
		}
		right := dfs()
		left := dfs()
		val, _ := strconv.Atoi(x)
		return &TreeNode{val, left, right}
	}
	return dfs()
}

// 中序遍历一般树，不是唯一答案，需要特定树，比如 BST

// 层序遍历

func WidthSerialize(root *TreeNode) string {
	sb := strings.Builder{}
	queue := make([]*TreeNode, 0)
	queue = append(queue, root)
	for len(queue) > 0 {
		root = queue[0]
		queue = queue[1:]
		if root == nil {
			sb.WriteString(flag + split)
		} else {
			sb.WriteString(strconv.Itoa(root.Val) + split)
			queue = append(queue, root.Left)
			queue = append(queue, root.Right)
		}
	}
	return sb.String()
}

func WidthDeserialize(str string) *TreeNode {
	datas := strings.Split(str, split)
	datas = datas[:len(datas)-1] // 去掉末尾的空字符

	rootVal := datas[0]
	datas = datas[1:]
	if rootVal == flag { // 空树
		return nil
	}
	val, _ := strconv.Atoi(rootVal)
	root := &TreeNode{Val: val}
	queue := make([]*TreeNode, 0)
	queue = append(queue, root)
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		// left
		leftVal := datas[0]
		datas = datas[1:]
		if leftVal != flag {
			val, _ = strconv.Atoi(leftVal)
			left := &TreeNode{Val: val}
			node.Left = left
			queue = append(queue, left)
		}
		// right
		rightVal := datas[0]
		datas = datas[1:]
		if rightVal != flag {
			val, _ = strconv.Atoi(rightVal)
			right := &TreeNode{Val: val}
			node.Right = right
			queue = append(queue, right)
		}
	}
	return root
}
