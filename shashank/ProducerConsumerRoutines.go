package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int)
	closed := make(chan bool)
	go func(closed chan bool) {
		i := 0
	loop:
		for {
			select {
			case <-closed:
				close(ch)
				break loop
			default:
				fmt.Println("Produce", i)
				ch <- i
				i++
			}
		}
		fmt.Println("Producer done")

	}(closed)

	go func() {
		for r := range ch {
			fmt.Println("Consume", r)

		}
	}()
	fmt.Println("Will sleep")
	time.Sleep(1 * time.Millisecond)
	fmt.Println("Awake")
	close(closed)

}
