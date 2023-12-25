package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	cancel, _ := context.WithCancel(context.Background())
	deadline, ok := cancel.Deadline()
	fmt.Printf("deadline: %v, ok: %v\n", deadline, ok)
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	go doSomething(ctx)
	select {
	case <-ctx.Done():
		fmt.Printf("oh, no, I have exceeded the deadline\n")
		fmt.Println(ctx.Err().Error())
	}
	var t time.Time
	fmt.Printf("t: %v", t)
	time.Sleep(time.Second)
}

func doSomething(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("time out\n")
			fmt.Println(ctx.Err().Error())
			return
		default:
			fmt.Printf("do something...\n")

		}
		time.Sleep(time.Second)
	}
}

func watch(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("%s goroutine err: %s\n", name, ctx.Err().Error())
			fmt.Printf("%s goroutine stopped...\n", name)
			return
		default:
			fmt.Printf("%s goruntine executing\n", name)
			time.Sleep(time.Second)
		}
	}
}
