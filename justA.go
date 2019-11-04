package main

import (
	"fmt"
	"math/rand"
)

func Abs(value int) int {
	if value >= 0 {
		return value
	}
	return -value
}

func main() {
	for _, value := range rand.Perm(6) {
		fmt.Println(value)
	}
	for i := 0; i < 5; i++ {
		y := (i + 5) * 7
		fmt.Println(y)
	}
}
