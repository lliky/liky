package algorithm

import "fmt"

func Fib(n int) int {
	if n <= 1 {
		return n
	}
	return Fib(n-1) + Fib(n-2)
}

func Fib2(n int, prev, next int) int {
	if n == 0 {
		return prev
	}
	n--
	return Fib2(n, next, next+prev)
}

func Morris(root *TreeNode) {
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
			if mostRight.Right == nil { // 第一次到达
				mostRight.Right = cur
				cur = cur.Left
				continue
			} else { //  第二次到达
				mostRight.Right = nil
			}
		}
		fmt.Println(cur.Val)
		cur = cur.Right
	}
}
