package binary_tree

import "math"

// 平衡二叉树
//	性质：该树所有节点的左右子树高度差不超过 1

func BBT(root *Node) bool {
	return processBBT(root).isBBT
}

type BBTData struct {
	isBBT  bool
	height int
}

func processBBT(root *Node) BBTData {
	if root == nil {
		return BBTData{isBBT: true, height: 0}
	}
	left := processBBT(root.Left)
	right := processBBT(root.Right)
	isBBT := left.isBBT && right.isBBT
	if math.Abs(float64(left.height-right.height)) > 1 {
		isBBT = false
	}
	return BBTData{isBBT: isBBT, height: max(left.height, right.height) + 1}
}
