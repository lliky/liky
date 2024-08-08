package binary_tree

import (
	"fmt"
	"testing"
)

func NewBinaryTree() *TreeNode {
	head := &TreeNode{Val: 1}
	head.Left = &TreeNode{Val: 2}
	head.Right = &TreeNode{Val: 3}
	head.Left.Left = &TreeNode{Val: 4}
	head.Left.Right = &TreeNode{Val: 5}
	head.Right.Left = &TreeNode{Val: 6}
	head.Right.Right = &TreeNode{Val: 7}
	return head
}

func TestPreOrderRecur(t *testing.T) {
	head := NewBinaryTree()
	PreOrderRecur(head)
	fmt.Println()
}

func TestPreOrderUnRecur(t *testing.T) {
	head := NewBinaryTree()
	PreOrderUnRecur(head)
	fmt.Println()
}

func TestInOrderRecur(t *testing.T) {
	head := NewBinaryTree()
	InOrderRecur(head)
	fmt.Println()
}

func TestInOrderUnRecur(t *testing.T) {
	head := NewBinaryTree()
	InOrderUnRecur(head)
}

func TestPostOrderRecur(t *testing.T) {
	head := NewBinaryTree()
	PostOrderRecur(head)
	fmt.Println()
}

func TestPostOrderUnRecur(t *testing.T) {
	head := NewBinaryTree()
	PostOrderUnRecur(head)
}
