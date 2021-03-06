// Package main contains code of an example with one routine
// runs in background. We use un-buffer channels, so they work
// in sequences. But whether all asked jobs can be finished is not certain.
package main

import (
	"fmt"
	"log"
	"time"
)

// count is a counter of tasks
var count int

// msg is a channel for sending and receiving messages
var msg = make(chan string)

// When package starts, run independent in background. It communicates
// through channel msg
func init() {
	go independent(msg)
}

// run reads from a channel forever
// Here we show two forms of reading a channel with or without select.
// There is no difference in functionality.
func run(c chan bool) {
	for _ = range c {
		count += 1
		log.Println(count, ": job is running ...")
	}
}

// runSelect checks a channel using select. Here has only one channel, so
// it is not a correct use of it.
// there is no difference between run and runSelect
func runSelect(c chan bool) {
	for {
		select {
		case <-c:
			count += 1
			log.Println(count, ": job is running ...")
		}
	}
}

// independent runs in background when package is loaded.
// it has default case, so it is non-block.
// package's go routine, started in init()
func independent(c chan string) {
	for {
		select {
		case w := <-c:
			log.Println("I was asked to spread message:", w)
			// pretent this is a heavy func
			time.Sleep(500 * time.Nanosecond)
			log.Printf("\t\tJob is done: content `%s` has been spread.\n", w)
		default:
			fmt.Print(". ")
			// log.Println("This background stuff will take 500ns to do something.")
			// // keep screen clean for a while
			// time.Sleep(500 * time.Nanosecond)
		}
	}
}

func main() {
	msg <- "Hello, background guard is here to serve the world."

	max := 10
	fmt.Println("Please set the max times you want to sending message. Default is 10.")
	fmt.Scanf("%d", &max)
	fmt.Println("You have set to the max number of run to ", max)

	// set counter to start from 0
	count = 0
	// single un-buffered channel
	con := make(chan bool)
	go run(con)
	// go bkSelect(con)
	// check channel con max times
	for i := 0; i < max; i++ {
		con <- true
		if i%4 == 0 {
			// tell the background where we are
			msg <- fmt.Sprintf("Current counter is %d", i)
		}
	}

	// By comment out sleeping for 100ms, the chance of finishing independent goroutine is reduced.
	// time.Sleep(100 * time.Microsecond)
	log.Println("After given some spare time to finish, the difference between asked and run =", max-count)

	// Send message to the independent running in the background
	msg <- "Last message from background."
}
