package binary_tree

import (
	"fmt"
)

// 前序遍历

// PreOrderRecur 前序递归遍历
func PreOrderRecur(head *TreeNode) {
	if head == nil {
		return
	}
	fmt.Printf("%v ", head.Val)
	PreOrderRecur(head.Left)
	PreOrderRecur(head.Right)
}

// PreOrderUnRecur 前序非递归遍历
func PreOrderUnRecur(head *TreeNode) {
	if head == nil {
		return
	}
	s := make([]*TreeNode, 0)
	s = append(s, head)
	for len(s) > 0 {
		cur := s[len(s)-1]
		s = s[:len(s)-1]
		fmt.Printf("%v ", cur.Val)
		if cur.Right != nil {
			s = append(s, cur.Right)
		}
		if cur.Left != nil {
			s = append(s, cur.Left)
		}
	}
	fmt.Println()
}
func InOrderRecur(head *TreeNode) {
	if head == nil {
		return
	}
	InOrderRecur(head.Left)
	fmt.Printf("%v ", head.Val)
	InOrderRecur(head.Right)
}

func InOrderUnRecur(head *TreeNode) {
	if head == nil {
		return
	}
	s := make([]*TreeNode, 0)
	for len(s) > 0 || head != nil {
		if head != nil {
			s = append(s, head)
			head = head.Left
		} else {
			head = s[len(s)-1]
			s = s[:len(s)-1]
			fmt.Printf("%v ", head.Val)
			head = head.Right
		}
	}
	fmt.Println()
}

func PostOrderRecur(head *TreeNode) {
	if head == nil {
		return
	}
	PostOrderRecur(head.Left)
	PostOrderRecur(head.Right)
	fmt.Printf("%v ", head.Val)
}

func PostOrderUnRecur(head *TreeNode) {
	if head == nil {
		return
	}
	s, res := make([]*TreeNode, 0), make([]*TreeNode, 0)
	s = append(s, head)
	for len(s) > 0 {
		cur := s[len(s)-1]
		s = s[:len(s)-1]
		res = append(res, cur)
		if cur.Left != nil {
			s = append(s, cur.Left)
		}
		if cur.Right != nil {
			s = append(s, cur.Right)
		}
	}
	for i := len(res) - 1; i >= 0; i-- {
		fmt.Printf("%d ", res[i].Val)
	}
	fmt.Println()
}
