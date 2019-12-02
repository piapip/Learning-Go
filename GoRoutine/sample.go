package main

import (
	"fmt"
	"strconv"
)

//in Goroutine, there will never be a situation where 2 channels handle the same piece of information.

func main() {
	theMine := []string{"ore", "rock", "ore", "mud", "rock", "ore", "rock", "rock", "ore", "mud", "mud", "ore"}
	oreChan := make(chan string)
	minedOreChan := make(chan string)
	doneChan := make(chan string)

	go func(mine []string, founder chan<- string) {
		for count, item := range mine {
			if item == "ore" {
				fmt.Println("Sending ore number " + strconv.Itoa(count))
				founder <- item + strconv.Itoa(count)
			}
			if count == len(theMine)-1 {
				founder <- "No more ore"
			}
		}
	}(theMine, oreChan)

	go func(founder <-chan string, miner chan<- string) {
		for foundOre := range founder {
			fmt.Println("Miner: Received " + foundOre + " from finder")
			miner <- foundOre
		}
	}(oreChan, minedOreChan)

	go func(miner <-chan string, supervisor chan<- string) {
		for minedOre := range miner {
			fmt.Println("From miner: " + minedOre)
			fmt.Println("Smelting...\n" + minedOre + " has been smelted")
			if minedOre == "No more ore" {
				doneChan <- "All done"
			}
		}
	}(minedOreChan, doneChan)

	fmt.Println(<-doneChan)
	close(oreChan)
	close(minedOreChan)
	close(doneChan)

	// go func(mine []string) {
	// 	for count, item := range mine {
	// 		if item == "ore" {
	// 			fmt.Println("Sending ore number " + strconv.Itoa(count))
	// 			oreChan <- item + strconv.Itoa(count)
	// 		}
	// 		if count == len(mine)-1 {
	// 			oreChan <- "No more ore"
	// 		}
	// 	}
	// }(theMine)

	// go func() {
	// 	for foundOre := range oreChan {
	// 		if foundOre == "No more ore" {
	// 			close(oreChan)
	// 		} else {
	// 			fmt.Println("Miner: Received " + foundOre + " from finder")
	// 		}
	// 		minedOreChan <- foundOre
	// 	}
	// }()

	// go func() {
	// 	for minedOre := range minedOreChan {
	// 		if minedOre == "No more ore" {
	// 			close(minedOreChan)
	// 		} else {
	// 			fmt.Println("From miner: " + minedOre)
	// 			fmt.Println("Smelting...\n" + minedOre + " has been smelted")
	// 		}
	// 	}
	// 	doneChan <- "All done."
	// }()

	// <-doneChan
}
