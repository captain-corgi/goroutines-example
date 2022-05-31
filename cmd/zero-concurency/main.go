package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Create wait group
var wg sync.WaitGroup

// Create mutex
var mutex sync.Mutex

func main() {
	websites := []string{
		"https://stackoverflow.com/",
		"https://github.com/",
		"https://www.linkedin.com/",
		"http://medium.com/",
		"https://golang.org/",
		"https://www.udemy.com/",
		"https://www.coursera.org/",
		"https://wesionary.team/",
		"https://www.facebook.com/",
	}

	// Start timer
	start := time.Now()

	// Loop through websites
	for _, website := range websites {
		go getWebsite(website)
		// Add wg
		wg.Add(1)
	}
	// Wait
	wg.Wait()

	// End timer
	elapsed := time.Since(start)
	fmt.Printf("Time elapsed: %s\n", elapsed)
}

// getWebsite makes an HTTP request to a website
func getWebsite(url string) {
	// Defer end wait group
	defer wg.Done()

	// Make HTTP request
	if res, err := http.Get(url); err != nil {
		fmt.Println(url, "is down")

	} else {
		// Get lock
		mutex.Lock()
		// Defer unlock
		defer mutex.Unlock()

		fmt.Printf("[%d] %s is up\n", res.StatusCode, url)
	}
}
