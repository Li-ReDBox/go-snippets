package main

import (
	"fmt"
	"time"
)

var count int
var msg = make(chan string)

func init() {
	go independent(msg)
}

func background(c chan bool) {
	for {
		if <-c {
			count += 1
			fmt.Println(count, ": We were asked to run")
		}
	}
}

// there is no difference between background and bkSelect
func bkSelect(c chan bool) {
	for {
		select {
		case <-c:
			count += 1
			fmt.Println(count, ": We were asked to run")
		}
	}
}

func independent(c chan string) {
	for {
		select {
		case w := <-c:
			fmt.Println("I was asked to spread:", w)
		default:
			fmt.Println("I am awake to do my stuff... Done... See you in 500ns.")
			// keep screen clean for a while
			time.Sleep(500 * time.Nanosecond)
		}
	}
}

func main() {

	count = 0
	con := make(chan bool)
	go background(con)
	// go bkSelect(con)
	max := 9
	for i := 0; i < max; i++ {
		con <- true
	}
	msg <- "Fist, world is busy"

	time.Sleep(100 * time.Microsecond)
	fmt.Println("After given some spare time to finish, the difference between asked and run =", max-count)

	msg <- "Last, world is still busy"
}
