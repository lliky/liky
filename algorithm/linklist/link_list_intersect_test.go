package linklist

import "testing"

func TestGetIntersectNode(t *testing.T) {
	testCases := []struct {
		name string
		f    func() (*Node, *Node)
	}{
		{
			name: "no loop no intersect",
			f:    NewNoLoopNoIntersect,
		},
		{
			name: "no loop intersect",
			f:    NewNoLoopIntersect,
		},
		{
			name: "loop no intersect",
			f:    NewLoopNoIntersect,
		},
		{
			name: "loop same entry point intersect",
			f:    NewLoopSameEntryPointIntersect,
		},
		{
			name: "loop different entry point  intersect",
			f:    NewLoopDifferentEntryPointIntersect,
		},
	}
	for _, testCase := range testCases {
		headA, headB := testCase.f()
		t.Logf("\ntestCase name: %s\nlink list intersect: %v \n", testCase.name, GetIntersectNode(headA, headB))
	}
}

func NewNoLoopNoIntersect() (*Node, *Node) {
	headA := &Node{Val: 1}
	headA.Next = &Node{Val: 2}
	headA.Next.Next = &Node{Val: 3}
	headA.Next.Next.Next = &Node{Val: 4}

	headB := &Node{Val: 1}
	headB.Next = &Node{Val: 2}
	headB.Next.Next = &Node{Val: 3}
	headB.Next.Next.Next = &Node{Val: 4}
	return headA, headB
}
func NewNoLoopIntersect() (*Node, *Node) {
	headA := &Node{Val: 1}
	headA.Next = &Node{Val: 2}
	headA.Next.Next = &Node{Val: 3}
	headA.Next.Next.Next = &Node{Val: 4}
	headA.Print()
	headB := &Node{Val: 6}
	headB.Next = &Node{Val: 5}
	headB.Next.Next = &Node{Val: 7}
	headB.Next.Next.Next = headA.Next
	headB.Print()
	return headA, headB
}

func NewLoopNoIntersect() (*Node, *Node) {
	headA := &Node{Val: 1}
	headA.Next = &Node{Val: 2}
	headA.Next.Next = &Node{Val: 3}
	headA.Next.Next.Next = headA
	headA.Print()
	headB := &Node{Val: 6}
	headB.Next = &Node{Val: 5}
	headB.Next.Next = &Node{Val: 7}
	return headA, headB
}

func NewLoopSameEntryPointIntersect() (*Node, *Node) {
	headA := &Node{Val: 1}
	headA.Next = &Node{Val: 2}
	headA.Next.Next = &Node{Val: 3}
	headA.Next.Next.Next = &Node{Val: 4}
	headA.Next.Next.Next.Next = &Node{Val: 5}
	headA.Next.Next.Next.Next.Next = &Node{Val: 6}
	headA.Next.Next.Next.Next.Next.Next = headA.Next.Next
	headA.Print()
	headB := &Node{Val: 10}
	headB.Next = &Node{Val: 11}
	headB.Next.Next = &Node{Val: 12}
	headB.Next.Next.Next = headA.Next
	return headA, headB
}

func NewLoopDifferentEntryPointIntersect() (*Node, *Node) {
	headA := &Node{Val: 1}
	headA.Next = &Node{Val: 2}
	headA.Next.Next = &Node{Val: 3}
	headA.Next.Next.Next = &Node{Val: 4}
	headA.Next.Next.Next.Next = &Node{Val: 5}
	headA.Next.Next.Next.Next.Next = &Node{Val: 6}
	headA.Next.Next.Next.Next.Next.Next = &Node{Val: 7}
	headA.Next.Next.Next.Next.Next.Next.Next = &Node{Val: 8}
	headA.Next.Next.Next.Next.Next.Next = headA.Next.Next
	headB := &Node{Val: 10}
	headB.Next = &Node{Val: 11}
	headB.Next.Next = &Node{Val: 12}
	headB.Next.Next.Next = headA.Next.Next.Next.Next.Next
	return headA, headB
}
