package main

import (
	"fmt"
	"sync"
	"time"
)

type CountdownLatch struct {
	wg *sync.WaitGroup
}

func NewCountDownLatch(max int) *CountdownLatch {
	var wgRef sync.WaitGroup
	wgRef.Add(max)
	return &CountdownLatch{wg: &wgRef}
}

func (latch *CountdownLatch) CountDown() {
	fmt.Println("Count down...")
	latch.wg.Done()
}

func (latch *CountdownLatch) Await() {
	fmt.Println("Waiting for latch to open")
	latch.wg.Wait()
}

func main() {
	latch := NewCountDownLatch(3)
	go doWork(latch)
	go doWork(latch)
	go doWork(latch)
	latch.Await()
	fmt.Println("Done")
}

func doWork(cl *CountdownLatch) {
	time.Sleep(5 * time.Second)
	cl.CountDown()
}
