package main

import (
	"fmt"
	"time"
)

func main() {
	//fmt.Println("*************** FanIn **************")
	//FanIn()
	//time.Sleep(time.Second)
	fmt.Println("*************** FanOut **************")
	ch := Generate("hello")
	process := NewProcessor()
	for v := range ch {
		process.PostJob(v)
	}
	time.Sleep(time.Second)

}
