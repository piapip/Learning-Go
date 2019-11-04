package main

import (
	"fmt"
	"reflect"
	"time"
)

func main() {
	start := time.Date(1518, 11, 01, 0, 0, 0, 0, time.UTC)
	end := time.Date(1518, 10, 31, 23, 50, 0, 0, time.UTC)
	diff := start.Sub(end).Minutes()
	fmt.Println(reflect.TypeOf(diff))
}
