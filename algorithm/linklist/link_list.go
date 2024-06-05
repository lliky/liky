package linklist

import "fmt"

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
