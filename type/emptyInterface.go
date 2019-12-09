package main

import "fmt"

func testFail() {
	shouldReturnError := func(v ...interface{}) {
		fmt.Println(v)
	}

	shouldReturnError(1, "v", make(chan interface{}))
}
