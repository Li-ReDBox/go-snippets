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

// waitSome checks the channel certain times. Once done, notify caller var a channel.
func waitSome(c chan string, count int, d chan bool) {
	n := 0
	for n < count {
		select {
		case m := <-c:
			fmt.Println(m)
			n++
		}
	}
	d <- true
}

func main() {
	c := make(chan string)
	start := time.Now()
	total := 10
	// Call up 10 workers, the order of goroutines run is not certain
	for i := 0; i < total; i++ {
		go randomDelay(i, c)
	}

	// // This section demonstrates a blocking wait
	// fmt.Println("Now read from them all")
	// // Workers finish in the order of which takes less time, communicates through a single channel
	// for i := 0; i < 10; i++ {
	// 	fmt.Println(<-c)
	// }
	// fmt.Println("Total run time should be close to the time of the slowest worker", time.Since(start))

	// This section demonstrates a non-blocking wait, it is slower
	w := make(chan bool)
	// Wait is in background
	go waitSome(c, total, w)

	// check if total number of tasks finishes and if not, do other things
	for {
		select {
		case <-w:
			fmt.Println("\nTotal run time should be close to the time of the slowest worker + some overhead", time.Since(start))
			return
		default:
			fmt.Println(". ")
			time.Sleep(200 * time.Millisecond)
		}
	}
}
