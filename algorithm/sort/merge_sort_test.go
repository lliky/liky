package sort

import (
	"sort"
	"testing"
)

func TestMergeSort(t *testing.T) {
	testTimes := 500000

	for i := 0; i < testTimes; i++ {
		arr1 := GenerateRandomArray(100, 100)
		arr2 := make([]int, len(arr1))
		copy(arr2, arr1)
		MergeSort(arr1)
		sort.Sort(SortArr{arr2})
		if !isEqual(arr1, arr2) {
			t.Logf("arr1: %v\n", arr1)
			t.Logf("arr2: %v\n", arr2)
			t.Fatalf("failed")
		}
	}
	t.Logf("success")
}
