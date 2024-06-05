package linklist

import "testing"

func TestPrint(t *testing.T) {
	head := &Node{Val: 1}
	head.Next = &Node{Val: 2}
	head.Next.Next = &Node{Val: 3}
	head.Print()
}
