package binary_tree

import (
	"fmt"
	"testing"
)

func TestWidthSerialize(*testing.T) {
	root := &TreeNode{Val: 1}
	root.Left = &TreeNode{Val: 2}
	root.Right = &TreeNode{Val: 3}
	root.Left.Right = &TreeNode{Val: 4}
	root.Left.Right.Right = &TreeNode{Val: 5}
	fmt.Println(WidthSerialize(root))
}
