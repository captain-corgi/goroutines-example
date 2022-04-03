package main

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

func main() {
	// Start timer
	start := time.Now()
	defer func() {
		// End timer
		elapsed := time.Since(start)
		fmt.Printf("Time elapsed: %s\n", elapsed)
	}()

	// Make string channel
	c := make(chan string)

	websites := []string{
		"https://stackoverflow.com/",
		"https://github.com/",
		"https://www.linkedin.com/",
		"http://medium.com/",
		"https://golang.org/",
		"https://www.udemy.com/",
		"https://www.coursera.org/",
		"https://wesionary.team/",
		"https://conmeo.com.vn/",
	}

	// Loop through websites
	for _, website := range websites {
		go getWebsite(website, c)
	}

	// Iterating over the range of channel. So keeps receiving messages until channel is closed
	for i := 0; i < len(websites); i++ {
		fmt.Println(<-c)
	}

	// for msg := range c {
	// 	fmt.Println(msg)
	// }
	// Alternative way to iterate over channel
	// for {
	// 	msg, open := <-c
	// 	if !open {
	// 		break
	// 	}
	// 	fmt.Println(msg)
	// }
}

func getWebsite(website string, c chan string) {
	if res, err := http.Get(website); err != nil {
		c <- fmt.Sprintf("%s is down", website)
	} else {
		// Resolve domain to ip
		ip, _ := net.LookupIP(res.Request.URL.Hostname())
		c <- fmt.Sprintf("[%d] %s is up, \t\tIP: %s", res.StatusCode, website, ip[0])
	}
}
