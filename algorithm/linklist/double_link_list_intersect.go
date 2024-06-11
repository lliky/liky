package linklist

import "math"

// leetcode 160
// 如果两链表相交，返回第一个节点，若不相交，返回 nil

// 1. loop1, loop2 都为 nil
// 2. loop1, loop2 有 nil, 有不为 nil  结论不存在
// 3. loop1, loop2 都不为空

func GetIntersectionNode(head1, head2 *Node) *Node {
	if head1 == nil || head2 == nil {
		return nil
	}
	loop1 := GetLoopNode(head1)
	loop2 := GetLoopNode(head2)
	if loop1 == nil && loop2 == nil {
		return NoLoop(head1, head2)
	}
	if loop1 != nil && loop2 != nil {
		return BothLoop(head1, loop1, head2, loop2)
	}
	return nil
}

// NoLoop 两个链表无环的情况
// 1. 各自走到终点统计长度及最后节点
// 2. 让长的先走长度的差值
// 3. 一起走判断节点是否相等并返回
func NoLoop(head1, head2 *Node) *Node {
	cur1, cur2, n := head1, head2, 0
	for cur1.Next != nil {
		n++
		cur1 = cur1.Next
	}
	for cur2.Next != nil {
		n--
		cur2 = cur2.Next
	}
	if cur1 != cur2 { // 最后一个节点不等，说明不相交
		return nil
	}
	if n > 0 { // 谁长谁是 cur1
		cur1 = head1
		cur2 = head2
	} else {
		cur1 = head2
		cur2 = head1
	}
	n = int(math.Abs(float64(n)))
	for n > 0 {
		cur1 = cur1.Next
		n--
	}
	for cur1 != cur2 {
		cur1 = cur1.Next
		cur2 = cur2.Next
	}
	return cur2
}

// BothLoop 两个有环链表，返回第一个相交节点，如果不相交，返回 nil
// 1. 各自成
// 2. 同一个入环点
// 3. 不一样的入环点
func BothLoop(head1, loop1, head2, loop2 *Node) *Node {
	var cur1, cur2 *Node
	if loop1 == loop2 { // 情况 2，用无环相交处理，终止条件是 loop1
		n := 0
		cur1, cur2 = head1, head2
		for cur1.Next != loop1 {
			n++
			cur1 = cur1.Next
		}
		for cur2.Next != loop1 {
			n--
			cur2 = cur2.Next
		}
		if n > 0 {
			cur1 = head1
			cur2 = head2
		} else {
			cur1 = head2
			cur2 = head1
		}
		n = int(math.Abs(float64(n)))
		for n > 0 {
			cur1 = cur1.Next
		}
		for cur1 != cur2 {
			cur1 = cur1.Next
			cur2 = cur2.Next
		}
		return cur1
	} else { // 情况 1,3
		cur1 = loop1.Next
		for cur1 != loop1 {
			if cur1 == loop2 { // 情况 3
				return loop1
			}
			cur1 = cur1.Next
		}
		return nil // 情况 1
	}
}

// GetLoopNode 找到链表第一个入环节点，如果无环，返回 nil
func GetLoopNode(head *Node) *Node {
	if head == nil || head.Next == nil || head.Next.Next == nil {
		return nil
	}
	slow, fast := head.Next, head.Next.Next
	for slow != fast {
		if fast.Next == nil || fast.Next.Next == nil {
			return nil
		}
		slow = slow.Next
		fast = fast.Next.Next
	}
	fast = head
	for slow != fast {
		fast = fast.Next
		slow = slow.Next
	}
	return slow
}
