package binary_tree

import (
	"fmt"
	"testing"
)

func TestMorris(t *testing.T) {
	root := &TreeNode{Val: 1}
	root.Left = &TreeNode{Val: 2}
	root.Right = &TreeNode{Val: 3}
	root.Left.Left = &TreeNode{Val: 4}
	root.Left.Right = &TreeNode{Val: 5}
	root.Right.Left = &TreeNode{Val: 6}
	root.Right.Right = &TreeNode{Val: 7}
	Morris(root)
}

func TestPreMorris(t *testing.T) {
	root := NewBinaryTree()
	PreMorris(root)
	PreOrderRecur(root)
}

func TestInMorris(t *testing.T) {
	root := NewBinaryTree()
	InMorris(root)
	InOrderRecur(root)
}

func TestPostMorris(t *testing.T) {
	root := NewBinaryTree()
	PostMorris(root)
	fmt.Println("=============")
	PostOrderRecur(root)
}

func TestPrint(t *testing.T) {
	root := NewBinaryTree()
	Print(root)
}
