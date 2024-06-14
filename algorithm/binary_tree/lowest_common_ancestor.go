package binary_tree

// 最近公共祖先
// leetcode 236

func lowestCommonAncestor(root, p, q *Node) *Node {
	// 每个节点记录父节点
	parentMap := make(map[*Node]*Node)
	parentMap[root] = nil
	processLCA(root, parentMap)

	cur := p
	pMap := make(map[*Node]struct{})
	pMap[p] = struct{}{}
	for cur != nil {
		cur = parentMap[cur]
		pMap[cur] = struct{}{}
	}
	cur = q
	for cur != nil {
		if _, ok := pMap[cur]; ok {
			return cur
		}
		cur = parentMap[cur]
	}
	return nil
}

func processLCA(root *Node, parentMap map[*Node]*Node) {
	if root == nil {
		return
	}
	if root.Left != nil {
		parentMap[root.Left] = root
		processLCA(root.Left, parentMap)
	}
	if root.Right != nil {
		parentMap[root.Right] = root
		processLCA(root.Right, parentMap)
	}
}

func LCA1(root, p, q *Node) *Node {
	if p == root || q == root || root == nil {
		return root
	}
	left := LCA1(root.Left, p, q)
	right := LCA1(root.Right, p, q)
	if left != nil && right != nil {
		return root
	}
	if left != nil || right != nil {
		if left != nil {
			return left
		}
		return right
	}
	return nil
}
