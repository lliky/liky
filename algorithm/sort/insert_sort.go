package sort

func InsertSort(arr []int) {
	if len(arr) < 2 {
		return
	}
	for i := 0; i < len(arr); i++ { // 0~i 做到有序
		for j := i - 1; j >= 0 && arr[j+1] < arr[j]; j-- {
			swap2(arr, j, j+1)
		}
	}
}
