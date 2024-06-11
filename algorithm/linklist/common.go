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
	node := findLoopNode(head)
	if node == nil {
		fmt.Println("The link list is no loop")
		cur := head
		for cur != nil {
			fmt.Printf(" %d ", cur.Val)
			cur = cur.Next
			if cur != nil {
				fmt.Printf("->")
			} else {
				fmt.Println()
			}
		}
	} else {
		fmt.Println("The link list is loop")
		cur := head
		var flag bool
		for {
			if cur == node {
				flag = true
			}
			fmt.Printf(" %d ->", cur.Val)
			cur = cur.Next
			if flag && cur == node {
				fmt.Printf("%d\n", cur.Val)
				break
			}
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

func swap(arr []*Node, i, j int) {
	arr[i], arr[j] = arr[j], arr[i]
}

type RandomNode struct {
	Val    int
	Next   *RandomNode
	Random *RandomNode
}

func (head *RandomNode) Print() {
	cur := head
	fmt.Println("*********** begin print **********")
	if cur == nil {
		fmt.Println("The link list is empty")
	}
	for cur != nil {
		fmt.Printf(" %d ", cur.Val)
		cur = cur.Next
		if cur != nil {
			fmt.Printf("->")
		} else {
			fmt.Println()
		}
	}

	fmt.Println("************ end print ***********")
}
