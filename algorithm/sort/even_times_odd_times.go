package sort

import "fmt"

func PrintOddTimesNum1(arr []int) {
	var eor = 0
	for _, v := range arr {
		eor ^= v
	}
	fmt.Printf("the result is: %v\n", eor)
}

func PrintOddTimesNum2(arr []int) {
	var eor int
	for _, v := range arr {
		eor ^= v
	}
	right := eor & -eor
	var other int
	for _, v := range arr {
		if right&v == 0 {
			other ^= v
		}
	}
	fmt.Printf("the result is: %d, %d\n", other, other^eor)
}
