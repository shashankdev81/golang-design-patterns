package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type CyclicBarrier struct {
	maxParticipants  int
	currParticipants int32
	barrierChan      chan int
	orDone           chan interface{}
	once             sync.Once
}

func (barrier *CyclicBarrier) await() {
	fmt.Println("Waiting on barrier")
	var curr int32
	curr = atomic.AddInt32(&barrier.currParticipants, 1)
	if int(curr) > barrier.maxParticipants {
		panic("This barrier can only handle these many participants: " + strconv.Itoa(barrier.maxParticipants))
	}
	fmt.Println("Will send value on channel", barrier.currParticipants)
	barrier.barrierChan <- int(curr)
	for {
		select {
		case <-barrier.orDone:
			fmt.Println("Barrier broken")
			return
		}
	}
}
func NewBarrier(max int) *CyclicBarrier {
	if max < 1 {
		panic("Please provide a positive value for max participants")
	}
	barrier := &CyclicBarrier{maxParticipants: max, barrierChan: make(chan int, max), orDone: make(chan interface{})}
	go func() {
		for {
			select {
			case curr := <-barrier.barrierChan:
				fmt.Println("Participant entering barrier", curr)
				if curr == barrier.maxParticipants {
					fmt.Println("Will close barrier")
					//once not needed
					barrier.once.Do(func() { close(barrier.orDone) })
				}
			default:
				//fmt.Println("Nothing happened")
			}
		}
	}()
	return barrier
}

func main() {
	//barrier := NewBarrier(-1)
	barrier := NewBarrier(3)
	go workConcurrently(barrier)
	go workConcurrently(barrier)
	go workConcurrently(barrier)
	//go workConcurrently(barrier)
	time.Sleep(250 * time.Millisecond)

}

func workConcurrently(barrier *CyclicBarrier) {
	fmt.Println("Invoking a go routine and doing some work")
	r := rand.Intn(20)
	time.Sleep(time.Duration(r) * time.Millisecond)
	barrier.await()
}
