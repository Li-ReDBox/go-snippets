package main

import (
	"fmt"
	"time"
)

// f gets an int value from a channel (right), plus 1 and passes on to
// another channel (left). It sets up one piece of a chain.
// from Rob Pike's talk
func f(left, right chan int) {
	left <- 1 + <-right
}

// from Rob Pike's talk
func original() {
	const n = 10000
	leftmost := make(chan int)
	right := leftmost
	left := leftmost
	for i := 0; i < n; i++ {
		right = make(chan int)
		go f(left, right)
		left = right
	}
	go func(c chan int) { c <- 1 }(right)
	fmt.Println(<-leftmost)
}

// A simplified version of Rob's
func demo() {
	const n = 10000
	// Starts from left, leftmost is left at the beginning
	left := make(chan int)
	leftmost := left
	var right chan int
	// set up a chain from left to right
	for i := 0; i < n; i++ {
		right = make(chan int)
		go f(left, right)
		left = right
	}
	// the last right is rightmost
	// go func(c chan int) { c <- 1 }(right)
	right <- 1
	fmt.Println(<-leftmost)
}

// This is what daisy chain has been setup for solving.
func straight() {
	const n = 10000
	right := 1
	left := right
	for i := 0; i < n; i++ {
		left = right + 1
		right = left
	}
	fmt.Println("By conventional", left)
}

func main() {
	start := time.Now()
	demo()
	elapsed := time.Since(start)
	fmt.Printf("Time took using channels was %v\n", elapsed.Microseconds())
	start = time.Now()
	straight()
	elapsed = time.Since(start)
	fmt.Printf("Time took by straight addition was %v\n", elapsed.Microseconds())
}
