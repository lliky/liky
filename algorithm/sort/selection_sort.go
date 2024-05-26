package sort

func SelectionSort(arr []int) {
	if len(arr) < 2 {
		return
	}
	for i := 0; i < len(arr)-1; i++ { // i ~ N-1
		index := i
		for j := i + 1; j < len(arr); j++ { // i ~ N-1 上找最小的下标
			if arr[j] < arr[index] {
				index = j
			}
		}
		swap1(arr, i, index)
	}
}
