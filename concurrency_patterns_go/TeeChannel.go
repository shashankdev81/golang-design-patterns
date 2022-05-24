package main

import "fmt"

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

	//fork := func(done <-chan interface{}, values <-chan int) chan int {
	//	valueStream := make(chan int)
	//	go func() {
	//		defer close(valueStream)
	//		for v := range values {
	//			select {
	//			case <-done:
	//				return
	//			case valueStream <- v:
	//			}
	//		}
	//	}()
	//	return valueStream
	//}

	//orDone := func(done <-chan interface{}, c <-chan int) <-chan int {
	//	valStream := make(chan int)
	//	go func() {
	//		defer close(valStream)
	//		for {
	//			select {
	//			case <-done:
	//				return
	//			case v, ok := <-c:
	//				if ok == false {
	//					return
	//				}
	//				select {
	//				case valStream <- v:
	//				case <-done:
	//				}
	//			}
	//		}
	//	}()
	//	return valStream
	//}

	out1, out2 := tee(done, generator(done, newSlice(1, 10)))

	printChan(out1)
	printChan(out2)
}

func tee(done <-chan interface{}, values <-chan int) (_, _ <-chan int) {

	out1 := make(chan int)
	out2 := make(chan int)
	go func() {
		defer close(out1)
		defer close(out2)
		for val := range values {
			var out1, out2 = out1, out2
			for i := 0; i < 2; i++ {
				select {
				case <-done:
					return
				case out1 <- val:
					out1 = nil
				case out2 <- val:
					out2 = nil
				}
			}
		}
	}()

	return out1, out2
}

func printChan(out <-chan int) {
	fmt.Println("")
	fmt.Println("Printing channel:")
	for j := range out {
		fmt.Printf("%v ", j)
	}
}
