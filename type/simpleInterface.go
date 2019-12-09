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

//ReadCloser a test interface
type ReadCloser interface {
	Read(b []byte) (n int, err error)
	Close()
}

//ReadAndClose is just a test function for ReadCloser interface, and struct with Read and Close function can be passed into this function
func ReadAndClose(r ReadCloser, buf []byte) (n int, err error) {
	var nr int
	for len(buf) > 0 && err != nil {
		nr, err = r.Read(buf)
		n += nr
		buf = buf[nr:]
	}
	r.Close()
	return n, err
}
