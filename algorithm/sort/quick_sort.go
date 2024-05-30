package sort

import (
	"golang.org/x/exp/rand"
	"time"
)

/*
	荷兰国旗问题
	问题一
	给定一个数组 arr，和一个数 num，请把小于等于 num 的数放在数组的左边，大于 num 的数放在
数组的右边。要求额外空间复杂度 O(1), 时间复杂度 O(N)
	问题二(荷兰国旗问题)
	给定一个数组 arr，和一个数 num，请把小于 num 的数放在数组的左边，等于 num 的数放在数组的中间，
大于 num 的数放在数组的右边。要求额外空间复杂度 O(1), 时间复杂度 O(N)
*/

/*
1. [i] <= num, [i]与区间下一位交换，i++
2. [i] > num, i++
*/
func Question1(nums []int, num int) {
	index := -1
	for i := 0; i < len(nums); i++ {
		if nums[i] <= num {
			index++
			swap1(nums, index, i)
		}
	}
}

/*
1. [i] < num, [i] 与左区间下一位交换，左区间右扩, i++
2. [i] = num, i++
3. [i] > num, [i] 与右区间的前一位交换，右区间左扩, i 不变
终止条件：i == 右区间第一个, 循环条件就是 i < r
*/
func Question2(nums []int, num int) {
	l, r := -1, len(nums)
	for i := 0; i < r; {
		if nums[i] < num {
			l++
			swap1(nums, l, i)
			i++
		} else if nums[i] == num {
			i++
		} else {
			r--
			swap1(nums, i, r)
		}
	}
}

/*
	快排 1.0，
	以区间最后一个数划分，question 1
*/

/*
	快排 2.0
	以区间最后一个数划分，question 2 荷兰国旗问题
*/

/*
	1.0，2.0 最差时间复杂度是O(N^2), 原因：划分值很偏，导致退化 N^2
	快排 3.0
	随机选择一个数，然后和最后一个交换，question2
*/

func QuickSort(nums []int) {
	if len(nums) < 2 {
		return
	}
	quickSort(nums, 0, len(nums)-1)
}

// nums[l...r] 有序
func quickSort(nums []int, l, r int) {
	if l < r {
		rand.Seed(uint64(time.Now().UnixNano()))
		swap1(nums, l+rand.Intn(r-l+1), r)
		p := partition(nums, l, r)
		quickSort(nums, l, p[0]-1)
		quickSort(nums, p[1]+1, r)
	}
}

// 处理 arr[l..r]的函数
// 默认以 arr[r] 做划分
// 返回等于区域(左边界，右边界)，所以返回一个长度为 2 的数组 res
func partition(nums []int, l, r int) []int {
	less, more := l-1, r
	for l < more {
		if nums[l] < nums[r] {
			less++
			swap1(nums, less, l)
			l++
		} else if nums[l] > nums[r] {
			more--
			swap1(nums, more, l)
		} else {
			l++
		}
	}
	swap1(nums, r, more)
	return []int{less + 1, more}
}

// 时间复杂度O(N*logN) ，空间复杂度O(logN)
