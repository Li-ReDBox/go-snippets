package main

import (
	"fmt"
	"time"
)

var count int

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

func main() {
	count = 0
	con := make(chan bool)
	go background(con)
	// go bkSelect(con)
	max := 9
	for i := 0; i < max; i++ {
		con <- true
	}
	time.Sleep(100 * time.Microsecond)
	fmt.Println("After given some spare time to finish, the difference between asked and run =", max-count)
}
