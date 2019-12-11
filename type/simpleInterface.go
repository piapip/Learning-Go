package main

import (
	"fmt"
)

//Animal is an interface for animals
type Animal interface {
	Speak() string
}

//Dog a test struct
type Dog struct {
}

//Speak for dog
func (d Dog) Speak() string {
	return "Woof woof"
}

//Cat a test struct
type Cat struct {
}

//Speak for cat
func (c Cat) Speak() string {
	return "Meow meow"
}

//Snake a test struct
type Snake struct {
}

//Speak for snake
func (s Snake) Speak() string {
	return "Hiss hiss"
}

//Fox a meme
type Fox struct {
}

//Speak for snake
func (f Fox) Speak() string {
	return "??????"
}

func testAnimalInterface() {
	animals := []Animal{Dog{}, Cat{}, Snake{}, Fox{}}
	for _, animal := range animals {
		fmt.Println(animal.Speak())
	}
}

//PrintThis just a test
func PrintThis(v []interface{}) {
	for _, element := range v {
		fmt.Println(element)
	}
}

func aFailedTest() {
	names := []string{"Yoyoyo", "pipi"}
	// vals := make([]interface{}, len(names))
	// for i, v := range names {
	// 	vals[i] = v
	// }
	// PrintThis(vals)
	PrintThis(names)
}
