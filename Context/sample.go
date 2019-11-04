package main

import (
	"bufio"
	"context"
	"os"
	"time"
)

func main() {
	ctx := context.Background()
	ctx, cancel = context.WithCancel(ctx)

	go func() {
		s := bufio.NewScanner(os.Stdin)
		s.Scan()
		cancel()
	}()
	sleepAndTalk(ctx, 5*time.Second(), "hello")
}
