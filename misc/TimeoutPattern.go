package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	//d := time.Now().Add(50 * time.Millisecond)
	//based on absolute time
	//ctx, cancel := context.WithDeadline(context.Background(), d)
	//based on relative time
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	func1(ctx)
}

func func1(ctx context.Context) {
	func2(ctx)
}

func func2(ctx context.Context) {
	func3(ctx)
}

func func3(ctx context.Context) {
	for i := 0; ; i++ {
		select {
		case <-time.After(1 * time.Second):
			fmt.Println("overslept")
			return
		case <-ctx.Done():
			fmt.Println(ctx.Err())
			return
		default:
			fmt.Println("Current value of i=", i)
		}
	}
}
