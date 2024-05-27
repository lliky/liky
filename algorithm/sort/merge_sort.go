package sort

func MergeSort(arr []int) {
	if len(arr) < 2 {
		return
	}
	process(arr, 0, len(arr)-1)
}

// master 公式
// T(N) = 2 * T(N/2) + O(N)
func process(arr []int, L, R int) {
	if L == R {
		return
	}
	mid := L + (R-L)>>1
	process(arr, L, mid)
	process(arr, mid+1, R)
	merge(arr, L, mid, R)
}

func merge(arr []int, L, M, R int) {
	var tmp = make([]int, 0, R-L+1)
	p1, p2 := L, M+1
	for p1 <= M && p2 <= R {
		if arr[p1] <= arr[p2] {
			tmp = append(tmp, arr[p1])
			p1++
		} else {
			tmp = append(tmp, arr[p2])
			p2++
		}
	}
	for p1 <= M {
		tmp = append(tmp, arr[p1])
		p1++
	}
	for p2 <= R {
		tmp = append(tmp, arr[p2])
		p2++
	}
	for i, v := range tmp {
		arr[L+i] = v
	}
}
