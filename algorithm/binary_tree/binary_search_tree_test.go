package binary_tree

import (
	"fmt"
	"testing"
)

func TestBST(t *testing.T) {

}

func TestBST2(t *testing.T) {
	root := &TreeNode{Val: 5}
	root.Left = &TreeNode{Val: 4}
	root.Right = &TreeNode{Val: 6}
	root.Right.Left = &TreeNode{Val: 3}
	root.Right.Right = &TreeNode{Val: 7}
	fmt.Println(BST2(root))
}
