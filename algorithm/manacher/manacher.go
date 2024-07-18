package manacher

import (
	"fmt"
)

func MaxLcpsLength(s string) int {
	if len(s) == 0 {
		return 0
	}
	res, _ := manacher1(s)
	return res - 1
}

func manacherString(s string) []byte {
	res := make([]byte, 2*len(s)+1)
	index := 0
	for i := 0; i < len(res); i++ {
		if i&1 == 0 {
			res[i] = '#'
		} else {
			res[i] = s[index]
			index++
		}
	}
	return res
}

func manacher1(s string) (r, index int) {
	str := manacherString(s)
	d := make([]int, len(str))
	R, C := -1, 0 // R 代表右边界, C 代表中心点
	for i := 0; i < len(str); i++ {
		if R < i {
			//暴力扩
			for i-d[i] >= 0 && i+d[i] < len(str) && str[i-d[i]] == str[i+d[i]] {
				d[i]++
			}
		} else {
			// L = 2 *C - R
			// i' = 2 *C - i   l' = i' - d[i'] + 1
			if 2*C-i-d[2*C-i]+1 > 2*C-R { // i' 的范围完全在 [L...R] 里面
				d[i] = d[2*C-i]
			} else if 2*C-i-d[2*C-i]+1 < 2*C-R { //  i' 的范围部分不在 [L...R] 里面
				d[i] = R - i + 1
			} else { // i' 左边界 l' 恰好和 L 重合
				d[i] = R - i
				for i-d[i] >= 0 && i+d[i] < len(str) && str[i-d[i]] == str[i+d[i]] {
					d[i]++
				}
			}
		}
		if r < d[i] {
			r = d[i]
			index = i
		}
		if C+d[i]-1 >= R {
			C = i
			R = C + d[i] - 1
		}
		fmt.Printf("R = %d, C = %d\n", R, C)
	}
	fmt.Println(d)
	return
}

// 返回最长回文长度
func manacher(s string) int {
	str := manacherString(s)
	d := make([]int, len(str))      // 回文半径数组
	var res = -1                    // 保存最大值
	R, C := -1, 0                   // R 回文边界, C 中心位置
	for i := 0; i < len(str); i++ { // 每一个位置都求回文半径
		// [i] 至少回文半径区域，先赋值给 d[i]
		if R <= i {
			d[i] = 1
		} else {
			// i'位置 C - (i -C) = 2 * C -i
			d[i] = min(d[2*C-i], R-i)
		}

		for i-d[i] >= 0 && i+d[i] < len(str) && str[i-d[i]] == str[i+d[i]] {
			d[i]++
		}
		if i+d[i]-1 > R {
			C = i
			R = i + d[i] - 1
		}
		fmt.Printf("R = %d, C = %d\n", R, C)
		res = max(res, d[i])
	}
	fmt.Println(d)
	return res - 1
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func manacher2(s string) int {
	str := manacherString(s)
	d := make([]int, len(str))
	R, C := -1, -1 // R 最左边边界, C 中心点
	res := -1
	for i := 0; i < len(str); i++ {
		// 说明 i == R 的情况，两者都可以
		if R <= i { // i 在 R 外面
			d[i] = 1 // 可写可不写，不写就要和自己比一次
			for i-d[i] >= 0 && i+d[i] < len(str) && str[i-d[i]] == str[i+d[i]] {
				d[i]++
			}
		} else { // i  属于 [L...R] 区间
			// i' (2*C - i) 左右边界 [2*C - i - d[2 *C -i] + 1, 2 *C - i + d[2*C - i] -1 ]
			// C 左右边界[2*C - R, R]
			if 2*C-i-d[2*C-i]+1 > 2*C-R { // i' 属于 C 里面
				d[i] = d[2*C-i]
			} else if 2*C-i-d[2*C-i]+1 < 2*C-R { // i' 左边界不属于 C
				d[i] = R - i + 1
			} else {
				// i' 恰好和 L 重合
				d[i] = R - i + 1
				for i-d[i] >= 0 && i+d[i] < len(str) && str[i-d[i]] == str[i+d[i]] {
					d[i]++
				}
			}
		}
		if i+d[i]-1 > R {
			R = i + d[i] - 1
			C = i
		}
		fmt.Printf("R = %d, C = %d\n", R, C)
		res = max(res, d[i])
	}
	fmt.Println(d)
	return res - 1
}

//func manacher(s string)  {
//	str := manacherString(s)
//	pArr := make([]int, len(str))
//	R, C := 0, -1
//	for i := 0; i < len(str);i++ {
//		if i 在 R 外部 {
//			 从 i 开始往两边暴力扩
//		} else {
//			if i‘ 回文区域彻底在 L...R 内 {
//				pArr[i] = O(1) 表达式
//			} else if i'' 回文区域有一部分在 L... R 外{
//				pArr[i] = O(1) 表达式
//			} else { i'' 回文区域和 L...R 的左边界压线
//				从 R 之外的字符开始，往外扩，然后确定 pArr[i] 的答案
//
//			}
//		}
//
//		更新 R,C
//	}
//}
