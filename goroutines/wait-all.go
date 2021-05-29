package main

import (
	"fmt"
	"math/rand"
	"time"
)

// randomDelay is a worker function takes random number of ms to finish
func randomDelay(o int, c chan string) {
	rnd := time.Duration(rand.Intn(1e3))
	fmt.Printf("This is func%d to sleep %dms\n", o, rnd)
	time.Sleep(rnd * time.Millisecond)
	c <- fmt.Sprintf("Slept for %d ms in func%d", rnd, o)
}

func main() {
	c := make(chan string)
	start := time.Now()
	// Call up 10 workers, the order of goroutines run is not certain
	for i := 0; i < 10; i++ {
		go randomDelay(i, c)
	}
	fmt.Println("Now read from them all")

	// Workers finish in the order of which takes less time, communicates through a single channel
	for i := 0; i < 10; i++ {
		fmt.Println(<-c)
	}
	fmt.Println("Total run time should be close to the time of the slowest worker", time.Since(start))
}
