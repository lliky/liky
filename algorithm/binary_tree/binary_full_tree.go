package binary_tree

// 满二叉树，每层节点都是满的

func BFT(root *Node) bool {
	data := processBFT(root)
	return data.n == 1<<data.level-1
}

type BFTData struct {
	level int
	n     int
}

func processBFT(root *Node) BFTData {
	if root == nil {
		return BFTData{level: 0, n: 0}
	}
	left := processBFT(root.Left)
	right := processBFT(root.Right)
	return BFTData{
		level: max(left.level, right.level) + 1,
		n:     left.n + right.n + 1,
	}
}
