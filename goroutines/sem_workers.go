// Inspired by Bryan C. Mills's 2018 talk: Rethinking Classical Concurrency Patterns
package main

import (
	"fmt"
	"time"
)

const limit = 50

type token struct{}

func do(s int) {
	fmt.Println("You have said", s)
	time.Sleep(100 * time.Millisecond)
}

func main() {
	sem := make(chan token, limit)

	start := time.Now()
	for task := 1; task <= 100; task++ {
		sem <- token{} // start the counter
		go func(task int) {
			do(task)
			<-sem // reduce the counter
		}(task)
	}

	for n := limit; n > 0; n-- {
		sem <- token{}
	}

	fmt.Println("Done and goodbye after", time.Since(start))
}
