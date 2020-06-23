package datastore

import (
	"entities"
)

type Dog struct Animal

func (d Dog) Speak() {
	fmt.Println("Wan wan")
}

//Eat is a test func
func (d Dog) Eat() {
	fmt.Println("Chewing bone")
}