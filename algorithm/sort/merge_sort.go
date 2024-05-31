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

// 归并排序的扩展
// 1. 小和问题
// 2. 逆序对
// 3. leetcode 315

/*
小和问题
在一个数组中，每一个数左边比当前数小的数累加起来，叫做这个数组的小和。求一个数组的小和
*/
func SumSmall(nums []int) int {
	if len(nums) < 2 {
		return 0
	}
	return processSumSmall(nums, 0, len(nums)-1)
}

// 在 nums[L...R] 排序，并且求小和
func processSumSmall(nums []int, l, r int) int {
	if l == r {
		return 0
	}
	mid := l + (r-l)>>1
	return processSumSmall(nums, l, mid) +
		processSumSmall(nums, mid+1, r) +
		mergeSumSmall(nums, l, mid, r)
}

func mergeSumSmall(nums []int, l, mid, r int) int {
	tmp := make([]int, 0, r-l+1)
	p1, p2 := l, mid+1
	sum := 0
	for p1 <= mid && p2 <= r {
		if nums[p1] < nums[p2] { // 和归并不同点在于，如果两个数相等，那么右边的数需要先放到 tmp 数组中
			sum += nums[p1] * (r - p2 + 1) // 如果左边数(p1)小于右边数(p2)，那么区间[p2,l] 的数都大于 [p1], 所以都得加起来
			tmp = append(tmp, nums[p1])
			p1++
		} else {
			tmp = append(tmp, nums[p2])
			p2++
		}
	}
	for p1 <= mid {
		tmp = append(tmp, nums[p1])
		p1++
	}
	for p2 <= mid {
		tmp = append(tmp, nums[p1])
		p2++
	}
	for i, v := range tmp {
		nums[l+i] = v
	}
	return sum
}

/*
	逆序对问题
	在一个数组中，左边的数如果比右边数大，则这两个数构成一个逆序对，请打印所有逆序对。
	比如：[1,3,2,5,0],逆序对：
	(1,0),(3,2),(3,0),(2,0),(5,0)
*/

func ReversePair(nums []int) int {
	if len(nums) < 2 {
		return 0
	}
	return processReversePair(nums, 0, len(nums)-1)
}

// 在 nums[l...r] 中找多少逆序对，并排序
func processReversePair(nums []int, l, r int) int {
	if l == r {
		return 0
	}
	mid := l + (r-l)>>1
	return processReversePair(nums, l, mid) +
		processReversePair(nums, mid+1, r) +
		mergeReversePair(nums, l, mid, r)
}

func mergeReversePair(nums []int, l, mid, r int) int {
	tmp := make([]int, 0, r-l+1)
	p1, p2 := l, mid+1
	res := 0
	for p1 <= mid && p2 <= r {
		if nums[p1] <= nums[p2] {
			tmp = append(tmp, nums[p1])
			p1++
		} else {
			tmp = append(tmp, nums[p2])
			res += mid - p1 + 1 // 这里是如果以 mid 为界，右边的数(p2)小于左边的数(p1)，那么区间[p1, mid] 所有数，都大于p2
			//var a = p1 // 如果打印，就在这里打印逆序对
			//for a <= mid {
			//	fmt.Printf("(%d,%d)\n", nums[a], nums[p2])
			//	a++
			//}
			p2++
		}
	}
	for p1 <= mid {
		tmp = append(tmp, nums[p1])
		p1++
	}
	for p2 <= r {
		tmp = append(tmp, nums[p2])
		p2++
	}
	for i, v := range tmp {
		nums[l+i] = v
	}
	return res
}
