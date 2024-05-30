package sort

import (
	"fmt"
	"sort"
	"testing"
)

func TestQuestion1(t *testing.T) {
	nums := []int{3, 5, 6, 4, 1, 5, 6, 7, 9, 5}
	Question1(nums, 5)
	fmt.Printf("nums: %v\n", nums)
}

func TestQuestion2(t *testing.T) {
	nums := []int{1, 2, 5, 3, 4, 5, 6}
	Question2(nums, 5)
	fmt.Println(nums)
}

func TestQuickSort(t *testing.T) {
	testTimes := 500000

	for i := 0; i < testTimes; i++ {
		arr1 := GenerateRandomArray(1000, 100)
		arr2 := make([]int, len(arr1))
		copy(arr2, arr1)
		QuickSort(arr1)
		sort.Sort(SortArr{arr2})
		if !isEqual(arr1, arr2) {
			t.Logf("arr1: %v\n", arr1)
			t.Logf("arr2: %v\n", arr2)
			t.Fatalf("failed")
		}
	}
	t.Logf("success")
}
