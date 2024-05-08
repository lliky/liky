package slice

import (
	"fmt"
	"testing"
)

func Test_1(t *testing.T) {
	s := make([]int, 10)
	s = append(s, 10)
	t.Logf("s: %v, len of s: %d, cap of s: %d", s, len(s), cap(s))
}

func Test_2(t *testing.T) {
	s := make([]int, 0, 10)
	s = append(s, 10)
	t.Logf("s: %v, len of s: %d, cap of s: %d", s, len(s), cap(s))
}

func Test_3(t *testing.T) {
	s := make([]int, 10, 11)
	s = append(s, 10)
	t.Logf("s: %v, len of s: %d, cap of s: %d", s, len(s), cap(s))
}

func Test_4(t *testing.T) {
	s := make([]int, 10, 12)
	s1 := s[8:]

	// [0,0], 2, 4
	t.Logf("s1: %v, len of s1: %d, cap of s1: %d", s1, len(s1), cap(s1))
}

func Test_5(t *testing.T) {
	s := make([]int, 10, 12)
	s1 := s[8:9]
	// [0], 1, 4
	t.Logf("s1: %v, len of s1: %d, cap of s1: %d", s1, len(s1), cap(s1))
}

func Test_6(t *testing.T) {
	s := make([]int, 10, 12)
	s1 := s[8:]
	s1[0] = -1
	// [0 0 0 0 0 0 0 0 -1 0]
	t.Logf("s: %v", s)
}

func Test_7(t *testing.T) {
	s := make([]int, 10, 12)
	//v := s[10]
	// 是否会越界？
	// 是
	t.Logf("s: %v", s)
}

func Test_8(t *testing.T) {
	s := make([]int, 10, 12)
	s1 := s[8:]
	s1 = append(s1, []int{10, 11, 12}...)
	//v  := s[10]
	// 数组访问是否越界
	// 是
	t.Logf("s: %v", s)
}

func Test_9(t *testing.T) {
	s := make([]int, 10, 12)
	s1 := s[8:]
	changeSlice1(s1)
	// [0 0 0 0 0 0 0 0 -1 0]
	t.Logf("s: %v", s)
}

func Test_10(t *testing.T) {
	s := make([]int, 10, 12)
	s1 := s[8:]
	changeSlice2(s1)
	// [0 0 0 0 0 0 0 0 0 0], 10, 12
	t.Logf("s: %v, len of s: %d, cap of s: %d", s, len(s), cap(s))
	// [0 0], 2, 4 错误的
	t.Logf("s1: %v, len of s1: %d, cap of s1: %d", s1, len(s1), cap(s1))
}

func Test_11(t *testing.T) {
	s := []int{0, 1, 2, 3, 4}
	s = append(s[:2], s[3:]...)
	// [0,1,3,4], 4, 5
	t.Logf("s: %v, len of s: %d, cap of s: %d", s, len(s), cap(s))
	//v := s[4] // 越界
}

func Test_12(t *testing.T) {
	s := make([]int, 512)
	s = append(s, 1)
	// 512, 扩容
	t.Logf("s: %v, len of s: %d, cap of s: %d", s, len(s), cap(s))
}
func changeSlice1(s1 []int) {
	s1[0] = -1
}
func changeSlice2(s1 []int) {
	// [0 0], 2, 4
	fmt.Printf("s1: %v, len of s1: %d, cap of s1: %d\n", s1, len(s1), cap(s1))
	s1 = append(s1, 10)
	// [0 0 10] 3 4
	fmt.Printf("s1: %v, len of s1: %d, cap of s1: %d\n", s1, len(s1), cap(s1))
}
