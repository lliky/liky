package binary_tree

import (
	"fmt"
	"testing"
)

func NewBinaryTree() *Node {
	head := &Node{Val: 1}
	head.Left = &Node{Val: 2}
	head.Right = &Node{Val: 3}
	head.Left.Left = &Node{Val: 4}
	head.Left.Right = &Node{Val: 5}
	head.Right.Left = &Node{Val: 6}
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
