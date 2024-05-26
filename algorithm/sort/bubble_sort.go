package sort

func BubbleSort(arr []int) {
	if len(arr) < 2 {
		return
	}
	for i := len(arr) - 1; i > 0; i-- {
		for j := 0; j < i; j++ {
			if arr[j] > arr[j+1] {
				swap2(arr, j, j+1)
			}
		}
	}
}
