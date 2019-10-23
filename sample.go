package main

import (
	"fmt"
	"strconv"
)

const NEED int = 4

func main() {
	theMine := []string{"ore", "rock", "ore", "mud", "rock", "ore", "rock", "rock", "ore", "mud", "mud", "ore"}
	oreChan := make(chan string)
	minedOreChan := make(chan string)
	doneChan := make(chan string)
	go func(mine []string) {
		for count, item := range theMine {
			if item == "ore" {
				fmt.Println("Sending ore number " + strconv.Itoa(count))
				oreChan <- item + strconv.Itoa(count)
			}
			if count == len(theMine)-1 {
				oreChan <- "No more ore"
			}
		}
	}(theMine)

	go func() {
		for foundOre := range oreChan {
			if foundOre == "No more ore" {
				close(oreChan)
			} else {
				fmt.Println("Miner: Received " + foundOre + " from finder")
			}
			minedOreChan <- foundOre
		}
	}()

	go func() {
		for minedOre := range minedOreChan {
			if minedOre == "No more ore" {
				close(minedOreChan)
			} else {
				fmt.Println("From miner: " + minedOre)
				fmt.Println("Smelting...\n" + minedOre + " has been smelted")
			}
		}
		doneChan <- "All done."
	}()

	<-doneChan
}
