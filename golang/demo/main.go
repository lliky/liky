package main

import "fmt"

func main() {
	var a = make(map[string]int)
	exp := a["bc"]
	a["bc"] = a["bc"] + 1
	fmt.Println(exp, a)
}
