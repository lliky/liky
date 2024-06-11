package binary_tree

import "fmt"

// 打印折痕

func PrintAllFolds(n int) {
	processFold(1, n, true)
}

// i 是节点的层数，N 是一共几层，down == true 是 凹，down == false 是凸
func processFold(i, N int, down bool) {
	if i > N {
		return
	}
	processFold(i+1, N, true)
	if down {
		fmt.Printf(" 凹 ")
	} else {
		fmt.Printf(" %s ", "凸")
	}
	processFold(i+1, N, false)
}
