package main

import (
	"fmt"
	"time"
)

var cache string

// Use a cache and the same interface for retrieving data
// either compute it from scratch or from cache
func readFromSomewhere() chan string {
	c := make(chan string)
	go func() {
		if len(cache) > 0 {
			fmt.Println("We have already had it, put it to channel now")
		} else {
			fmt.Println("Pretend to run a IO heavy stuff somewhere.")
			time.Sleep(1 * time.Second)
			cache = "your data is here"
		}
		c <- cache
	}()
	return c
}

func main() {
	question := readFromSomewhere()
	// no reader, no wait, so discard this, let main to waste it
	readFromSomewhere()
	wait := make(chan bool)
	for _ = range [5]int{} {
		go func(w chan bool) {
			fmt.Println(<-readFromSomewhere())
			w <- true
		}(wait)
		_ = <-wait
	}
	// consume once
	fmt.Println("Print previously received content:")
	fmt.Println(<-question)
}
