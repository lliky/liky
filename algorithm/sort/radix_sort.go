package sort

import "math"

func RadixSort(nums []int) {
	if len(nums) < 2 {
		return
	}
	radixSort(nums, 0, len(nums)-1, maxBits(nums))
}

// nums[l...r] 排序
func radixSort(nums []int, l, r, digit int) {
	radix := 10
	i, j := 0, 0
	bucket := make([]int, r-l+1)
	for d := 1; d <= digit; d++ { // 有多少位就进入多少次
		// 10 个空间
		// count[0] 当前位(d位)是 0 的数字有多少
		// count[1] 当前位(d位)是 0~1 的数字有多少
		// count[2] 当前位(d位)是 0~2 的数字有多少
		// count[i] 当前位(d位)是 0~i 的数字有多少
		count := make([]int, radix)
		for i := l; i <= r; i++ {
			j := getDigit(nums[i], d)
			count[j]++
		}
		for i := 1; i < radix; i++ {
			count[i] += count[i-1]
		}
		for i = r; i >= l; i-- {
			j := getDigit(nums[i], d)
			bucket[count[j]-1] = nums[i]
			count[j]--
		}
		i, j = l, 0
		for i <= r {
			nums[i] = bucket[j]
			i++
			j++
		}
	}
}

func maxBits(nums []int) int {
	m := math.MinInt
	for _, v := range nums {
		m = max(m, v)
	}
	res := 0
	for m != 0 {
		res++
		m /= 10
	}
	return res
}

func getDigit(num, d int) int {
	var res int
	for d > 0 {
		res = num % 10
		num /= 10
		d--
	}
	return res
}
