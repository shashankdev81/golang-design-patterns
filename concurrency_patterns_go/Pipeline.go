package main

import (
	"fmt"
)

var multiplier = 5
var additive = 10

//type intFunc func(a int, b int) int

func main() {
	multiply := func(values []int, factor int) []int {
		multipliedValues := make([]int, len(values))

		for i, v := range values {
			multipliedValues[i] = v * factor
		}

		return multipliedValues
	}

	addition := func(values []int, factor int) []int {
		addedValues := make([]int, len(values))

		for i, v := range values {
			addedValues[i] = v + factor
		}

		return addedValues
	}

	divide := func(values []int, denom int) []int {
		dividedValues := make([]int, len(values))

		for i, v := range values {
			dividedValues[i] = v / denom
		}
		return dividedValues
	}

	/*
		A simple batch processing pipeline consisting of functions that have same signature and same i/p and o/p
		Pros: good expressive power
		Cons: large memory footprint, high verbosity

	*/
	result1 := divide(multiply(addition([]int{2, 4, 6}, 3), 5), 4)
	fmt.Println("Batch processing pipeline:", result1)

	/*
		Functional programming... functions as params
		Pros: excellent expressive power and re-use
		Cons: large memory footprint

	*/

	mul := func(x int, y int) int { return x * y }
	add := func(x int, y int) int { return x + y }
	div := func(x int, y int) int { return x / y }

	apply := func(values []int, factor int, intFunc func(a int, b int) int) []int {
		dividedValues := make([]int, len(values))

		for i, v := range values {
			dividedValues[i] = intFunc(v, factor)
		}
		return dividedValues

	}

	result2 := apply(apply(apply([]int{2, 4, 6}, 3, add), 5, mul), 4, div)
	fmt.Println("Batch processing pipeline functional way:", result2)
	/*
		A simple stream processing pipeline consisting of functions that have same signature and same i/p and o/p
		Pros: less memory footprint
		Cons: less expressive power, n function calls
	*/
	streamFunc := func(values []int, addFactor int, mulFactor int, divFactor int) []int {
		dividedValues := make([]int, len(values))
		for i, v := range values {
			dividedValues[i] = div(mul(add(v, addFactor), mulFactor), divFactor)
		}
		return dividedValues
	}
	fmt.Println("Stream processing pipeline functional way:", streamFunc([]int{2, 4, 6}, 3, 5, 4))

	/*
		A channel led concurrent design
		Pros Concurrency baked into design at all stages, stream based pipeline, less memory footprint, use range to exit when channel closes

		There is a recurrence-relation at play here (In mathematics, a recurrence relation is an equation that expresses
		the nth term of a sequence as a function of the k preceding terms, for some fixed k ...)
		At the beginning of the pipeline, we’ve established that we must convert
		discrete values into a channel. There are two points in this process which must be preemptable:
		1. Creation of the discrete value that is not nearly instantaneous.
		2. Sending of the discrete value on its channel.
		 The second is handled via our select statement and done channel which ensures that generator is preemptable even
		if it is blocked attempting to write to intStream.

	*/
	generator := func(done <-chan interface{}, integers ...int) <-chan int {
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

	divStream := func(done <-chan interface{}, intStream <-chan int, divFactor int) <-chan int {
		resultStream := make(chan int)
		go func() {
			defer close(resultStream)
			for i := range intStream {
				select {
				case <-done:
					return
				case resultStream <- i / divFactor:
				}
			}
		}()
		return resultStream
	}

	done := make(chan interface{})
	defer close(done)
	intStream := generator(done, 2, 4, 6)
	pipeline := divStream(done, multiplyStream(done, addStream(done, intStream, 3), 5), 4)
	fmt.Printf("Stream processing pipeline concurrent design:")
	fmt.Printf("[")
	for v := range pipeline {
		fmt.Printf("%v ", v)
	}
	fmt.Printf("]")

}

//multiply := func (values []int, multiplier int) []int {
//	multipliedValues := make([]int, len(values))
//	for i, v := range values {
//		multipliedValues[i] = v * multiplier
//	}
//	return multipliedValues
//}
//
//add := func (values []int, additive int) []int {
//	addedValues := make([]int, len(values))
//	for i, v := range values {
//		addedValues[i] = v + additive
//	}
//	return addedValues
//}
