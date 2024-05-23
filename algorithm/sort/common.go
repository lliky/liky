package sort

func swap1(arr []int, i, j int) {
	arr[i], arr[j] = arr[j], arr[i]
}

func swap2(arr []int, i, j int) {
	arr[i] = arr[i] ^ arr[j]
	arr[j] = arr[i] ^ arr[j]
	arr[i] = arr[i] ^ arr[j]
}
