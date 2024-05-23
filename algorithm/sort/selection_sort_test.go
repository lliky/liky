package sort

import (
	"testing"
)

func TestSelection_sort(t *testing.T) {
	arr := []int{2, 3, 4, 6, 2, 0, -1}
	t.Logf("arr before sort: %v\n", arr)
	Selection_sort(arr)
	t.Logf("arr after sort: %v\n", arr)
}
