package priority_queue

import (
	heap2 "container/heap"
	"fmt"
	"testing"
)

func TestPriorityQueue_Len(t *testing.T) {
	d := &PriorityQueue{}
	heap2.Init(d)

	heap2.Push(d, 5)
	heap2.Push(d, 3)
	heap2.Push(d, 7)
	heap2.Push(d, 1)
	for d.Len() > 0 {
		fmt.Println(d.Pop())
	}
}
