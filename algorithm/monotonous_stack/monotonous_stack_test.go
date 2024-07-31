package monotonous_stack

import (
	"fmt"
	"testing"
)

func TestMonotonousStackNoRepeat(t *testing.T) {
	nums := []int{7, 5, 6, 3, 4, 8, 9}
	fmt.Println(MonotonousStackNoRepeat(nums))
}

func TestMonotonousStackRepeat(t *testing.T) {
	nums := []int{5, 7, 3, 1, 4, 8, 4, 1, 1}
	fmt.Println(MonotonousStackRepeat(nums))
}
