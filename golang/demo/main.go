package main

import (
	"fmt"
	"time"
)

type AA struct {
	Name string
}

func main() {
	fmt.Println(time.Now().Add(6 * time.Minute))
}

func alignUp(n, a uintptr) uintptr {
	return (n + a - 1) &^ (a - 1)
}
