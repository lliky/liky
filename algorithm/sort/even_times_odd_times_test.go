package sort

import "testing"

func TestPrintOddTimesNum1(t *testing.T) {
	arr := []int{1, 1, 2, 2, 3, 4, 4}
	PrintOddTimesNum1(arr)
}

func TestPrintOddTimesNum2(t *testing.T) {
	arr := []int{1, 1, 2, 2, 3, 4, 4, 5, 1, 5}
	PrintOddTimesNum2(arr)
}
