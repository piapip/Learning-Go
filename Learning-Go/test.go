package main

import (
	"fmt"
)

// func stackExample() {
// 	stackSlice := make([]byte, 512)
// 	s := runtime.Stack(stackSlice, false)
// 	fmt.Printf("\n%s", stackSlice[0:s])
// }

// //First ...
// func First() {
// 	Second()
// }

// //Second ...
// func Second() {
// 	Third()
// }

// //Third ...
// func Third() {
// 	for c := 0; c < 5; c++ {
// 		pc, file, line, something := runtime.Caller(c)
// 		fileArr := strings.Split(file, "/")
// 		file = fileArr[len(fileArr)-1]
// 		fmt.Println("Caller " + strconv.Itoa(c))
// 		fmt.Printf("pc: %v\nfile: %v\nline: %v\nsomething: %v\n", pc, file, line, something)

// 		fn := runtime.FuncForPC(pc).Name()
// 		fnArr := strings.Split(fn, ".")
// 		fn = fnArr[len(fnArr)-1]
// 		fmt.Println(fn)

// 		fmt.Println(" ")
// 	}
// }

// func test(args ...string) {
// 	fmt.Println(args)
// }

// func testStartSpan(ctx context.Context) {
// 	fmt.Println(ctx)
// 	ctx, span := trace.StartSpan(ctx, "main")
// 	fmt.Println(ctx)
// 	fmt.Println(span)
// 	defer span.End()

// 	for i := 0; i < 10; i++ {
// 		sample := []string{strconv.Itoa(i)}
// 		test(sample...)
// 	}
// }

var mine1 []string = []string{"ore", "ore", "rock", "mud", "ore", "rock"}

func findOre() func() (int, int) {
	oreFound := 0
	return func() (int, int) {
		material := mine1[0]
		mine1[0] = mine1[len(mine1)-1]
		mine1[len(mine1)-1] = ""
		mine1 = mine1[:len(mine1)-1]
		// fmt.Println(mine)
		if material == "ore" {
			oreFound++
			return oreFound, 1
		}
		return oreFound, 0
	}
}

func main() {
	miner1 := findOre()

	for len(mine1) > 0 {
		fmt.Println(miner1())
	}

	myFirstChannel := make(chan string)
	myFirstChannel <- "hello"
	myVariable := <-myFirstChannel
	fmt.Println(myVariable)
}

// func main() {
// fmt.Println("######### STACK ################")
// stackExample()
// fmt.Println("\n\n######### CALLER ################")
// First()

// ctx := context.Background()
// testStartSpan(ctx)

// }
