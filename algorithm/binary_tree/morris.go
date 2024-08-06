package binary_tree

func Morris(root *TreeNode) {
	if root == nil {
		return
	}
	cur := root
	for cur != nil {
		mostRight := cur.Left
		if mostRight != nil { // 有左子树
			for mostRight.Right != nil && mostRight.Right != cur {
				mostRight = mostRight.Right
			}
			if mostRight.Right == nil { // 第一次来到 cur
				mostRight.Right = cur
				cur = cur.Left
				continue
			} else { // 第二次来到 cur
				mostRight.Right = nil
			}
		}
		cur = cur.Right
	}
}
