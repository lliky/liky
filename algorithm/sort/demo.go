package sort

import (
	"math"
)

func DemoMergeSort(nums []int) {
	demoProcess(nums, 0, len(nums)-1)
}

func demoProcess(nums []int, L, R int) {
	if L == R {
		return
	}
	mid := L + (R-L)>>1
	demoProcess(nums, L, mid)
	demoProcess(nums, mid+1, R)
	demoMerge(nums, L, mid, R)
}

func demoMerge(nums []int, L, M, R int) {
	tmp := make([]int, 0, R-L+1)
	p1, p2 := L, M+1
	for p1 <= M && p2 <= R {
		if nums[p1] <= nums[p2] {
			tmp = append(tmp, nums[p1])
			p1++
		} else {
			tmp = append(tmp, nums[p2])
			p2++
		}
	}
	for p1 <= M {
		tmp = append(tmp, nums[p1])
		p1++
	}
	for p2 <= R {
		tmp = append(tmp, nums[p2])
		p2++
	}
	for i, v := range tmp {
		nums[L+i] = v
	}
}

func DemoSelect(nums []int) {
	if len(nums) < 2 {
		return
	}
	for i := 0; i < len(nums); i++ {
		minIndex := i
		for j := i + 1; j < len(nums); j++ {
			if nums[j] < nums[minIndex] {
				minIndex = j
			}
		}
		nums[i], nums[minIndex] = nums[minIndex], nums[i]
	}
}

func DemoQuickSort(nums []int) {
	if len(nums) < 2 {
		return
	}
	demoQuickSort(nums, 0, len(nums)-1)
}

func demoQuickSort(nums []int, l, r int) {
	if l >= r {
		return
	}
	p := demoPartition(nums, l, r)
	demoQuickSort(nums, l, p[0]-1)
	demoQuickSort(nums, p[1]+1, r)
}

func demoPartition(nums []int, l, r int) []int {
	less, more := l-1, r
	for l < more {
		if nums[l] < nums[r] {
			less++
			swap2(nums, less, l)
			l++
		} else if nums[l] > nums[r] {
			more--
			swap2(nums, l, more)
		} else {
			l++
		}
	}
	swap2(nums, more, r) // 把最后一个换一下
	return []int{less + 1, more}
}

func DemoRadix(nums []int) {
	radix := 10
	digit := demoMaxBits(nums)
	for d := 1; d <= digit; d++ { //最大位数，需要几进几出桶
		count := make([]int, radix)
		for i := 0; i < len(nums); i++ {
			j := demoGetDigits(nums[i], d)
			count[j]++
		}
		for i := 1; i < radix; i++ { // 前缀和
			count[i] += count[i-1]
		}
		bucket := make([]int, len(nums))
		for i := len(nums) - 1; i >= 0; i-- {
			j := demoGetDigits(nums[i], d)
			bucket[count[j]-1] = nums[i]
			count[j]--
		}
		for i, v := range bucket {
			nums[i] = v
		}
	}
}

func demoGetDigits(num, d int) int {
	var v, c int
	for num != 0 {
		v = num % 10
		num /= 10
		c++
		if c == d {
			break
		}
	}
	return v
}

func demoMaxBits(nums []int) int {
	var max = math.MinInt
	for _, v := range nums {
		if max < v {
			max = v
		}
	}
	var res = 0
	for max != 0 {
		max /= 10
		res++
	}
	return res
}
