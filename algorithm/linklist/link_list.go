package linklist

import (
	"fmt"
)

// 1. 反转单向链表

// ReverseSingleLinkList1 迭代
func ReverseSingleLinkList1(head *Node) *Node {
	if head == nil || head.Next == nil {
		return head
	}
	pre := &Node{}
	for head != nil {
		tmp := head.Next
		head.Next = pre.Next
		pre.Next = head
		head = tmp
	}
	return pre.Next
}

// ReverseSingleLinkList2 用栈
func ReverseSingleLinkList2(head *Node) *Node {
	if head == nil || head.Next == nil {
		return head
	}
	help := make([]*Node, 0)
	for head != nil {
		help = append(help, head)
		head = head.Next
	}
	for i := len(help) - 1; i > 0; i-- {
		help[i].Next = help[i-1]
	}
	help[0].Next = nil
	return help[len(help)-1]
}

// ReverseSingleLinkList3 递归
func ReverseSingleLinkList3(head *Node) *Node {
	if head == nil || head.Next == nil {
		return head
	}
	cur := ReverseSingleLinkList3(head.Next)
	head.Next.Next = head // 指回去(head 的下一个节点指向 head)
	head.Next = nil       // 	切断 head 与下一个节点的联系, 成为为尾节点
	return cur
}

// 2. 给定两个有序链表的头指针 head1 和 head2，打印两个链表的公共部分。

func PrintLinkListCommonPart(head1, head2 *Node) {
	cur1, cur2 := head1, head2
	for cur1 != nil && cur2 != nil {
		if cur1.Val == cur2.Val {
			fmt.Printf(" %d ", cur1.Val)
			cur1 = cur1.Next
			cur2 = cur2.Next
		} else if cur1.Val > cur2.Val {
			cur2 = cur2.Next
		} else {
			cur1 = cur1.Next
		}
	}
	fmt.Println()
}

// 3. 将两个升序链表合并为一个新的 升序 链表并返回。新链表是通过拼接给定的两个链表的所有节点组成的。
// leetcode 21

func MergeTwoLists(head1, head2 *Node) *Node {
	cur1, cur2 := head1, head2
	head := new(Node)
	cur := head
	for cur1 != nil && cur2 != nil {
		if cur1.Val <= cur2.Val {
			cur.Next = cur1
			cur1 = cur1.Next
		} else {
			cur.Next = cur2
			cur2 = cur2.Next
		}
		cur = cur.Next
	}
	if cur1 != nil {
		cur.Next = cur1
	}
	if cur.Next != nil {
		cur.Next = cur2
	}
	return head.Next
}

// 4. 回文

// IsPalindrome1 用栈
func IsPalindrome1(head *Node) bool {
	if head == nil || head.Next == nil {
		return true
	}
	slow, fast := head, head
	help := make([]*Node, 0)
	for fast.Next != nil && fast.Next.Next != nil {
		help = append(help, slow)
		slow = slow.Next
		fast = fast.Next.Next
	}
	if fast.Next != nil {
		help = append(help, slow)
	}
	cur, i := slow.Next, len(help)-1
	for cur != nil {
		if cur.Val != help[i].Val {
			return false
		}
		cur = cur.Next
		i--
	}
	return true
}

func IsPalindrome2(head *Node) bool {
	if head == nil || head.Next == nil {
		return true
	}
	slow, fast := head, head
	for fast.Next != nil && fast.Next.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
	}
	// reverse next
	pre, cur := slow, slow.Next
	pre.Next = nil
	for cur != nil {
		tmp := cur.Next
		cur.Next = pre.Next
		pre.Next = cur
		cur = tmp
	}
	head.Print()
	prev, next := head, slow.Next
	for next != nil {
		if prev.Val != next.Val {
			return false
		}
		next = next.Next
		prev = prev.Next
	}
	// recover list
	pre, cur = slow, slow.Next
	pre.Next = nil
	for cur != nil {
		tmp := cur.Next
		cur.Next = pre.Next
		pre.Next = cur
		cur = tmp
	}
	head.Print()
	return true
}

// 5.单链表按某值划分

// PartitionByVal1 放到数组中，荷兰国旗
func PartitionByVal1(head *Node, x int) *Node {
	if head == nil || head.Next == nil {
		return head
	}
	help := make([]*Node, 0)
	for head != nil {
		help = append(help, head)
		head = head.Next
	}
	partition(help, x)
	for i := 0; i < len(help)-1; i++ {
		help[i].Next = help[i+1]
	}
	help[len(help)-1].Next = nil
	return help[0]
}

func partition(nums []*Node, x int) {
	left, right := -1, len(nums)
	i := 0
	for i < right {
		if nums[i].Val < x {
			left++
			swap(nums, i, left)
			i++
		} else if nums[i].Val > x {
			right--
			swap(nums, i, right)
		} else {
			i++
		}
	}
}

func PartitionByVal2(head *Node, x int) *Node {
	var sH, sT, mH, mT, lH, lT *Node
	cur := head
	for cur != nil {
		if cur.Val < x {
			if sH == nil {
				sH = cur
				sT = cur
			} else {
				sT.Next = cur
				sT = sT.Next
			}
		} else if cur.Val > x {
			if lH == nil {
				lH = cur
				lT = cur
			} else {
				lT.Next = cur
				lT = lT.Next
			}
		} else {
			if mH == nil {
				mH = cur
				mT = cur
			} else {
				mT.Next = cur
				mH = lT.Next
			}
		}
		cur = cur.Next
	}
	//  merge
	if sH != nil && mH != nil && lH != nil {
		sT.Next = mH
		mT.Next = lH
		lT.Next = nil
		return sH
	}
	if sH == nil {
		if mH == nil {
			return lH
		} else {
			mT.Next = lH
			if lT != nil {
				lT.Next = nil
			}
			return mT
		}
	}
	if mH == nil {
		if sH == nil {
			return lH
		} else {
			sT.Next = lH
			if lT != nil {
				lT.Next = nil
			}
			return sH
		}
	}
	if lH == nil {
		if sH == nil {
			return mH
		} else {
			sT.Next = mH
			if mT != nil {
				mT.Next = nil
			}
			return sH
		}
	}
	return nil
}

//5.复制含有随机指针节点的链表

// CopyRandomList1 用 hash 辅助空间
func CopyRandomList1(head *RandomNode) *RandomNode {
	if head == nil {
		return head
	}
	cur := head
	help := make(map[*RandomNode]*RandomNode)
	for cur != nil {
		node := &RandomNode{Val: cur.Val}
		help[cur] = node
		cur = cur.Next
	}
	cur = head
	d := new(RandomNode)
	pre := d
	for cur != nil {
		pre.Next = help[cur]
		if cur.Random != nil {
			pre.Next.Random = help[cur.Random]
		}
		cur = cur.Next
		pre = pre.Next
	}
	return d.Next
}

func CopyRandomList2(head *RandomNode) *RandomNode {
	if head == nil {
		return head
	}
	cur := head
	// copy new node
	for cur != nil {
		node := &RandomNode{Val: cur.Val}
		node.Next = cur.Next
		cur.Next = node
		cur = node.Next
	}
	head.Print()
	// copy random node
	cur = head
	for cur != nil {
		if cur.Random != nil {
			cur.Next.Random = cur.Random.Next
		}
		cur = cur.Next.Next
	}
	// delete old node
	cur = head
	d := new(RandomNode)
	pre := d
	for cur != nil {
		pre.Next = cur.Next
		cur.Next = cur.Next.Next
		cur = cur.Next
		pre = pre.Next
	}
	head.Print()
	return d.Next
}
