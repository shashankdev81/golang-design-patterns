package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	var wg sync.WaitGroup
	var concurrency int64
	var MAX int64
	MAX = 10
	work := func() {
		fmt.Println("Don't chase growth; instead make yourself useful")
		time.Sleep(3 * time.Second)
	}
	for i := 0; i < 30; i++ {
		executeConcurrently(&concurrency, MAX, work, &wg)
	}
	wg.Wait()

}

func executeConcurrently(concurrency *int64, MAX int64, work func(), wg *sync.WaitGroup) {
	wg.Add(1)
	for {
		if *concurrency < MAX && atomic.CompareAndSwapInt64(concurrency, *concurrency, *concurrency+1) {
			go func(f func()) {
				defer wg.Done()
				f()
				atomic.AddInt64(concurrency, -1)
			}(work)
			break
		}
	}
}
