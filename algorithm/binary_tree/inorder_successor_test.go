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

func NewBST() *Node {
	root := &Node{Val: 10}
	root.Left = &Node{Val: 5}
	root.Right = &Node{Val: 17}
	root.Left.Left = &Node{Val: 2}
	root.Left.Right = &Node{Val: 7}
	root.Right.Left = &Node{Val: 14}
	root.Right.Right = &Node{Val: 25}
	root.Left.Left.Left = &Node{Val: 1}
	root.Left.Left.Right = &Node{Val: 3}
	root.Left.Right.Left = &Node{Val: 6}
	root.Left.Right.Right = &Node{Val: 9}
	return root
}
