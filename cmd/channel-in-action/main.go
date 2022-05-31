package main

import (
	"fmt"
	"time"
)

func main() {
	// Create a channel
	ch := make(chan string, 1)
	// Start the greeter to provide a greeting
	go greet(ch)
	// Sleep for 5 seconds
	time.Sleep(5 * time.Second)
	fmt.Println("Main ready!")
	// Receive the greeting
	msg := <-ch
	// Sleep 2 seconds and print
	time.Sleep(2 * time.Second)
	fmt.Println("Greeting received")
	fmt.Println(msg)
}

// greet writes a greet to the given channel and then says goodbye
func greet(ch chan<- string) {
	fmt.Printf("Greeter ready! \nGreeter waiting to send greeting...\n")
	ch <- "Hello, world!"
	fmt.Println("Greeter completed sending greeting!")
}
