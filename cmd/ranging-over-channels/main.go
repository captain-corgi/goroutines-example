package main

import (
	"fmt"
	"time"
)

var (
	// greetings in many languages
	greetings = []string{"Hello!", "Bonjour!", "Hola!", "Xin chào!", "Guten Tag!"}
	// goodbyes in many languages
	goodbyes = []string{"Goodbye!", "Au revoir!", "Ciao!", "Tạm biệt!", "Adeus!"}
)

func main() {
	// Create a channel
	ch := make(chan string, 1)
	ch2 := make(chan string, 1)
	// Start the greeter to provide a greeting
	go greet(greetings, ch)
	go greet(goodbyes, ch2)
	// Sleep for 1 seconds
	time.Sleep(1 * time.Second)
	fmt.Println("Main ready!")

	for {
		select {
		case msg1, ok := <-ch:
			if !ok {
				ch = nil
				break
			}
			printGreeting(msg1)
		case msg2, ok := <-ch2:
			if !ok {
				ch2 = nil
				break
			}
			printGreeting(msg2)
		default:
			fmt.Println("No message received")
			return
		}
	}
}

// greet writes a greet to the given channel and then says goodbye
func greet(greetings []string, ch chan<- string) {
	fmt.Printf("Greeter ready! \nGreeter waiting to send greeting...\n")
	for _, msg := range greetings {
		ch <- msg
	}
	close(ch)
	fmt.Println("Greeter completed sending greeting!")
}

// printGreeting sleeps and prints the greeting given
func printGreeting(msg string) {
	// Sleep 500 ms and print
	time.Sleep(500 * time.Millisecond)
	fmt.Printf("Greeting received: %s\n", msg)
}
