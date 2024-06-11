package binary_tree

import (
	"fmt"
	"testing"
)

func TestBST(t *testing.T) {

}

func TestBST2(t *testing.T) {
	root := &Node{Val: 5}
	root.Left = &Node{Val: 4}
	root.Right = &Node{Val: 6}
	root.Right.Left = &Node{Val: 3}
	root.Right.Right = &Node{Val: 7}
	fmt.Println(BST2(root))
}
