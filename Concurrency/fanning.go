package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func testFanIn() {
	repeatFn := func(done <-chan interface{}, fn func() interface{}) <-chan interface{} {
		valueStream := make(chan interface{})
		go func() {
			defer close(valueStream)
			for {
				select {
				case <-done:
					return
				case valueStream <- fn():
				}
			}
		}()
		return valueStream
	}

	take := func(done <-chan interface{}, valueStream <-chan interface{}, num int) <-chan interface{} {
		takeStream := make(chan interface{})
		go func() {
			defer close(takeStream)
			for i := 0; i < num; i++ {
				select {
				case <-done:
					return
				case takeStream <- <-valueStream:
				}
			}
		}()
		return takeStream
	}

	toInt := func(done <-chan interface{}, valueStream <-chan interface{}) <-chan int {
		intStream := make(chan int)
		go func() {
			defer close(intStream)
			for v := range valueStream {
				select {
				case <-done:
					return
				case intStream <- v.(int):
				}
			}
		}()
		return intStream
	}

	primeFinder := func(done <-chan interface{}, intStream <-chan int) <-chan interface{} {
		primeStream := make(chan interface{})
		go func() {
			defer close(primeStream)
			for integer := range intStream {
				integer--
				prime := true
				for divisor := integer - 1; divisor > 1; divisor-- {
					if integer%divisor == 0 {
						prime = false
						break
					}
				}

				if prime {
					select {
					case <-done:
						return
					case primeStream <- integer:
					}
				}
			}
		}()
		return primeStream
	}

	done := make(chan interface{})
	defer close(done)
	start := time.Now()

	rd := func() interface{} {
		return rand.Intn(50000000)
	}
	randIntStream := toInt(done, repeatFn(done, rd))

	numFinders := runtime.NumCPU()
	fmt.Println("Spinning upto ", numFinders, " prime finders")
	finders := make([]<-chan interface{}, numFinders)
	for i := 0; i < numFinders; i++ {
		finders[i] = primeFinder(done, randIntStream)
	}
	fmt.Println("Primes: ")
	for prime := range take(done, fanIn(done, finders...), 10) {
		fmt.Printf("%v ", prime)
	}
	fmt.Printf("\nSearch took: %v\n", time.Since(start))
}

func testBadExample() {
	repeatFn := func(done <-chan interface{}, fn func() interface{}) <-chan interface{} {
		valueStream := make(chan interface{})
		go func() {
			defer close(valueStream)
			for {
				select {
				case <-done:
					return
				case valueStream <- fn():

				}
			}
		}()
		return valueStream
	}

	take := func(done <-chan interface{}, valueStream <-chan interface{}, num int) <-chan interface{} {
		takeStream := make(chan interface{})
		go func() {
			defer close(takeStream)
			for i := 0; i < num; i++ {
				select {
				case <-done:
					return
				case takeStream <- <-valueStream:
				}
			}
		}()
		return takeStream
	}

	toInt := func(done <-chan interface{}, valueStream <-chan interface{}) <-chan int {
		intStream := make(chan int)
		go func() {
			defer close(intStream)
			for v := range valueStream {
				select {
				case <-done:
					return
				case intStream <- v.(int):
				}
			}
		}()
		return intStream
	}

	//find the closest prime to the int. And why is it chan interface{}... and not chan int... a lil bit weird...
	//oooh, because take channel only accepts chan interface{}
	primeFinder := func(done <-chan interface{}, intStream <-chan int) <-chan interface{} {
		primeStream := make(chan interface{})
		go func() {
			defer close(primeStream)
			for integer := range intStream {
				prime := true
				for divisor := integer - 1; divisor > 1; divisor-- {
					if integer%divisor == 0 {
						prime = false
						break
					}
				}

				if prime {
					select {
					case <-done:
						return
					case primeStream <- integer:
					}
				}
			}
		}()
		return primeStream
	}

	rd := func() interface{} {
		return rand.Intn(50000000)
	}

	done := make(chan interface{})
	defer close(done)

	start := time.Now()
	randIntStream := toInt(done, repeatFn(done, rd))
	fmt.Println("Primes:")
	for prime := range take(done, primeFinder(done, randIntStream), 10) {
		fmt.Printf("%v ", prime)
	}
	fmt.Printf("\nSearch took: %v\n", time.Since(start))
}

func fanIn(done <-chan interface{}, channels ...<-chan interface{}) <-chan interface{} {
	var wg sync.WaitGroup
	multiplexedStream := make(chan interface{})

	muliplex := func(c <-chan interface{}) {
		defer wg.Done()
		for i := range c {
			select {
			case <-done:
				return
			case multiplexedStream <- i:
			}
		}
	}

	wg.Add(len(channels))
	for _, c := range channels {
		go muliplex(c)
	}

	go func() {
		wg.Wait()
		close(multiplexedStream)
	}()

	// wg.Wait()
	// close(multiplexedStream)

	return multiplexedStream
}
