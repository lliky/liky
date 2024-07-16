package manacher

func MaxLcpsLength(s string) int {
	if len(s) == 0 {
		return 0
	}
	return 0
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

func manacher(s string)  {
	str := manacherString(s)
	pArr := make([]int, len(str))
	R, C := 0, -1
	for i := 0; i < len(str);i++ {
		if i 在 R 外部 {
			 从 i 开始往两边暴力扩
		} else {
			if i‘ 回文区域彻底在 L...R 内 {
				pArr[i] = O(1) 表达式
			} else if i'' 回文区域有一部分在 L... R 外{
				pArr[i] = O(1) 表达式
			} else { i'' 回文区域和 L...R 的左边界压线
				从 R 之外的字符开始，往外扩，然后确定 pArr[i] 的答案

			}
		}

		更新 R,C 
	}
}
