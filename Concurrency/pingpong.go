package main

import (
	"fmt"
	"time"
)

func player1(table chan int) {
	for {
		ball := <-table
		ball = 1 - ball
		fmt.Println("Player 1: ", ball)
		time.Sleep(100 * time.Millisecond)
		table <- ball
	}
}

func player2(table chan int) {
	for {
		ball := <-table
		ball = 1 - ball
		fmt.Println("Player 2: ", ball)
		time.Sleep(100 * time.Millisecond)
		table <- ball
	}
}

func pingpong() {
	var Ball int
	table := make(chan int)
	go player1(table)
	go player2(table)

	table <- Ball
	time.Sleep(1 * time.Second)
	<-table
}
