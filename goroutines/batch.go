package main

import (
	"fmt"
	"math/rand"
	"time"
)

func worker(d chan int) {
	rnd := rand.Intn(1e3)
	fmt.Printf("Worker will have to work for %dms\n", rnd)
	time.Sleep(time.Duration(rnd) * time.Millisecond)
	d <- rnd
}

func batch(d chan bool) {
	start := time.Now()
	fmt.Println("Start a batch")
	c := make(chan int)
	// call workers straightaway
	for i := 0; i < 10; i++ {
		go worker(c)
	}

	// now check/wait them to be finished
	n := 0
	for n < 10 {
		fmt.Println(<-c, "seconds")
		n++
	}
	fmt.Println("Total run time should be close to the time of the slowest worker + some overhead", time.Since(start))
	d <- true
}

func main() {
	// this is a timeout solution, so there will be only certain amount of batches will run/finish
	boom := time.After(2 * time.Second)

	tick := make(chan bool)
	go batch(tick)
	n := 1
	for {
		select {
		case <-tick:
			fmt.Println("One batch done. So far: ", n, "\n")
			go batch(tick)
			n++
		case <-boom:
			fmt.Println("Too much, going home. BOOM!")
			return
		}
	}
}
