package cond

import (
	"fmt"
	"sync"
	"time"
)

func Cond() {
	var wg sync.WaitGroup
	var condition = false
	wg.Add(2)
	c := sync.NewCond(&sync.Mutex{})
	go func() {
		defer wg.Done()
		c.L.Lock()
		if !condition {
			fmt.Println("goroutine 1 wait")
			c.Wait()
		}
		fmt.Println("goroutine 1 continue execution")
		c.L.Unlock()
	}()
	go func() {
		defer wg.Done()
		c.L.Lock()
		for !condition {
			fmt.Println("goroutine 2 wait")
			c.Wait()
		}
		fmt.Println("goroutine 2 continue execution")
		c.L.Unlock()
	}()

	time.Sleep(time.Second)
	c.L.Lock()
	fmt.Println("main goroutine ready")
	condition = true
	c.Broadcast()
	fmt.Println("main goroutine broadcast")
	c.L.Unlock()
	wg.Wait()
}
