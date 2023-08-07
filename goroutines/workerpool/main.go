package main

import (
	"fmt"
	"math/rand"
	"runtime"

	// "sync"
	"time"
)

func main() {
	demoPool2()
}

func work(id, load int) {
	fmt.Println("    Working on job", id, "with load of", load, "ms")
	time.Sleep(time.Duration(load) * time.Millisecond)
}

// One worker, work in sequential
// func main() {
// 	fmt.Println("Hello workers")

// 	for i := 0; i < 10; i++ {
// 		load := rand.Intn(1e2)
// 		work(i, load)
// 	}
// }

// This is naive demo of running on-demand number of goroutines
// In real world example, there will be scheduler to recive jobs and create workers in goroutines.
// // runs in sequential using goroutines
// // when there are too many jobs, there will be as many as goroutines
// // there is a possiblity if there is no upper limit, there will be too many goroutines
// func main() {
// 	// create a worker pool
// 	wg := &sync.WaitGroup{}
// 	for i := 0; i < 10; i++ {
// 		jobId := i
// 		load := rand.Intn(1e2)
// 		wg.Add(1)
// 		go func(jobId, load int) {
// 			work(jobId, load)
// 			wg.Done()
// 		}(jobId, load)
// 	}
// 	wg.Wait()
// 	fmt.Println("Completed runs in sequentia using goroutines")
// 	// Unordered Output:
// 	// Working on  3 with load of 59 ms
// 	// Working on  4 with load of 81 ms
// 	// Working on  9 with load of 0 ms
// 	// Working on  6 with load of 25 ms
// 	// Working on  7 with load of 40 ms
// 	// Working on  1 with load of 87 ms
// 	// Working on  5 with load of 18 ms
// 	// Working on  8 with load of 56 ms
// 	// Working on  2 with load of 47 ms
// 	// Working on  0 with load of 81 ms
// 	// Completed runs in sequentia using goroutines
// }

type Job struct {
	ID, Load int
}

func worker(workerID int, job chan Job, results chan string) {
	for detail := range job {
		fmt.Println("Worker", workerID, "received job", detail.ID, detail.Load)
		work(detail.ID, detail.Load)
		results <- fmt.Sprintf("Worker %d has done job %d", workerID, detail.ID)
	}
}

// func main() {
// 	const (
// 		poolSize = 3
// 		jobCount = 10
// 	)

// 	// In gobyexamples, worker accept two channels: one for in, one for out
// 	// The channel are buffered by the number of jobs to run. This is not
// 	// great, isn't?

// 	jobs := make(chan Job, 2*poolSize+1)
// 	results := make(chan string)

// 	for i := 0; i < poolSize; i++ {
// 		workerId := i
// 		go worker(workerId, jobs, results)
// 	}

// 	for j := 0; j < jobCount; j++ {
// 		job := Job{
// 			ID:   j,
// 			Load: rand.Intn(1e3),
// 		}
// 		jobs <- job
// 	}
// 	close(jobs)

// 	// wait for them to finish
// 	for j := 0; j < jobCount; j++ {
// 		fmt.Println(<-results)
// 	}
// 	fmt.Println("Completed runs with a worker pool of size", poolSize)
// }

// // Worker pool
// func pool(job chan Job)  {
// 	poolSize := 3
// 	workers := []chan bool
// 	// A worker pool of size of 3

// 	for i := 0; i < poolSize; i++ {
// 		ch := worker(job)
// 		workers.append(ch)
// 	}
// 	for {
// 		select wor
// 	}
// }

func demoPool2() {
	const (
		poolSize = 3
		jobCount = 7
	)

	// In gobyexamples, worker accept two channels: one for in, one for out
	// The channel are buffered by the number of jobs to run. This is not
	// great, isn't?

	// Here we use non-buffered channels
	jobs := make(chan Job)
	results := make(chan string)

	// We create a pool of poolSize worker
	var i int
	for i = 0; i < poolSize; i++ {
		workerId := i
		go worker(workerId, jobs, results)
	}

	var job Job
	i = 0
	job = Job{
		ID:   i,
		Load: rand.Intn(1e3),
	}

	// There are poolSize of workers, so at the start, we can send poolSize jobs without deadlock
	// After that we have to monitoring the channels to make sure they are moving forward
	resultCount := 0
	for i = 0; i < jobCount; {
		select {
		case jobs <- job:
			i++
			job = Job{
				ID:   i,
				Load: rand.Intn(1e3),
			}
		case r := <-results:
			fmt.Println(r)
			resultCount = resultCount + 1

		default:
			time.Sleep(5 * time.Microsecond)
		}
	}

	// Once all jobs are sent to workers, there will be poolSize number of results have not retrieved:
	fmt.Println("so far", resultCount, "has done, the rest will be retrieved below:")
	for i = resultCount; i < jobCount; i++ {
		fmt.Println(<-results)
	}
	fmt.Printf("\nNumber of gorountine before closing jobs channel: %d\n", runtime.NumGoroutine()-1)

	// close jobs channel, so worker routines can exit.
	// In this example, it is not critical but will be critical when this goroutine is created by others
	// and run for a while
	close(jobs)

	// give runtime a bit of time to let it to do it jobs
	time.Sleep(100 * time.Millisecond)

	fmt.Println("Completed runs with a worker pool of size", poolSize)
	fmt.Println("Number of gorountine now:", runtime.NumGoroutine()-1)
}

type Pool struct {
	Size    int
	Jobs    chan Job
	Results chan string
	counter int
}

func (p *Pool) Init() {
	p.Jobs = make(chan Job, p.Size+1)
	p.Results = make(chan string)
	p.counter = 0

	for i := 0; i < p.Size; i++ {
		workID := i
		go worker(workID, p.Jobs, p.Results)
	}
}

func (p *Pool) Process(j Job) {
	p.counter += 1
	p.Jobs <- j
}

// func (p *Pool) Conclude() {
// 	for r := range  {
// 		fmt.Println(<-p.Results)
// 	}
// }

// demoJobProcessing demos a simple worker pool processes jobs.
func demoJobProcessing() {
	// create a size of 2 pool
	p := &Pool{Size: 2}
	p.Init()

	// we have five jobs to send out
	for i := 0; i < 5; i++ {
		job := Job{
			ID:   i,
			Load: rand.Intn(1e3),
		}

		p.Process(job)
	}
	//go p.Conclude()
	close(p.Jobs)
	for i := 0; i < 5; i++ {
		fmt.Println(<-p.Results)
	}
}
