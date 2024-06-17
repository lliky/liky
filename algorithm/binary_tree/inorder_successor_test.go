package binary_tree

import (
	"testing"
)

func TestInorderSuccessor1(t *testing.T) {
	root := NewBST()
	p := root.Left.Right
	InorderSuccessor1(root, p)
}

func TestInorderSuccessor3(t *testing.T) {
	root := NewBST()
	p := root.Left.Right
	a := InorderSuccessor3(root, p)
	t.Logf("%v", a)
}

func NewBST() *TreeNode {
	root := &TreeNode{Val: 10}
	root.Left = &TreeNode{Val: 5}
	root.Right = &TreeNode{Val: 17}
	root.Left.Left = &TreeNode{Val: 2}
	root.Left.Right = &TreeNode{Val: 7}
	root.Right.Left = &TreeNode{Val: 14}
	root.Right.Right = &TreeNode{Val: 25}
	root.Left.Left.Left = &TreeNode{Val: 1}
	root.Left.Left.Right = &TreeNode{Val: 3}
	root.Left.Right.Left = &TreeNode{Val: 6}
	root.Left.Right.Right = &TreeNode{Val: 9}
	return root
}
