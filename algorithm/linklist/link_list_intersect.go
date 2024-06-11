package linklist

// 相交链表  leetcode 160 变式，链表可能有环
// 给两个单链表的头节点 headA 和 headB，找出两个链表相交的起始节点，若不存在，返回 nil

// 1. 两个不成环的相交
// 2. 一个成环，一个不成环   不可能存在这种交点
// 3. 两个都成环相交

func GetIntersectNode(headA, headB *Node) *Node {
	if headA == nil || headB == nil {
		return nil
	}
	loopA, loopB := findLoopNode(headA), findLoopNode(headB)

	if loopA == nil && loopB == nil { // 1
		return noLoop(headA, headB)
	}
	if loopA != nil && loopB != nil { // 3
		return bothLoop(headA, loopA, headB, loopB)
	}
	return nil // 3
}

// 两个没有环的链表，返回相交点
func noLoop(headA, headB *Node) *Node {
	curA, curB := headA, headB
	for curA != curB {
		if curA == nil {
			curA = headB
		} else {
			curA = curA.Next
		}

		if curB == nil {
			curB = headA
		} else {
			curB = curB.Next
		}
	}
	return curB
}

// 1. 各自成环，无相交
// 2. 同一入环点
// 3. 不同入环点
func bothLoop(headA, loopA, headB, loopB *Node) *Node {
	curA, curB := headA, headB
	if loopA == loopB { // 2
		for curA != curB {
			if curA == loopA {
				curA = headB
			} else {
				curA = curA.Next
			}
			if curB == loopA {
				curB = headA
			} else {
				curB = curB.Next
			}
		}
		return curA
	} else { // 1 and 3
		curA = loopA.Next
		for curA != loopA {
			if curA == loopB { // 3, 返回 loopA, loopB 都行
				return curA
			}
			curA = curA.Next
		}
		return nil // 1
	}
}

// 找链表有环的交点，如果不是环，返回 nil
func findLoopNode(head *Node) *Node {
	// 1. 快慢指针, 走到相同节点
	// 2. 快指针回到 head, 以慢指针速度走，快慢指针，下次遇到就是就是环的入口
	if head == nil || head.Next == nil || head.Next.Next == nil {
		return nil
	}
	slow := head.Next
	fast := head.Next.Next
	for slow != fast {
		if fast.Next == nil || fast.Next.Next == nil {
			return nil
		}
		slow = slow.Next
		fast = fast.Next.Next
	}
	fast = head
	for slow != fast {
		slow = slow.Next
		fast = fast.Next
	}
	return slow
}
