package main

import (
	"fmt"
	"sync"
)

func main() {
	// fmt.Println(getNumber1())
	// fmt.Println(getNumber1n1())
	// fmt.Println(<-getNumber1n2())
	// fmt.Println(getNumber2())
	// fmt.Println(getNumber3())
	// fmt.Println(getNumber3n1())

	fmt.Println(getNumber())
}

func getNumber() int {
	var i = 0
	fmt.Println("increase i by 1")
	fmt.Println("increase i by 1")
	fmt.Println("increase i by 1")
	go func() {
		fmt.Println("increase i by 1")
		i = 5
	}()
	fmt.Println("increase i by 1")
	fmt.Println("increase i by 1")
	fmt.Println("increase i by 1")
	return i
}

// getNumber1 use channel to receive written data
func getNumber1() int {
	var i = 0
	iCh := make(chan int) // Create a channel
	go func() {
		i = 5
		iCh <- i // put i into channel
	}()
	return <-iCh // return value from channel (Wait operation)
}

// getNumber1n1 block write operation by a channel
func getNumber1n1() int {
	var i = 0
	iCh := make(chan interface{}) // Create a channel
	go func() {
		i = 5
		iCh <- "" // put anything into channel
	}()
	<-iCh // Wait & discard value from channel (Wait operation)
	return i
}

// getNumber1n2 return written data as a channel
func getNumber1n2() chan int {
	var i = 0
	iCh := make(chan int) // Create a channel
	go func() {
		i = 5
		iCh <- i // put i into channel
	}()
	return iCh // return channel
}

// getNumber2 block program until write operation complete using wait group
func getNumber2() int {
	var i = 0
	var wg sync.WaitGroup // Init wait group
	wg.Add(1)             // Wait 1 unit
	go func() {
		i = 5
		wg.Done() // Done 1 unit
	}()
	wg.Wait() // Block program until nothing left to wait (Wait operation)
	return i
}

// getNumber3 use embedded mutex for lock read & write operation
//
//	NOTE: Locking value just resolve race condition,
//	ensure that read/write operation on i is sync
//
//	https://www.sohamkamani.com/golang/data-races/
func getNumber3() int {
	i := &SafeInt{}
	go func() {
		i.Set(5)
	}()
	return i.Get()
}

// SafeInt a custom struct which store int value & mutex lock
type SafeInt struct {
	val int
	m   sync.Mutex
}

// Set lock variable & write value
func (r *SafeInt) Set(val int) {
	r.m.Lock()
	defer r.m.Unlock()

	r.val = val
}

// Set lock variable & return value
func (r *SafeInt) Get() int {
	r.m.Lock()
	defer r.m.Unlock()

	return r.val
}

// getNumber3n1 use single instance mutex for locking i
func getNumber3n1() int {
	var i = 0
	var m sync.Mutex

	go func() {
		m.Lock()
		defer m.Unlock()

		i = 5
	}()

	m.Lock()
	defer m.Unlock()
	return i
}
