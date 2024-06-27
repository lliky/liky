package greedy

import "fmt"

func Hanoi(n int) {
	if n > 0 {
		processHanoi(n, "左", "右", "中")
	}
}

// 1～i圆盘的目标是 from -> to，other 是另外一个

// 1 ～ i-1 from -> other
// i  from -> to
// 1 ~i-1 other -> to
func processHanoi(i int, start, end, other string) {
	if i == 1 {
		fmt.Println("Move 1 from", start, "to", end)
	} else {
		processHanoi(i-1, start, other, end)
		fmt.Println("Move", i, "from", start, "to", end)
		processHanoi(i-1, other, end, start)
	}
}
