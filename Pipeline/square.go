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

func merge(fraction ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	output := func(c <-chan int) {
		for n := range c {
			out <- n
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

	for value := range merge(fraction1, fraction2) {
		fmt.Println(value)
	}
}
