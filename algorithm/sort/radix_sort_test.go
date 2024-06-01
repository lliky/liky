package sort

import (
	"fmt"
	"testing"
)

func TestGetDigit(t *testing.T) {
	num := 1248
	fmt.Println(getDigit(num, 1))
	fmt.Println(getDigit(num, 2))
	fmt.Println(getDigit(num, 3))
	fmt.Println(getDigit(num, 4))
}

func TestRadixSort(t *testing.T) {
	num := []int{12, 16, 6, 100, 2333, 45}
	RadixSort(num)
	fmt.Println(num)
}
