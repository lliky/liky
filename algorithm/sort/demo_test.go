package sort

import (
	"fmt"
	"sort"
	"testing"
)

func TestDemoMergeSort(t *testing.T) {
	testTimes := 500000
	for i := 0; i < testTimes; i++ {
		arr1 := GenerateRandomArray(100, 100)
		arr2 := make([]int, len(arr1))
		copy(arr2, arr1)
		InsertSort(arr1)
		sort.Sort(SortArr{arr2})
		if !isEqual(arr1, arr2) {
			t.Logf("arr1: %v\n", arr1)
			t.Logf("arr2: %v\n", arr2)
			t.Fatalf("failed")
		}
	}
	t.Logf("success")
}

func TestDemoInsert(t *testing.T) {
	testTimes := 50000
	for i := 0; i < testTimes; i++ {
		arr1 := GenerateRandomArray(1000, 100)
		arr2 := make([]int, len(arr1))
		copy(arr2, arr1)
		DemoSelect(arr1)
		sort.Sort(SortArr{arr2})
		if !isEqual(arr1, arr2) {
			t.Logf("arr1: %v\n", arr1)
			t.Logf("arr2: %v\n", arr2)
			t.Fatalf("failed")
		}
	}
	t.Logf("success")
}

func TestDemoQuickSort(t *testing.T) {
	testTimes := 50000
	for i := 0; i < testTimes; i++ {
		arr1 := GenerateRandomArray(100, 100)
		arr2 := make([]int, len(arr1))
		copy(arr2, arr1)
		DemoQuickSort(arr1)
		sort.Sort(SortArr{arr2})
		if !isEqual(arr1, arr2) {
			t.Logf("arr1: %v\n", arr1)
			t.Logf("arr2: %v\n", arr2)
			t.Fatalf("failed")
		}
	}
	t.Logf("success")
}

func TestA(t *testing.T) {
	heap := NewPriorityQueue()
	heap.Add(1)
	heap.Add(4)
	fmt.Println(heap.Pop())
	heap.Add(2)
	heap.Add(5)
	heap.Add(3)
	heap.Add(7)
	for !heap.IsEmpty() {
		fmt.Println(heap.Pop())
	}
}
