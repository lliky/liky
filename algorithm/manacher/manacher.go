package manacher

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
	}
	return
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
