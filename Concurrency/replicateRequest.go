package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// func test() {
// 	doWork := func(done <-chan interface{}, id int, wg *sync.WaitGroup, result chan<- int) {
// 		started := time.Now()
// 		defer wg.Done()
// 		// Simulate random load
// 		simulatedLoadTime := time.Duration(1+rand.Intn(5)) * time.Second
// 		select {
// 		case <-done:
// 		case <-time.After(simulatedLoadTime):
// 		}
// 		select {
// 		case <-done:
// 		case result <- id:
// 		}
// 		took := time.Since(started)
// 		// Display how long handlers would have taken
// 		if took < simulatedLoadTime {
// 			took = simulatedLoadTime
// 		}
// 		fmt.Printf("%v took %v\n", id, took)
// 	}
// 	done := make(chan interface{})
// 	result := make(chan int)
// 	var wg sync.WaitGroup
// 	wg.Add(10)
// 	for i := 0; i < 10; i++ {
// 		go doWork(done, i, &wg, result)
// 	}
// 	firstReturned := <-result //this line is the pattern's special spot.
// 	// the reason being for this little shit to be here ís because of the time.After up there. If I let the whole thing to run and call
// 	close(done)
// 	wg.Wait()
// 	// firstReturned := <-result
// 	fmt.Printf("Received an answer from #%v\n", firstReturned)
// }

func test() {
	doWork := func(done <-chan interface{}, id int, wg *sync.WaitGroup, result chan<- int) {
		defer wg.Done()

		started := time.Now()
		simulatedLoadTime := time.Duration(1+rand.Intn(5)) * time.Second

		// select {
		// case <-done:
		// case <-time.After(simulatedLoadTime):
		// }
		// select {
		// case <-done:
		// case result <- id:
		// }

		select {
		case <-done:
		case <-time.After(simulatedLoadTime):
		case result <- id:
		}

		// Display how long handlers would have taken
		took := time.Since(started)
		if took < simulatedLoadTime {
			took = simulatedLoadTime
		}
		fmt.Printf("%v took %f\n", id, took.Seconds())
	}
	done := make(chan interface{})
	// result := make(chan int)
	result := make(chan int, 1)

	wg := &sync.WaitGroup{}
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go doWork(done, i, wg, result)
	}
	// firstReturned := <-result //this line is the pattern's special spot.
	time.Sleep(time.Second) //this line is a must if you want to run it this way.
	// so the time.After up there can be triggered and maintain the order.
	close(done)
	wg.Wait()

	firstReturned := <-result //this line is the pattern's special spot.
	fmt.Printf("Received an answer from #%v\n", firstReturned)
}
