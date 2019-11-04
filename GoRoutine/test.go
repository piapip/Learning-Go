package main

import "fmt"

//unbuffered channel must be woken up by putting it another Goroutine. REMEMBER that main itself is already a Goroutine.
//buffered channel can work directly from the main function without creating any other Goroutine. But its size is limited.
func main() {
	mess := make(chan string)
	doneChan := make(chan string)
	go func() {
		mess <- "Yo"
		mess <- "pop"
		mess <- "No more mess"
	}()
	go func() {
		for message := range mess {
			if message == "No more mess" {
				close(mess)
			} else {
				fmt.Println(message)
			}
		}
		doneChan <- "Dekita"
	}()
	<-doneChan
}

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

// func main() {
// fmt.Println("######### STACK ################")
// stackExample()
// fmt.Println("\n\n######### CALLER ################")
// First()

// ctx := context.Background()
// testStartSpan(ctx)

// }
