package sort

import (
	"golang.org/x/exp/rand"
	"time"
)

func swap1(arr []int, i, j int) {
	arr[i], arr[j] = arr[j], arr[i]
}

// i, j 不能指向同一个内存地址，否则变为 0
func swap2(arr []int, i, j int) {
	if i == j {
		return
	}
	arr[i] = arr[i] ^ arr[j]
	arr[j] = arr[i] ^ arr[j]
	arr[i] = arr[i] ^ arr[j]
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
