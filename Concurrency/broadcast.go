package main

import (
	"fmt"
	"sync"
)

func testBroadCast() {
	type Button struct {
		Clicked *sync.Cond
	}
	button := Button{Clicked: sync.NewCond(&sync.Mutex{})}

	subscribe := func(c *sync.Cond, fn func()) {
		var goroutineRunning sync.WaitGroup
		goroutineRunning.Add(1)
		go func() {
			goroutineRunning.Done()
			c.L.Lock()
			defer c.L.Unlock()
			c.Wait() // This piece of shit will be in standing by mode right here.
			fn()
		}()
		goroutineRunning.Wait()
	}

	var clickRegister sync.WaitGroup

	clickRegister.Add(3)
	subscribe(button.Clicked, func() { // I mean this
		fmt.Println("Doing this")
		clickRegister.Done()
	})

	subscribe(button.Clicked, func() { // and this
		fmt.Println("Doing that")
		clickRegister.Done()
	})

	subscribe(button.Clicked, func() { //and this too will stand still at c.Wait()
		fmt.Println("Not doing anything")
		clickRegister.Done()
	})

	button.Clicked.Broadcast() //until this piece of shit is triggered
	clickRegister.Wait()
}
