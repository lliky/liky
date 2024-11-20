package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var ch chan string

func init() {
	ch = make(chan string, 10)
}

type aaa struct {
	aMap map[string]struct{}
	lock sync.Mutex
}

func main() {
	a := aaa{
		aMap: map[string]struct{}{},
	}

	go func() {
		for k := range ch {
			a.lock.Lock()
			a.aMap[k] = struct{}{}
			a.lock.Unlock()
			fmt.Println("add: ", k)
		}
	}()

	go func() {
		for {

			a.lock.Lock()
			if len(a.aMap) > 0 {
				for k, _ := range a.aMap {
					a.lock.Unlock()
					// select {
					// case <-time.After(time.Second):
					// 	fmt.Println("skip: ")
					// default:
					a.lock.Lock()
					delete(a.aMap, k)
					a.lock.Unlock()
					fmt.Println("del: ", k)

					// }
					break
				}
			} else {
				a.lock.Unlock()
			}
			// time.Sleep(time.Millisecond)
		}
	}()
	go func() {
		var i int
		for {
			ip := fmt.Sprintf("a: %d", time.Now().Local().UnixNano())
			ch <- ip
			//time.Sleep(time.Second)
			time.Sleep(time.Microsecond * time.Duration(rand.Int63n(64)))
			if i == 1000 {
				time.Sleep(time.Minute)
			} else {
				i++
			}
		}
	}()
	select {}
}
