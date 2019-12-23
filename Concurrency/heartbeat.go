package main

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func testGenInt() {
	doWork := func(done <-chan interface{}, nums ...int) (<-chan interface{}, <-chan int) {
		heartbeatStream := make(chan interface{}, 1)
		intStream := make(chan int)
		go func() {
			defer close(heartbeatStream)
			defer close(intStream)

			//could put some delay here if the test takes time to run.
			time.Sleep(2 * time.Second)

			sendPulse := func() {
				select {
				case heartbeatStream <- struct{}{}:
				default:
				}
			}

			sendResult := func(n int) {
				select {
				case <-done:
					return
				case intStream <- n:
				}
			}

			for _, n := range nums {
				sendPulse()
				sendResult(n)
			}
		}()

		return heartbeatStream, intStream
	}

	//Without heartbeat, can't control the flow, it becomes nondeterministic
	theBadTest = func(t *testing.T) {
		done := make(chan interface{})
		defer close(done)

		intSlice := []int{0, 1, 2, 3, 5}
		_, results := doWork(done, intSlice...)
		for i, expected := range intSlice {
			select {
			case r := <-results:
				if r != expected {
					t.Errorf(
						"index %v: expected %v, but received %v,",
						i,
						expected,
						r,
					)
				}
			case <-time.After(1 * time.Second):
				t.Fatal("test timed out")
			}
		}
	}

	//use heartbeat, 1 heartbeat means 1 result, flow becomes easy to control
	theGoodTest = func(t *testing.T) {
		done := make(chan interface{})
		defer close(done)
		intSlice := []int{0, 1, 2, 3, 5}
		// _, results := doWork(done, intSlice...)
		heartbeat, results := doWork(done, intSlice...)

		<-heartbeat //remember even though it only has 1 slot, so it's blocking the flow of the doWork.
		//so the doWork is not actually finished yet, it's pending for the pump signal
		i := 0
		for r := range results {
			if expected := intSlice[i]; r != expected {
				t.Errorf("index %v: expected %v, but received %v,", i, expected,
					r)
			}
			i++
		}
	}
}

//this thing sends out results instantly... Quite weird for a goroutine stuff. Good for testing they say.
func testWeirdPulse() {
	doWork := func(done <-chan interface{}) (<-chan interface{}, <-chan int) {
		heartbeatStream := make(chan interface{}, 1)
		workStream := make(chan int)
		go func() {
			defer close(heartbeatStream)
			defer close(workStream)

			//could put some delay here if the test takes time to run.

			sendPulse := func() {
				select {
				case heartbeatStream <- struct{}{}:
				default:
				}
			}

			sendResult := func() {
				select {
				case <-done:
					return
				case workStream <- rand.Intn(10):
				}
			}

			for i := 0; i < 10; i++ {
				sendPulse()
				sendResult()
			}
		}()

		return heartbeatStream, workStream
	}
	done := make(chan interface{})
	defer close(done)
	heartbeat, results := doWork(done)
	for {
		select {
		case _, ok := <-heartbeat:
			if ok {
				fmt.Println("pulse")
			} else {
				return
			}
		case r, ok := <-results:
			if ok {
				fmt.Printf("results %v\n", r)
			} else {
				return
			}
		}
	}
}

func testBadHeartbeat() {
	doWork := func(done <-chan interface{}, pulseInterval time.Duration) (<-chan interface{}, <-chan time.Time) {
		heartbeat := make(chan interface{})
		results := make(chan time.Time)
		go func() {
			// defer close(heartbeat)
			// defer close(results)

			pulse := time.Tick(pulseInterval)
			workGen := time.Tick(2 * pulseInterval)

			sendPulse := func() {
				select {
				case heartbeat <- struct{}{}:
				default:
				}
			}

			sendResult := func(r time.Time) {
				for {
					select {
					case <-pulse:
						sendPulse()
					case results <- r:
						return
					}
				}
			}

			// for{
			for i := 0; i < 2; i++ {
				select {
				case <-done:
					return
				case <-pulse:
					sendPulse()
				case r := <-workGen:
					sendResult(r)
				}
			}
		}()
		return heartbeat, results
	}
	done := make(chan interface{})
	time.AfterFunc(10*time.Second, func() { close(done) })
	const timeout = 2 * time.Second
	heartbeat, results := doWork(done, timeout/2)
	for {
		select {
		case _, ok := <-heartbeat:
			if ok == false {
				return
			}
			fmt.Println("Pulse")
		case r, ok := <-results:
			if ok == false {
				return
			}
			fmt.Println("Results ", r.Second())
		case <-time.After(timeout):
			fmt.Println("worker goroutine is not healthy!")
			return
		}
	}
}

func testGoodHeartbeat() {
	doWork := func(done <-chan interface{}, pulseInterval time.Duration) (<-chan interface{}, <-chan time.Time) {
		heartbeat := make(chan interface{})
		results := make(chan time.Time)
		go func() {
			defer close(heartbeat)
			defer close(results)

			pulse := time.Tick(pulseInterval)
			workGen := time.Tick(2 * pulseInterval)

			sendPulse := func() {
				select {
				case heartbeat <- struct{}{}:
				default:
				}
			}

			sendResult := func(r time.Time) {
				for {
					select {
					case <-done:
						return
					case <-pulse:
						sendPulse()
					case results <- r:
						return
					}
				}
			}

			for {
				select {
				case <-done:
					return
				case <-pulse:
					sendPulse()
				case r := <-workGen:
					sendResult(r)
				}
			}
		}()
		return heartbeat, results
	}
	done := make(chan interface{})
	time.AfterFunc(11*time.Second, func() {
		close(done)
	})
	const timeout = 2 * time.Second
	heartbeat, results := doWork(done, timeout/2)
	for {
		select {
		case _, ok := <-heartbeat:
			if ok == false {
				return
			}
			fmt.Println("Pulse")
		case r, ok := <-results:
			if ok == false {
				return
			}
			fmt.Println("Results ", r.Second())
		case <-time.After(timeout):
			fmt.Println("If this shows up then there's a goroutine problem")
			return
		}
	}
}
