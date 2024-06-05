package linklist

import "fmt"

type PrintList interface {
	Print()
}

type Node struct {
	Val  int
	Next *Node
}

func (head *Node) Print() {
	fmt.Println("*********** begin print **********")
	if head == nil {
		fmt.Println("The link list is empty")
	}
	for head != nil {
		fmt.Printf(" %d ", head.Val)
		head = head.Next
		if head != nil {
			fmt.Printf("->")
		} else {
			fmt.Println()
		}
	}

	fmt.Println("************ end print ***********")
}

type DoubleNode struct {
	Val  int
	Prev *DoubleNode
	Next *DoubleNode
}

func (d *DoubleNode) Print() {
	fmt.Println("*********** begin print **********")
	if d == nil {
		fmt.Println("The link list is empty")
	}
	for d.Prev != nil {
		d = d.Prev
	}
	for d != nil {
		fmt.Printf(" %d ", d.Val)
		d = d.Next
		if d != nil {
			fmt.Printf("->")
		} else {
			fmt.Println()
		}
	}
	fmt.Println("************ end print ***********")
}
