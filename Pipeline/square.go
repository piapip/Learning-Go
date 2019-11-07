package main

import (
	"fmt"
	"sync"
)

func gen(done <-chan struct{}, nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, value := range nums {
			select {
			case out <- value:
			case <-done:
				return
			}
		}
	}()
	return out
}

func sq(done <-chan struct{}, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for value := range in {
			select {
			case out <- value * value:
			case <-done:
				return
			}
		}
	}()
	return out
}

func merge(done <-chan struct{}, fractions ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	output := func(fraction <-chan int) {
		defer wg.Done()
		for value := range fraction {
			select {
			case out <- value:
			case <-done:
				return
			}
		}
	}

	wg.Add(len(fractions))

	//this will give out 1 set of answer, all channels are packed into this one channel
	// go func() {
	// 	for _, channel := range fractions {
	// 		output(channel)
	// 	}
	// }()

	//but this will give out multiple set of answer, each channel will be handled on its own channel
	for _, channel := range fractions {
		go output(channel)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func main() {
	done := make(chan struct{})
	defer close(done)

	fraction1 := sq(done, gen(done, 1, 2))
	fraction2 := sq(done, gen(done, 3, 4))

	wholePiece := merge(done, fraction1, fraction2)

	for value := range wholePiece {
		fmt.Println(value)
	}
}
