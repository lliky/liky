package sort

import (
	"golang.org/x/exp/rand"
	"time"
)

func swap1(arr []int, i, j int) {
	arr[i], arr[j] = arr[j], arr[i]
}

func swap2(arr []int, i, j int) {
	if i == j { // i, j 不能指向同一个内存地址，否则变为 0
		return
	}
	arr[i] = arr[i] ^ arr[j]
	arr[j] = arr[i] ^ arr[j]
	arr[i] = arr[i] ^ arr[j]
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func GenerateRandomArray(maxSize, maxValue int) []int {
	rand.Seed(uint64(time.Now().UnixNano()))
	arr := make([]int, rand.Intn(maxSize+1))
	for i := 0; i < len(arr); i++ {
		arr[i] = rand.Intn(maxValue+1) - rand.Intn(maxValue)
	}
	return arr
}

func isEqual(arr1, arr2 []int) bool {
	if (arr1 == nil && arr2 != nil) || (arr1 != nil && arr2 == nil) {
		return false
	}
	if arr1 == nil && arr2 == nil {
		return true
	}
	if len(arr1) != len(arr2) {
		return false
	}
	for i := 0; i < len(arr1); i++ {
		if arr1[i] != arr2[i] {
			return false
		}
	}
	return true
}

type SortArr struct {
	Arr []int
}

func (a SortArr) Len() int {
	return len(a.Arr)
}

func (a SortArr) Less(i, j int) bool {
	return a.Arr[i] < a.Arr[j]
}

func (a SortArr) Swap(i, j int) {
	a.Arr[i], a.Arr[j] = a.Arr[j], a.Arr[i]
}

type PriorityQueue struct {
	arr      []int
	heapSize int
}

func NewPriorityQueue() *PriorityQueue {
	return &PriorityQueue{
		arr: make([]int, 0),
	}
}

func (p *PriorityQueue) Add(num int) {
	p.arr = append(p.arr, num)
	p.heapSize++
	priorityQueueHeapInsert(p.arr, p.heapSize-1)
}

func (p *PriorityQueue) Pop() int {
	res := p.arr[0]
	swap2(p.arr, 0, len(p.arr)-1)
	p.arr = p.arr[:len(p.arr)-1]
	p.heapSize--
	priorityQueueHeapify(p.arr, 0, p.heapSize)
	return res
}

func (p *PriorityQueue) IsEmpty() bool {
	return p.heapSize == 0
}

func priorityQueueHeapInsert(nums []int, index int) {
	for nums[index] > nums[(index-1)/2] {
		swap2(nums, index, (index-1)/2)
		index = (index - 1) / 2
	}
}

func priorityQueueHeapify(nums []int, index, heapSize int) {
	left := 2*index + 1
	for left < heapSize { // 说明有孩子
		largest := left
		if left+1 < heapSize && nums[left+1] > nums[left] {
			largest = left + 1
		}
		if nums[index] >= nums[largest] {
			break
		}
		swap2(nums, index, largest)
		index = largest
		left = 2*index + 1
	}

}
