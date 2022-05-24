package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	done := make(chan interface{})
	defer close(done)

	newSlice := func(start int, end int) []int {
		s := make([]int, end-start)
		count := 0
		for i := start; i < end; i++ {
			s[count] = i
			count++
		}
		return s
	}
	generator := func(done <-chan interface{}, integers []int) <-chan int {
		intStream := make(chan int)
		go func() {
			defer close(intStream)
			for _, v := range integers {
				select {
				case <-done:
					return
				case intStream <- v:
				}
			}
		}()
		return intStream
	}

	/* forks data from the input stream and passes on to a newly created channel*/
	fork := func(done <-chan interface{}, values <-chan int) chan int {
		valueStream := make(chan int)
		go func() {
			defer close(valueStream)
			for v := range values {
				select {
				case <-done:
					return
				case valueStream <- v:
				}
			}
		}()
		return valueStream
	}

	multiplyStream := func(done <-chan interface{}, intStream <-chan int, mulFactor int) <-chan int {
		resultStream := make(chan int)
		go func() {
			defer close(resultStream)
			for i := range intStream {
				select {
				case <-done:
					return
				case resultStream <- i * mulFactor:
				}
			}
		}()
		return resultStream
	}
	addStream := func(done <-chan interface{}, intStream <-chan int, addFactor int) <-chan int {
		resultStream := make(chan int)
		go func() {
			defer close(resultStream)
			for i := range intStream {
				select {
				case <-done:
					return
				case resultStream <- i + addFactor:
				}
			}
		}()
		return resultStream
	}

	concurrentSolution(generator, done, newSlice, addStream, multiplyStream, fork)
	//simpleSolution(generator, done, newSlice, addStream, multiplyStream, fork)
}

func concurrentSolution(generator func(done <-chan interface{}, integers []int) <-chan int, done chan interface{}, newSlice func(start int, end int) []int, addStream func(done <-chan interface{}, intStream <-chan int, addFactor int) <-chan int, multiplyStream func(done <-chan interface{}, intStream <-chan int, mulFactor int) <-chan int, fork func(done <-chan interface{}, values <-chan int) chan int) {
	cpus := runtime.NumCPU()
	fanOutChans := make([]<-chan int, cpus)
	input := generator(done, newSlice(1, 1000))
	for i := 0; i < cpus; i++ {
		fanOutChans[i] = addStream(done, multiplyStream(done, fork(done, input), 2), 2)
	}
	//time.Sleep(100 * time.Millisecond)
	//joinedChan := join(done, fanOutChans)

	var wg sync.WaitGroup
	multiplexedStream := make(chan int)
	wg.Add(len(fanOutChans))
	multiplex := func(c <-chan int) {
		defer wg.Done()
		for i := range c {
			//basically do some work and then write to a common channel
			time.Sleep(1 * time.Millisecond)
			select {
			case <-done:
				return
			case multiplexedStream <- i:
			}
		}
	}
	// Select from all the channels
	for _, ch := range fanOutChans {
		go multiplex(ch)
	}

	// Wait for all the reads to complete
	/* following code will deadlock since we will be waiting on lock but the subsequent code
	that's supposed to read wont be invoked ever
	*/
	//wg.Wait()
	//close(multiplexedStream)

	go func() {
		//we want to wait for any channels to finish so that we can go ahead and close the multiplexedStream
		wg.Wait()
		close(multiplexedStream)
	}()

	fmt.Println("")
	fmt.Println("Printing results of fork join")
	total := 0
	for j := range multiplexedStream {
		total = total + j
	}
	fmt.Printf("%v ", total)
}

func simpleSolution(generator func(done <-chan interface{}, integers []int) <-chan int, done chan interface{}, newSlice func(start int, end int) []int, addStream func(done <-chan interface{}, intStream <-chan int, addFactor int) <-chan int, multiplyStream func(done <-chan interface{}, intStream <-chan int, mulFactor int) <-chan int, fork func(done <-chan interface{}, values <-chan int) chan int) {
	input := generator(done, newSlice(1, 1000))
	output := addStream(done, multiplyStream(done, input, 2), 2)

	fmt.Println("")
	fmt.Println("Printing results of fork join")
	total := 0
	for j := range output {
		total = total + j
	}
	fmt.Printf("%v ", total)
}
