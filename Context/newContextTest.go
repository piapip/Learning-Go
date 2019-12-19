package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func test() {
	handleResponse := func(ctx context.Context) {
		fmt.Printf("handling response for %v (%v)\n", ctx.Value("userID"), ctx.Value("authToken"))
	}

	processRequest := func(userID, authToken string) {
		ctx := context.WithValue(context.Background(), interface{}("userID"), userID) //using interface to suppresss the warning
		ctx = context.WithValue(ctx, interface{}("authToken"), authToken)
		handleResponse(ctx)
	}

	processRequest("jane", "abc123")
}

func testDeadline() {
	local := func(ctx context.Context) (string, error) {
		if deadline, ok := ctx.Deadline(); ok {
			if deadline.Sub(time.Now().Add(1*time.Minute)) <= 0 {
				return "", context.DeadlineExceeded
			}
		}
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-time.After(10 * time.Minute):
		}

		return "EN/US", nil
	}

	genGreeting := func(ctx context.Context) (string, error) {
		ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
		defer cancel()

		switch local, err := local(ctx); {
		case err != nil:
			return "", err
		case local == "EN/US":
			return "Hello", nil
		}
		return "", fmt.Errorf("Unsupported locale")
	}

	genFarewell := func(ctx context.Context) (string, error) {
		switch local, err := local(ctx); {
		case err != nil:
			return "", err
		case local == "EN/US":
			return "Goodbye", nil
		}
		return "", fmt.Errorf("Unsupported locale")
	}

	printGreeting := func(ctx context.Context) error {
		greeting, err := genGreeting(ctx)
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println(greeting, "world!")
		return nil
	}

	printFarewell := func(ctx context.Context) error {
		farewell, err := genFarewell(ctx)
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println(farewell)
		return nil
	}

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := printGreeting(ctx); err != nil {
			fmt.Printf("cannot print greeting: %v\n", err)
			cancel()
		}
	}()

	wg.Add(1)
	go func() {
		if err := printFarewell(ctx); err != nil {
			fmt.Printf("cannot print greeting: %v\n", err)
			cancel()
		}
	}()
	wg.Wait()
}

func testNonContext() {
	local := func(done <-chan interface{}) (string, error) {
		select {
		case <-done:
			return "", fmt.Errorf("Canceled")
		case <-time.After(2 * time.Second):
			return "EN/US", nil
		}
	}

	genGreeting := func(done <-chan interface{}) (string, error) {
		switch local, err := local(done); {
		case err != nil:
			return "", err
		case local == "EN/US":
			return "Hello", nil
		}
		return "", fmt.Errorf("Unsupported locale")
	}

	genFarewell := func(done <-chan interface{}) (string, error) {
		switch local, err := local(done); {
		case err != nil:
			return "", err
		case local == "EN/US":
			return "Goodbye", nil
		}
		return "", fmt.Errorf("Unsupported locale")
	}

	printGreeting := func(done <-chan interface{}) error {
		greeting, err := genGreeting(done)
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println(greeting, "world!")
		return nil
	}

	printFarewell := func(done <-chan interface{}) error {
		farewell, err := genFarewell(done)
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println(farewell)
		return nil
	}

	var wg sync.WaitGroup
	done := make(chan interface{})
	defer close(done)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := printGreeting(done); err != nil {
			fmt.Println(err)
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := printFarewell(done); err != nil {
			fmt.Println(err)
			return
		}
	}()
	wg.Wait()
}
