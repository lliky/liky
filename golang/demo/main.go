package main

import "fmt"

func main() {
	ch1 := make(chan int)
	ch1 <- 1
	select {
	case v := <-ch1:
		fmt.Println(v)
	default:
		fmt.Println("default")
	}
}
