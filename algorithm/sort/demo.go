package sort

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
