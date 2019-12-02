package main

import (
	"fmt"
	"time"
)

func preventGoroutineLeak() {

	doWork := func(done <-chan interface{}, strings <-chan string) <-chan interface{} {
		terminated := make(chan interface{})
		go func() {
			defer fmt.Println("doWork exited.")
			defer close(terminated)
			for {
				select {
				case s := <-strings:
					terminated <- s + " 1"
				case <-done:
					fmt.Println("Got here")
					return
				}
			}
		}()
		return terminated
	}

	done := make(chan interface{})

	strings := make(chan string)
	go func() {
		strings <- "Hello"
		strings <- "world"
	}()

	terminated := doWork(done, strings)

	// terminated := doWork(done, nil)
	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("Canceling doWork goroutine...")
		close(done)
	}()

	<-terminated
	fmt.Println("Done")

	// doWork := func(lock *sync.Mutex, strings <-chan string) <-chan interface{} {
	// 	complete := make(chan interface{})
	// 	go func() {
	// 		lock.Lock()
	// 		defer lock.Unlock()
	// 		defer fmt.Println("doWork exited.")
	// 		defer close(complete)
	// 		for s := range strings {
	// 			fmt.Println(s)
	// 			s = s + " 1"
	// 		}
	// 	}()

	// 	return complete
	// }

	// strings := make(chan string)
	// var wg sync.WaitGroup
	// var lock sync.Mutex
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	defer close(strings)
	// 	strings <- "Hello"
	// 	strings <- "world"
	// }()

	// channel := doWork(&lock, strings)
	// wg.Wait()

	// fmt.Println(<-channel)
	// fmt.Println("done")

}
