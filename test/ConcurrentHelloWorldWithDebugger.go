package main

import (
	"fmt"
	"github.com/dlsniper/debugger"
)

func main() {
	chan1 := make(chan string)
	chan2 := make(chan string)

	go say1("world", chan1)
	go say1("hello", chan2)

	res1 := <-chan1
	res2 := <-chan2

	fmt.Printf("Chan1 %v, chan2 %v", res1, res2)
}

func say1(s string, chan1 chan string) {
	fmt.Println("Inside say")
	debugger.SetLabels(func() []string {
		return []string{
			"displayString", s,
		}
	})
	chan1 <- s
}
