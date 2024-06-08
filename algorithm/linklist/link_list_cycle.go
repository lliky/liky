package linklist

// HasCycle leetcode 141
func HasCycle(head *Node) bool {
	if head == nil || head.Next == nil || head.Next.Next == nil {
		return false
	}
	slow := head.Next
	fast := head.Next.Next
	for slow != fast {
		if fast.Next == nil || fast.Next.Next == nil {
			return false
		}
		slow = slow.Next
		fast = fast.Next.Next
	}
	return true
}

// FindFirstIntersectNode leetcode 142
// 找到链表第一个入环节点，如果无环，返回 nil
func FindFirstIntersectNode(head *Node) *Node {
	if head == nil || head.Next == nil || head.Next.Next == nil {
		return nil
	}
	slow, fast := head.Next, head.Next.Next
	for slow != fast { //找到相遇节点
		if fast.Next == nil || fast.Next.Next == nil {
			return nil
		}
		slow = slow.Next
		fast = fast.Next.Next
	}
	fast = head //  找第一个相遇点
	for slow != fast {
		slow = slow.Next
		fast = fast.Next
	}
	return slow
}
