package main

import "fmt"

func main() {
	fmt.Println(make(chan int))
	for i := 1; i < 16; i++ {
		and := i & 3
		fmt.Println(i, " & 3 = ", and)
	}
	fmt.Println(alignUp(12, 4))
}

func alignUp(n, a uintptr) uintptr {
	return (n + a - 1) &^ (a - 1)
}
