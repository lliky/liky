package algorithm

import (
	"fmt"
	"testing"
)

func TestFib(t *testing.T) {
	fmt.Println(Fib(45))
}

func TestFib2(t *testing.T) {
	fmt.Println(Fib2(45, 0, 1))
}

func TestA(t *testing.T) {
	fmt.Println(-1 / 2)
}

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
