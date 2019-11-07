package main

import (
	"context"
	"fmt"
)

func main() {
	// gen generates integers in a separate goroutine and
	// sends them to the returned channel.
	// The callers of gen need to cancel the context once
	// they are done consuming generated integers not to leak
	// the internal goroutine started by gen.

	//put a Goroutine into the passed context
	gen := func(ctx context.Context) <-chan int {
		out := make(chan int)
		n := 1
		go func() {
			// defer close(out) where is this line?
			for {
				select {
				case out <- n:
					n++
				case <-ctx.Done():
					return //signal to end all of this code's careers
				}
			}
		}()
		return out
	}

	//context.Background() here will be the background for the context.WithCancel

	//	==============================
	//	|  WithCancel                |
	//	|  cancel                    |
	//	|                            |
	//	|                            |
	//	|                            |
	//	|                            |
	//	==============================
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for value := range gen(ctx) {
		fmt.Println(value)
		if value == 5 {
			break
		}
	}
}
