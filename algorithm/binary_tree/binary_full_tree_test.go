package binary_tree

import (
	"fmt"
	"testing"
)

func TestBCT(t *testing.T) {
	head := NewBinaryTree()
	head.Right.Right = &Node{Val: 7}
	fmt.Println(BFT(head))
}
