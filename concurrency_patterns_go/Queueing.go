package main

import (
	"fmt"
	"time"
)

//the utility of introducing a queue, isn’t that the run time of
//one of stages has been reduced, but rather that the time it’s in a blocking state is reduced.
func main() {
	done := make(chan interface{})
	defer close(done)

	//connStream := make(chan int)
	connStream := make(chan int, 100)
	defer close(connStream)

	httpHandler := func(done <-chan interface{}, connections <-chan int) <-chan int {
		go func() {
			for c := range connections {
				select {
				case <-done:
					return
				case connStream <- c:
					{
						time.Sleep(1 * time.Second)
					}
				}
			}
		}()
		return connStream
	}

	//connStream2 := make(chan int)
	connStream2 := make(chan int, 100)
	defer close(connStream2)

	acceptConnection := func(done chan interface{}, handler <-chan int) chan int {
		go func() {
			for c := range handler {
				select {
				case <-done:
					return
				case connStream2 <- c:
					{
						time.Sleep(2 * time.Second)

					}
				}
			}
		}()
		return connStream
	}

	//results := make(chan int)
	results := make(chan int, 100)
	defer close(results)
	processRequest := func(done chan interface{}, connections <-chan int) chan int {
		go func() {
			for c := range connections {
				select {
				case <-done:
					return
				case results <- c:
					{
						time.Sleep(1 * time.Second)
					}
				}
			}
		}()
		return results
	}

	//inStream := make(chan int)
	inStream := make(chan int, 1000)
	defer close(inStream)
	outcomes := processRequest(done, acceptConnection(done, httpHandler(done, inStream)))
	go func() {
		for res := range outcomes {
			fmt.Println("Connection processed successfully", res)
		}
	}()
	conn := 0
	for start := time.Now(); time.Since(start) < 10*time.Second; {
		conn++
		fmt.Println("Accepted  connections", conn)
		inStream <- conn
	}
	time.Sleep(10 * time.Second)

}
