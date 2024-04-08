package main

import (
	"fmt"
	"time"
)

func Generate(data string) <-chan string {
	ch := make(chan string)
	go func() {
		for {
			ch <- data
			time.Sleep(time.Millisecond * 100)
		}
	}()
	return ch
}

func FanIn() {
	ch1 := Generate("1")
	ch2 := Generate("2")
	ch3 := Generate("3")

	fanin := make(chan string)

	go func() {
		for {
			select {
			case str1 := <-ch1:
				fanin <- str1
			case str2 := <-ch2:
				fanin <- str2
			case str3 := <-ch3:
				fanin <- str3
			}
		}
	}()

	go func() {
		for v := range fanin {
			fmt.Println(v)
		}
	}()
}
