package manacher

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestManacher(t *testing.T) {
	s1 := "1221"
	fmt.Println(string(manacherString(s1)))
	s2 := "121"
	fmt.Println(string(manacherString(s2)))
}

func TestManacher1(t *testing.T) {
	s := "aba"
	fmt.Println(manacher1(s))
}

func TestMaxLcpsLength(t *testing.T) {
	testCases := []struct {
		s        string
		expected int
	}{
		{
			s:        "aaaaa",
			expected: 5,
		},
		{
			s:        "abcdefg",
			expected: 1,
		},
		{
			s:        "123216",
			expected: 5,
		},
		{
			s:        "1221",
			expected: 4,
		},
	}
	for _, testCase := range testCases {
		actual := MaxLcpsLength(testCase.s)
		require.Equal(t, testCase.expected, actual)
	}
}

func TestMaxLcpsLength2(t *testing.T) {
	s := "aba"
	//1 2 1 4 1 2 1
	fmt.Println(manacher(s))
}

func TestMaxLcpsLength3(t *testing.T) {
	s := "aba"
	//1 2 1 4 1 2 1
	fmt.Println(manacher2(s))
}
