package main

import "fmt"

// A Tree is a binary tree with integer values.
type Tree struct {
	Left  *Tree
	Value int
	Right *Tree
}

func main() {
	var a [22]int
	fmt.Println(a[0])
}
