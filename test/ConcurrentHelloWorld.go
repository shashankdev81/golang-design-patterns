package main

import (
	"fmt"
)

func main() {
	chan1 := make(chan string)
	chan2 := make(chan string)

	RoutineAnnotator(func() {
		say("world", chan1)
	}, []string{"string", "world"})
	RoutineAnnotator(func() {
		say("hello", chan2)
	}, []string{"string", "world"})

	res1 := <-chan1
	res2 := <-chan2

	fmt.Printf("Chan1 %v, chan2 %v", res1, res2)
}

func say(s string, chan1 chan string) {
	fmt.Println("Inside say")
	chan1 <- s
}
