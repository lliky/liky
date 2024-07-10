package kmp

// StrStr leetcode 28
func StrStr(s, substr string) int {
	if len(substr) < 1 || len(s) < len(substr) {
		return -1
	}
	i1, i2 := 0, 0
	next := getNext(substr)               // O(M)
	for i1 < len(s) && i2 < len(substr) { // O(N)
		if s[i1] == substr[i2] {
			i1++
			i2++
		} else if next[i2] == -1 { // substr 中比对的位置已经无法往前跳了
			i1++
		} else {
			i2 = next[i2]
		}
	}
	// 条件终止 i1 或 i2 越界
	if i2 == len(substr) {
		return i1 - i2
	}
	return -1
}

func getNext(s string) []int {
	if len(s) < 2 {
		return []int{-1}
	}
	next := make([]int, len(s))
	next[0] = -1
	next[1] = 0
	i, cn := 2, 0 // i 表示 next 的位置, cn 表示 0 ～ cn-1 位置前缀后缀最长, 也表示 next[cn] 这个位置
	for i < len(s) {
		if s[i-1] == s[cn] {
			cn++
			next[i] = cn
			i++
		} else if cn > 0 { // cn 位置的字符和 i-1 位置的字符配不上，接着找前缀后缀
			cn = next[cn]
		} else { // 前缀后缀为 0
			next[i] = 0
			i++
		}
	}

	return next
}
