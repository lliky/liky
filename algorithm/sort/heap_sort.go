package sort

// HeapSort
/*
	1.让数组成为大根堆
	i.可以 heapInsert
	ii.可以从中间位置开始 heapify
	2. 第一位数和堆最后一个交换，heapSize--, heapify 过程
*/
func HeapSort(nums []int) {
	if len(nums) < 2 {
		return
	}
	// 成为大根堆

	//for i := range nums { // O(N)
	//	HeapInsert(nums, i) // O(logN)
	//}
	for i := (len(nums) - 1) / 2; i >= 0; i-- {
		Heapify(nums, i, len(nums))
	}
	// 排序过程
	heapSize := len(nums)
	for heapSize > 0 { // O(N)
		swap2(nums, 0, heapSize-1) //O(1)
		heapSize--                 //O(1)
		Heapify(nums, 0, heapSize) //O(logN)
	}
}

// 大根堆
// 某个数出现在 index 位置，往上继续移动
// 时间复杂度O(logN)
func HeapInsert(nums []int, index int) {
	for nums[index] > nums[(index-1)/2] {
		swap2(nums, index, (index-1)/2)
		index = (index - 1) / 2
	}
}

// 某个数在 index 位置，能否往上移动
func Heapify(nums []int, index int, heapSize int) {
	left := 2*index + 1
	for left < heapSize { // 说明 index 下方还有孩子，heapSize 始终比最后一个 index 大1
		largest := left
		if left+1 < heapSize && nums[left+1] > nums[left] { //
			// 如果存在右孩子，并且右孩子的值比左孩子的值大，下标就给 largest
			largest = left + 1
		}
		// 父的值和最大孩子的值比较，谁大, 把下标给 largest
		if nums[largest] < nums[index] {
			//largest = index
			break // break 说明父的值比孩子都大，结束 heapify 过程
		}
		swap2(nums, largest, index)
		index = largest
		left = 2*index + 1
	}
}

// 堆排序扩展题目
/*
	已知一个几乎有序的数组，几乎有序是指，如果把数组排好顺序的话，每个元素移动的距离可以不超过 k, 并且 k 相对于
数组来说比较小。请选择一个合适的排序算法针对这个数据进行排序。
*/

func SortArrayDistancelessK(nums []int, k int) {
	heap := NewPriorityQueue()
	index := 0
	for ; index < min(len(nums), k); index++ {
		heap.Add(nums[index])
	}
	i := 0
	for ; index < len(nums); index++ {
		heap.Add(nums[index])
		nums[i] = heap.Pop()
		i++
	}
	for !heap.IsEmpty() {
		nums[i] = heap.Pop()
		i++
	}
}
