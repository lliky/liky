package sliding_window

// SlidingWindow leetcode239
func SlidingWindow(nums []int, k int) []int {
	if k < 1 || len(nums) < k {
		return nil
	}
	res := make([]int, len(nums)-k+1)
	index := 0
	qMax := make([]int, 0)
	for i := 0; i < len(nums); i++ {
		for len(qMax) > 0 && nums[qMax[len(qMax)-1]] <= nums[i] {
			qMax = qMax[:len(qMax)-1]
		}
		qMax = append(qMax, i)
		//if len(qMax) == k+1 { 这里不能以 队列的长度来判断，可参考 TestB
		if qMax[0] == i-k { // 说明第一个最大数超出滑动窗口的范围
			qMax = qMax[1:]
		}
		if i >= k-1 {
			res[index] = nums[qMax[0]]
			index++
		}
	}
	return res
}
