package main

import (
	"context"
	"fmt"
	"time"
)

func oldTest() {
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
			defer close(out)
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
	//	|  ctx1 (end when canceled)  |
	//	|  cancel1                   |
	//	|                            |
	//	==============================
	ctx1, cancel1 := context.WithCancel(context.Background())
	defer cancel1()

	for value := range gen(ctx1) {
		fmt.Println(value)
		if value == 5 {
			break
		}
	}

	//	==============================
	//	|  ctx1 (end when canceled)  |
	//	|  cancel1                   |
	//	|  ctx2 (end after 50ms)     |
	//	|  cancel2                   |
	//	==============================

	d := time.Now().Add(50 * time.Millisecond)
	ctx2, cancel2 := context.WithDeadline(context.Background(), d)
	defer cancel2()

	select {
	case <-time.After(1 * time.Second):
		fmt.Println("Overslept")
	case <-ctx2.Done():
		fmt.Println(ctx2.Err())
	}

	//	====================================================
	//	|  ctx1 (end when canceled)                        |
	//	|  cancel1                                         |
	//	|  ctx2 (end after 50ms)                           |
	//	|  cancel2                                         |
	//	|  ctx3 (a copy of ctx, store hash value there)    |
	//	====================================================

	type favContextKey string

	f := func(ctx context.Context, k favContextKey) {
		if v := ctx.Value(k); v != nil {
			fmt.Println("Found value: ", v)
			return
		}
		fmt.Println("Key not found: ", k)
	}

	k := favContextKey("Language")
	ctx := context.WithValue(context.Background(), k, "Go")
	f(ctx, k)
	f(ctx, favContextKey("Color"))

}
