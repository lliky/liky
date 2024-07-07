package hash_table

import (
	"fmt"
	"unsafe"
)

func systemIntBits() {
	var x int
	fmt.Println("int 位数为: ", unsafe.Sizeof(x)*8)

	a := 15
	fmt.Printf("a bit : %b, %v\n", a, a)
	a &= ^(1 << 2)
	fmt.Printf("a bit : %b, %v\n", a, a)
}

func bits() {
	arr := make([]int, 10) // 64bit * 10 -> 640bits

	// arr[0] int 0 ~ 63
	// arr[1] int 64 ~ 127
	// arr[2] int 128 ~ 191

	i := 178 // 想要获取第 178 位 bit 的状态

	numIndex := i / 64
	bitIndex := i % 64

	// 获取第 i 位 bit 的状态
	_ = (arr[numIndex] >> bitIndex) & 1

	// 将第 i 位 bit 置为 1
	arr[numIndex] = arr[numIndex] | (1 << bitIndex)

	// 将第 i 位 bit 置为 0
	arr[numIndex] = arr[numIndex] & ^(1 << bitIndex)
}
