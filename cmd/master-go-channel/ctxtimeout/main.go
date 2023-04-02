package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	signal := make(chan bool)
	go expensiveTask(ctx, signal)

	select {
	case <-ctx.Done():
		fmt.Println("Expensive task took too long to complete")
		return
	case <-signal:
		fmt.Println("Expensive task was completed on time")
	}
}

func expensiveTask(ctx context.Context, signal chan<- bool) {
	time.Sleep(6 * time.Second)
	signal <- true
}
