package main

import (
	"fmt"
	"time"
)

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
	<-or2(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)
	fmt.Printf("done after %v\n", time.Since(start))

	//start = time.Now()
	//<-or(
	//	sig(2*time.Second),
	//	sig(3*time.Second),
	//)
	//fmt.Printf("done after %v\n", time.Since(start))
	//
	//start = time.Now()
	//<-or(
	//	sig(2 * time.Second),
	//)
	//fmt.Printf("done after %v\n", time.Since(start))
}

func or2(channels ...<-chan interface{}) <-chan interface{} {

	orDone := make(chan interface{})
	fmt.Println("making new done channel")
	go func() {
		defer close(orDone)
		switch len(channels) {
		case 0:
			panic("Can't combine non existent channels")
		case 1:
			<-channels[0]
		default:
			select {
			case <-channels[0]:
			case <-channels[1]:
			case <-or2(append(channels[2:], orDone)...):
			}
		}
	}()
	return orDone
}

func or1(channels ...<-chan interface{}) <-chan interface{} {
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}
	orDone := make(chan interface{})
	fmt.Println("making new done channel")
	go func() {
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
			case <-or1(append(channels[3:], orDone)...):
			}
		}
	}()
	return orDone
}
