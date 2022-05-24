package snippets

import (
	"fmt"
	"time"
)

func main() {
	/*Useful generators*/
	repeat := func(done <-chan interface{}, values ...interface{}) chan interface{} {
		valueStream := make(chan interface{})
		go func() {
			defer close(valueStream)
			for {
				for _, v := range values {
					select {
					case <-done:
						return
					case valueStream <- v:
					}
				}
			}
		}()
		return valueStream
	}

	take := func(done <-chan interface{}, valueStream <-chan interface{}, num int) chan interface{} {
		takeStream := make(chan interface{})
		go func() {
			defer close(takeStream)
			for i := 0; i < num; i++ {
				select {
				case <-done:
					return
				case takeStream <- <-valueStream:
				}
			}
		}()
		return takeStream
	}

	done := make(chan interface{})
	defer close(done)
	fmt.Println("")
	fmt.Println("Will print 1 many times")
	for val := range take(done, repeat(done, 1), 10) {
		fmt.Printf("%v ", val)
	}
	fmt.Println("")

	/*Simple util funcs */
	toString := func(done <-chan interface{}, inputStream chan interface{}) chan interface{} {
		valueStream := make(chan interface{})
		go func() {
			defer close(valueStream)
			for v := range inputStream {
				select {
				case <-done:
					return
				case valueStream <- fmt.Sprint(v): //TODO ? does golang provide an overridable toStroing method for interface{} type
				}
			}
		}()
		return valueStream
	}
	transform := func(done <-chan interface{}, f func(val interface{}) interface{}, values ...interface{}) chan interface{} {
		valueStream := make(chan interface{})
		go func() {
			defer close(valueStream)
			for {
				for _, v := range values {
					select {
					case <-done:
						return
					case valueStream <- f(v):
					}
				}
			}
		}()
		return valueStream
	}
	done2 := make(chan interface{})
	defer close(done2)
	toStr := func(input interface{}) interface{} {
		return input
	}
	fmt.Println("Will convert slice to channel")
	go func() {
		for v := range toString(done2, transform(done2, toStr, 2, 4, 6)) {
			fmt.Printf("%v ", v)
		}
	}()
	time.Sleep(1 * time.Millisecond)
}
