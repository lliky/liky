package greedy

import "math"

func Num1(n int) int {
	if n < 1 {
		return 0
	}
	record := make([]int, n) // record[i] -> i 行的皇后，放在了第几列
	return process1(0, record, n)
}

// i : 目前来到了第 i 行
// record[0...i-1]:表示 i 之前的行，放过的皇后（任意两个皇后一定不共行，不共列，不共斜线）
// n : 整体一共多少行
// 返回值摆完所有的皇后，合理的摆法有多少种
func process1(i int, record []int, n int) int {
	if i == n { // 终止行
		return 1
	}
	var res = 0
	for j := 0; j < n; j++ { // 当前行在 i 行，尝试 i 行所有的列
		// 当前 i 行的皇后，放在 j 列，会不会和之前(0...i-1)的皇后，共行共列或共斜线
		// 如果是，无效
		// 如果不是，认为有效
		if isValid(record, i, j) {
			record[i] = j
			res += process1(i+1, record, n)
		}
	}
	return res
}

// record[0...i-1] 需要看，record[i...] 不需要
// 返回 i 行皇后，放在 j 列是否有效
func isValid(record []int, i, j int) bool {
	for k := 0; k < i; k++ {
		if j == record[k] || math.Abs(float64(record[k]-j)) == math.Abs(float64(i-k)) {
			return false
		}
	}
	return true
}
