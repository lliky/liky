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
