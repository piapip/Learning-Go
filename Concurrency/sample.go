package main

import (
	"fmt"
	"time"
)

func main() {
	data := 0
	go func() {
		data++
	}()
	time.Sleep(1 * time.Second)
	if data == 0 {
		fmt.Println("Yo")
	}
}
