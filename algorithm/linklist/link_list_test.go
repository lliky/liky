package linklist

import (
	"fmt"
	"testing"
)

func TestReverseSingleLinkList1(t *testing.T) {
	head := &Node{Val: 1}
	head.Next = &Node{Val: 2}
	head.Next.Next = &Node{Val: 3}
	head.Next.Next.Next = &Node{Val: 4}
	head.Next.Next.Next.Next = &Node{Val: 5}
	head.Print()
	cur := ReverseSingleLinkList1(head)
	cur.Print()
}

func TestReverseSingleLinkList2(t *testing.T) {
	head := &Node{Val: 1}
	head.Next = &Node{Val: 2}
	head.Next.Next = &Node{Val: 3}
	head.Next.Next.Next = &Node{Val: 4}
	head.Next.Next.Next.Next = &Node{Val: 5}
	head.Print()
	cur := ReverseSingleLinkList2(head)
	cur.Print()
}

func TestReverseSingleLinkList3(t *testing.T) {
	head := &Node{Val: 1}
	head.Next = &Node{Val: 2}
	head.Next.Next = &Node{Val: 3}
	head.Next.Next.Next = &Node{Val: 4}
	head.Next.Next.Next.Next = &Node{Val: 5}
	head.Print()
	cur := ReverseSingleLinkList3(head)
	cur.Print()
}

func TestPrintLinkListCommonPart(t *testing.T) {
	head1 := &Node{Val: 1}
	head1.Next = &Node{Val: 2}
	head1.Next.Next = &Node{Val: 3}
	head1.Next.Next.Next = &Node{Val: 4}
	head1.Next.Next.Next.Next = &Node{Val: 5}
	head2 := &Node{Val: 2}
	head2.Next = &Node{Val: 2}
	head2.Next.Next = &Node{Val: 2}
	head2.Next.Next.Next = &Node{Val: 4}
	head2.Next.Next.Next.Next = &Node{Val: 6}
	PrintLinkListCommonPart(head1, head2)
}

func TestMergeTwoLists(t *testing.T) {
	head1 := &Node{Val: 1}
	head1.Next = &Node{Val: 2}
	head1.Next.Next = &Node{Val: 3}
	head1.Next.Next.Next = &Node{Val: 4}
	head1.Next.Next.Next.Next = &Node{Val: 5}
	head2 := &Node{Val: 2}
	head2.Next = &Node{Val: 2}
	head2.Next.Next = &Node{Val: 2}
	head2.Next.Next.Next = &Node{Val: 4}
	head2.Next.Next.Next.Next = &Node{Val: 6}
	head := MergeTwoLists(head1, head2)
	head.Print()
}

func TestIsPalindrome1(t *testing.T) {
	head1 := &Node{Val: 1}
	head1.Next = &Node{Val: 2}
	head1.Next.Next = &Node{Val: 3}
	head1.Next.Next.Next = &Node{Val: 2}
	head1.Next.Next.Next.Next = &Node{Val: 1}
	fmt.Println(IsPalindrome1(head1))
}

func TestIsPalindrome2(t *testing.T) {
	head1 := &Node{Val: 1}
	head1.Next = &Node{Val: 2}
	head1.Next.Next = &Node{Val: 2}
	head1.Next.Next.Next = &Node{Val: 1}
	//head1.Next.Next.Next.Next = &Node{Val: 2}
	//head1.Next.Next.Next.Next.Next = &Node{Val: 1}
	fmt.Println(IsPalindrome2(head1))
}

func TestPartitionByVal1(t *testing.T) {
	head := &Node{Val: 5}
	head.Next = &Node{Val: 1}
	//head.Next.Next = &Node{Val: 4}
	//head.Next.Next.Next = &Node{Val: 3}
	//head.Next.Next.Next.Next = &Node{Val: 2}
	//head.Next.Next.Next.Next.Next = &Node{Val: 4}
	//head.Next.Next.Next.Next.Next.Next = &Node{Val: 1}
	//head.Next.Next.Next.Next.Next.Next.Next = &Node{Val: 2}
	head.Print()
	PartitionByVal1(head, 3).Print()
}

func TestPartitionByVal2(t *testing.T) {
	head := &Node{Val: 5}
	head.Next = &Node{Val: 1}
	//head.Next.Next = &Node{Val: 4}
	//head.Next.Next.Next = &Node{Val: 2}
	//head.Next.Next.Next.Next = &Node{Val: 4}
	//head.Next.Next.Next.Next.Next = &Node{Val: 1}
	//head.Next.Next.Next.Next.Next.Next = &Node{Val: 2}
	head.Print()
	PartitionByVal2(head, 3).Print()
}

func TestCopyRandomList1(t *testing.T) {
	head := &RandomNode{Val: 7}
	head.Next = &RandomNode{Val: 13}
	head.Next.Next = &RandomNode{Val: 11}
	head.Next.Next.Next = &RandomNode{Val: 10}
	head.Next.Next.Next.Next = &RandomNode{Val: 1}

	head.Next.Random = head
	head.Next.Next.Random = head.Next.Next.Next.Next
	head.Next.Next.Next.Random = head.Next.Next
	head.Next.Next.Next.Next.Random = head
	head.Print()
	head = CopyRandomList1(head)
	head.Print()
}

func TestCopyRandomList2(t *testing.T) {
	head := &RandomNode{Val: 7}
	head.Next = &RandomNode{Val: 13}
	head.Next.Next = &RandomNode{Val: 11}
	head.Next.Next.Next = &RandomNode{Val: 10}
	head.Next.Next.Next.Next = &RandomNode{Val: 1}

	head.Next.Random = head
	head.Next.Next.Random = head.Next.Next.Next.Next
	head.Next.Next.Next.Random = head.Next.Next
	head.Next.Next.Next.Next.Random = head
	head.Print()
	head = CopyRandomList2(head)
	head.Print()
}
