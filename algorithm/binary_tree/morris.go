package binary_tree

import "fmt"

func Morris(root *TreeNode) {
	if root == nil {
		return
	}
	cur := root
	for cur != nil {
		mostRight := cur.Left
		fmt.Printf("%d ", cur.Val)
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
				cur = cur.Right
			}
		} else {
			cur = cur.Right
		}
	}
	fmt.Println()
}

// PreMorris 前序遍历
// 遇到就打印
func PreMorris(root *TreeNode) {
	if root == nil {
		return
	}
	cur := root
	for cur != nil {
		mostRight := cur.Left
		if mostRight != nil { // 左子树
			for mostRight.Right != nil && mostRight.Right != cur {
				mostRight = mostRight.Right
			}
			if mostRight.Right == nil { // 第一次到 cur
				fmt.Printf("%d ", cur.Val)
				mostRight.Right = cur
				cur = cur.Left
			} else { // 第二次到达 cur
				mostRight.Right = nil
				cur = cur.Right
			}
		} else { // 左子树为空，进入右子树
			fmt.Printf("%d ", cur.Val)
			cur = cur.Right
		}
	}
	fmt.Println()
}

// InMorris 中序遍历
// 如果有第二次到达，第二次打印否则就第一次打印
func InMorris(root *TreeNode) {
	if root == nil {
		return
	}
	cur := root
	for cur != nil {
		mostRight := cur.Left
		if mostRight != nil {
			for mostRight.Right != nil && mostRight.Right != cur {
				mostRight = mostRight.Right
			}
			if mostRight.Right == nil {
				mostRight.Right = cur
				cur = cur.Left
			} else {
				fmt.Printf("%d ", cur.Val)
				mostRight.Right = nil
				cur = cur.Right
			}
		} else {
			fmt.Printf("%d ", cur.Val)
			cur = cur.Right
		}
	}
	fmt.Println()
}

// PostMorris 后序遍历
// a. 第二次节点，把节点的左子树右边节点从底往上打印
// b. 把根节点的右边节点从底往上打印
func PostMorris(root *TreeNode) {
	if root == nil {
		return
	}
	cur := root
	for cur != nil {
		mostRight := cur.Left
		if mostRight != nil {
			for mostRight.Right != nil && mostRight.Right != cur {
				mostRight = mostRight.Right
			}
			if mostRight.Right == nil {
				mostRight.Right = cur
				cur = cur.Left
			} else { // 第二次到达节点
				mostRight.Right = nil
				Print(cur.Left)
				cur = cur.Right
			}
		} else {
			cur = cur.Right
		}
	}
	Print(root)
}

func Print(root *TreeNode) {
	root = reverseNode1(root)
	printNode(root)
	reverseNode1(root)
}

func reverseNode1(head *TreeNode) *TreeNode {
	if head == nil || head.Right == nil {
		return head
	}
	newHead := reverseNode1(head.Right)
	head.Right.Right = head
	head.Right = nil
	return newHead
}

func reverseNode(head *TreeNode) *TreeNode {
	if head == nil || head.Right == nil {
		return head
	}
	d := new(TreeNode)
	cur := head
	for cur != nil {
		tmp := cur.Right
		cur.Right = d.Right
		d.Right = cur
		cur = tmp
	}
	return d.Right
}
func printNode(head *TreeNode) {
	if head == nil {
		return
	}
	cur := head
	for cur != nil {
		fmt.Printf("%d ", cur.Val)
		cur = cur.Right
	}
}
