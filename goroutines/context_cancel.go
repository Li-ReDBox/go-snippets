package main

import (
	"context"
	"flag"
	"fmt"
	"time"
)

func long(ctx context.Context) {
	// normal non-interrupt run will last 30*500ms = 15s
	start := time.Now()
	fmt.Println(start)
	for i := 0; i < 30; i++ {
		select {
		case <-ctx.Done():
			fmt.Println("Parent context called us to finish")
			fmt.Println("We have let func long run for ", time.Since(start))
			return
		default:
			fmt.Println("Running a heavy func by sleeping 500ms, run count =", i+1)
			time.Sleep(500 * time.Millisecond)
			fmt.Println("Heavy func finishied")
		}
	}
}

var runTime = flag.Int("run_time", 900, "Allowed time in ms for func long to run.")

func main() {
	flag.Parse()
	// Create a new context with cancel channel
	ctx, cancel := context.WithCancel(context.Background())

	// create a channel to ask func main to wait for func long to finish
	stop := make(chan bool)
	go func() {
		long(ctx)
		stop <- true
	}()

	// in another routine, wait for 900 ms, cancel the long run - simulate time out in http request
	go func(cancel context.CancelFunc) {
		time.Sleep(time.Duration(*runTime) * time.Millisecond)
		fmt.Printf("Will kill func long because we have let it run for %d ms\n", *runTime)
		cancel()
	}(cancel)

	// wait for the long run to finish, or it will be cancelled by another go routine
	<-stop
}
