package sort

import (
	"fmt"
	"sort"
	"testing"
)

func TestGetDigit(t *testing.T) {
	num := 1248
	fmt.Println(getDigit(num, 1))
	fmt.Println(getDigit(num, 2))
	fmt.Println(getDigit(num, 3))
	fmt.Println(getDigit(num, 4))
}

func TestRadixSort(t *testing.T) {
	testTimes := 50000

	for i := 0; i < testTimes; i++ {
		arr1 := GenerateRandomArray(100, 100)
		arr2 := make([]int, len(arr1))
		copy(arr2, arr1)
		RadixSort(arr1)
		sort.Sort(SortArr{arr2})
		if !isEqual(arr1, arr2) {
			t.Logf("arr1: %v\n", arr1)
			t.Logf("arr2: %v\n", arr2)
			t.Fatalf("failed")
		}
	}
	t.Logf("success")
}
