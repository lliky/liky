package binary_tree

type Node struct {
	Val   int
	Left  *Node
	Right *Node
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
