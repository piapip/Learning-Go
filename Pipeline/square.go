package main

import (
	"fmt"
	"sync"
)

func gen(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, value := range nums {
			out <- value
		}
		close(out)
	}()
	return out
}

func sq(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for value := range in {
			out <- value * value
		}
		close(out)
	}()
	return out
}

func merge(done <-chan struct{}, fraction ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	output := func(c <-chan int) {
		for n := range c {
			select {
			case out <- n:
			//This two crapshoots here will trigger this "select" case
			// done <- struct{}{}
			// done <- struct{}{}
			case <-done:
			}
		}
		wg.Done()
	}
	wg.Add(len(fraction))
	for _, channel := range fraction {
		go output(channel)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func main() {
	fraction1 := sq(gen(2, 3))
	fraction2 := sq(gen(4, 5))

	done := make(chan struct{}, 2)
	out := merge(done, fraction1, fraction2)
	fmt.Println(<-out)

	done <- struct{}{}
	done <- struct{}{}
	// for value := range merge(fraction1, fraction2) {
	// 	fmt.Println(value)
	// }
}
