package sliding_window

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestA(t *testing.T) {
	nums := []int{4, 3, 5, 4, 3, 3, 6, 7}
	fmt.Println(SlidingWindow(nums, 2))
}

func TestSlidingWindow(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		nums     []int
		k        int
		expected []int
	}{
		{
			nums:     []int{1, 2, 3, 1, 5},
			k:        1,
			expected: []int{1, 2, 3, 1, 5},
		},
		{
			nums:     []int{4, 3, 5, 4, 3, 3, 6, 7},
			k:        2,
			expected: []int{4, 5, 5, 4, 3, 6, 7},
		},
		{
			nums:     []int{4, 3, 5, 4, 3, 3, 6, 7},
			k:        3,
			expected: []int{5, 5, 5, 4, 6, 7},
		},
	}

	for _, testCase := range testCases {
		actual := SlidingWindow(testCase.nums, testCase.k)
		require.Equal(t, testCase.expected, actual)
	}
}

func TestB(t *testing.T) {
	nums := []int{5, 3, 3, 3, 3, 2, 2, 2, 2, 6}
	fmt.Println(SlidingWindow(nums, 2))
}
