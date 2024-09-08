package bitwise_operation

// 保证参数 n，不是 1 就是 0
// 1 -> 0
// 0 -> 1
func Flip(n int) int {
	return n ^ 1
}

// n 是非负数，返回 1
// n 是负数返回 0
func Sign(n int) int {
	return Flip(n>>63) & 1
}

// 数据可能溢出
func GetMax1(a, b int) int {
	c := a - b
	scA := Sign(c)   // a - b 为非负， scA 为 1；a - b 为负，scA 为 0
	scB := Flip(scA) // scA 为 0，scB 为 1；scA 为1， scB 为 0
	return scA*a + scB*b
}

func GetMax2(a, b int) int {
	c := a - b
	sa := Sign(a)
	sb := Sign(b)
	sc := Sign(c)
	difSab := sa ^ sb       // a 和 b 符号不一样，为 1；一样，为 0
	sameSab := Flip(difSab) // a 和 b 符号一样，为 1；不一样，为 0
	// 返回 a 的条件
	//	1. a, b 符号一样，且 a - b >= 0
	//	2. a, b 符号不一样，a >= 0
	returnA := sameSab*sc + difSab*sa
	returnB := Flip(returnA)
	return returnA*a + returnB*b
}

func Is2Power(n int) bool {
	return n&(n-1) == 0
}

func Is4Power(n int) bool {
	return (n&(n-1) == 0) && (n&0x55555555 != 0)
}
