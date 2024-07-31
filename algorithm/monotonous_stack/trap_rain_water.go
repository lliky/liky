package monotonous_stack

func TrapRainWater(nums []int) int {
	var res int
	stack := make([]int, 0)
	for i := 0; i < len(nums); i++ {
		for len(stack) > 0 && nums[stack[len(stack)-1]] <= nums[i] {
			stack = stack[:len(stack)-1]
		}
	}
	return res
}
