package monotonous_stack

import (
	"golang.org/x/exp/rand"
	"time"
)

// 在数组中想找到一个数，左边和右边比这个数小、且离这个数最近的位置

func MonotonousStackNoRepeat(nums []int) [][2]int {
	res := make([][2]int, len(nums)) // [2]int, [0]表示左边，[1]表示右边
	stack := make([]int, 0)
	for i := 0; i < len(nums); i++ {
		for len(stack) > 0 && nums[stack[len(stack)-1]] < nums[i] {
			index := stack[len(stack)-1]
			stack = stack[:len(stack)-1] // pop
			if len(stack) > 0 {
				res[index][0] = nums[stack[len(stack)-1]]
			} else {
				res[index][0] = -1
			}
			res[index][1] = nums[i]
		}
		stack = append(stack, i)
	}
	for len(stack) > 0 {
		index := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if len(stack) > 0 {
			res[index][0] = nums[stack[len(stack)-1]]
		} else {
			res[index][0] = -1
		}
		res[index][1] = -1
	}
	return res
}

func MonotonousStackRepeat(nums []int) [][2]int {
	res := make([][2]int, len(nums))
	stack := make([][]int, 0)
	for i := 0; i < len(nums); i++ {
		for len(stack) > 0 && nums[stack[len(stack)-1][0]] < nums[i] {
			indexs := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			var leftIndex int
			if len(stack) > 0 {
				leftIndex = nums[stack[len(stack)-1][0]]
			} else {
				leftIndex = -1
			}
			for _, index := range indexs {
				res[index][0] = leftIndex
				res[index][1] = nums[i]
			}
		}
		if len(stack) > 0 && nums[stack[len(stack)-1][0]] == nums[i] {
			stack[len(stack)-1] = append(stack[len(stack)-1], i)
		} else {
			tmp := []int{i}
			stack = append(stack, tmp)
		}
	}
	for len(stack) > 0 {
		indexs := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		var leftIndex int
		if len(stack) > 0 {
			leftIndex = nums[stack[len(stack)-1][0]]
		} else {
			leftIndex = -1
		}
		for _, index := range indexs {
			res[index][0] = leftIndex
			res[index][1] = -1
		}
	}
	return res
}

func getRandomArrayNoRepeat(size int) []int {
	rand.Seed(uint64(time.Now().UnixNano()))
	nums := make([]int, rand.Intn(size))
	for i := 0; i < len(nums); i++ {
		nums[i] = i
	}
	for i := 0; i < len(nums); i++ {
		swapIndex := rand.Intn(size)
		nums[swapIndex], nums[i] = nums[i], nums[swapIndex]
	}
	return nums
}
