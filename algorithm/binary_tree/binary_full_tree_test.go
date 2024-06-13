package binary_tree

import (
	"fmt"
	"testing"
)

func TestBFT(t *testing.T) {
	head := NewBinaryTree()
	head.Right.Right = &Node{Val: 7}
	fmt.Println(BFT(head))
}
