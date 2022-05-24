package main

import (
	"fmt"
	"time"
)

var ops int

func main() {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}
	start := time.Now()

	c1 := sig(1 * time.Hour)
	c2 := sig(2 * time.Minute)
	c3 := sig(3 * time.Hour)
	c4 := sig(4 * time.Minute)
	c5 := sig(2 * time.Second)
	c6 := sig(6 * time.Hour)
	c7 := sig(7 * time.Minute)
	c8 := sig(8 * time.Minute)
	c9 := sig(9 * time.Hour)
	c10 := sig(10 * time.Minute)

	<-or(c1, c2, c3, c4, c5, c6, c7, c8, c9, c10)
	checkIfChannelClosed(c1)
	checkIfChannelClosed(c2)
	checkIfChannelClosed(c3)
	checkIfChannelClosed(c4)
	checkIfChannelClosed(c5)
	checkIfChannelClosed(c6)
	checkIfChannelClosed(c7)
	checkIfChannelClosed(c8)
	checkIfChannelClosed(c9)
	checkIfChannelClosed(c10)
	fmt.Printf("done after %v\n", time.Since(start))

}

func checkIfChannelClosed(c1 <-chan interface{}) {
	ok := true
	select {
	case _, ok = <-c1:
	default:
	}
	if ok {
		fmt.Printf("Channel is not closed\n")
	} else {
		fmt.Printf("Channel is closed\n")
	}
}

func or(channels ...<-chan interface{}) <-chan interface{} {
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}
	orDone := make(chan interface{})
	ops = ops + 1
	go func(count int) {
		fmt.Printf("Starting Another go routine %v \n", count)
		defer fmt.Printf("Closing another go routine %v\n", count)
		defer close(orDone)
		switch len(channels) {
		case 2:
			select {
			case <-channels[0]:
			case <-channels[1]:

			}
		default:
			select {
			case <-channels[0]:

			case <-channels[1]:

			case <-channels[2]:

			case <-or(append(channels[3:], orDone)...):

			}
		}
	}(ops)
	return orDone
}
