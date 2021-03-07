// goroutines run independently. So when main goroutine stops, other goroutines
// called from main goroutine are killed. goroutines communicate through channels.
// un-buffered channels block. Once has been set for sending, the
// goroutine is blocked. The receiving goroutine is blocked when it
// starts receiving from a channel if the channel is empty.
// Channels must have been set up before they can be received. Compiler does not
// check the setup.

package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

// sendingAndReceiving shows that the common way of using a channel
func sendingAndReceiving() {
	c := make(chan int)
	// sending and receiving cannot be in the same function space. This is because that
	// once a channel has sent a message, the execution in sending goroutine will be
	// blocked until the channel has been received. So it has to have another goroutine
	// for receiving. normally is main.
	go func() {
		c <- 1
	}()
	fmt.Println(<-c)
}

func blocker(c chan int) {
	select {
	case n := <-c:
		fmt.Println("Channel is readable", n)
	}
	fmt.Println("After select-case")
}

func looper(c chan int) {
	for {
		// select is blocking, same as below for one channel
		// select {
		// case n := <-c:
		// 	fmt.Println("Channel is readable", n)
		// }
		n := <-c
		fmt.Println("Channel is readable", n)
		fmt.Println("After select-case")
	}
	fmt.Println("This will not been seen")
}

// Copied from Rob Pike's talk
func boring(msg string) <-chan string { // Returns receive-only channel of strings.
	c := make(chan string)
	go func() { // We launch the goroutine from inside the function.
		for i := 0; ; i++ {
			c <- fmt.Sprintf("%s %d", msg, i)
			// we are not using a real random number, so it always generate the same sequence
			t := rand.Intn(1e3)
			fmt.Printf("Will pause for %dms\n", t)
			time.Sleep(time.Duration(t) * time.Millisecond)
		}
	}()
	return c // Return the channel to the caller.
}

// if channel does not receive message, it finishes
// Copied from Rob Pike's talk
func slower() {
	c := boring("Joe")
	for {
		select {
		case s := <-c:
			fmt.Println(s)
		case <-time.After(500 * time.Millisecond):
			fmt.Println("You're too slow.")
			return
		}
	}
}

// go routine only runs for 5s, then it finishies
// Copied from Rob Pike's talk
func timeout() {
	c := boring("Joe")
	start := time.Now()

	timeout := time.After(5 * time.Second)
	fmt.Println("Service is available for 5s. After that we are done.")
	for {
		select {
		case s := <-c:
			fmt.Println(s)
		case <-timeout:
			elapsed := time.Since(start)
			fmt.Printf("You talk too much. I have listened to you for %v seconds!\n", elapsed.Seconds())
			return
		}
	}
}

// one go routine which reads channel once, but try to send more than one
func deadlock() {
	fmt.Println("Start")
	ch := make(chan int)
	// go blocker(ch)
	fmt.Println("We are done in main")
	ch <- 0
	// one go routine, next channel will not have a receiver, panic:
	// fatal error: all goroutines are asleep - deadlock!

	ch <- 1
	ch <- 2
}

// loop the channel forever
func nonStop() {
	fmt.Println("Start")
	ch := make(chan int)
	go looper(ch)
	fmt.Println("We are done in main")
	ch <- 0
	// one go routine, next channel will not have a receiver, panic:
	// fatal error: all goroutines are asleep - deadlock!

	// goroutine 1 [chan send]:
	// main.main()
	// 	/home/li/Documents/learning/golang/github_go-snippets/goroutines/blocks.go:26 +0x145
	// exit status 2
	time.Sleep(3 * time.Second)
	ch <- 1
	// the last one, there is no guarantee to finish when there is something blocks
	time.Sleep(1 * time.Second)
	ch <- 2
}

// sync demos how to synchronise goroutines
func sync() {
	// Total run time is 3 seconds, so anything takes longer thant this will be killed.
	start := time.Now()

	// run a time consuming routine without linking to main, it cannot finish.
	// Once main is done, this will be killed
	go func() {
		log.Println("Cannot finish: It will try to outlive the main func, which will not be possible")
		time.Sleep(5 * time.Second)
		log.Println("Cannot finish: out live main is not possible. Wrong!!!")
	}()

	c := make(chan string)

	go func(msg chan string) {
		log.Println("Deciding one: Will do my thing which takes 3 seconds.")
		time.Sleep(3000 * time.Millisecond)
		msg <- "Deciding one communicated"
		log.Println("Deciding one: Fully done")
	}(c)

	go func(msg chan string) {
		time.Sleep(1000 * time.Millisecond)
		log.Println("Long run: This function will jump line to communicate - it takes only 1 sec to communicate")
		fmt.Println("\t\t", "but after communication, it takes another 3s to do other things, so the last message will not be seen.")
		msg <- "Long run communicated"
		// this goroutine will do other stuff for 3 seconds, so it will not finish when
		// the main goroutine stops
		time.Sleep(3 * time.Second)
		// rand.Seed(time.Now().Unix())
		// time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
		log.Println("Long run: Jumping line routine done")
	}(c)

	log.Println("All have been set up, lets start.")

	for i := 0; i < 2; i++ {
		log.Println(<-c)
	}
	// log.Println(<-c)
	// log.Println(<-c)

	// now two properly set up channels are empty, if call another time, there is panic
	// log.Println(<-c)
	log.Println("main goroutine is leaving, but long run's last message has not been received.")
	log.Println("Running time is", time.Since(start))
}

func main() {

	// // one go routine which read channel forever, no error,
	// // but the last one may not have chance to run
	// nonStop()

	// // slower will be cut out after a channel has been process but no new message received
	// slower()

	// gorountine only serves certain time, then it shuts down
	timeout()
	// // one goroutine which read channel once, if send more than once, error
	// deadlock()

	// synchronise goroutines
	sync()
}
