package main

import (
	"container/heap"
	"fmt"
	"math/rand"
	"time"
)

type Request struct {
	fn func() int // The operation to perform.
	c  chan int   // The channel to return the result.
}

func workFn() int {
	fmt.Println("I did something for the balancer.")
	return 1
}

func furtherProcess(c int) {
	fmt.Println("This is the processor of ", c)
}

// An artificial but illustrative simulation of a requester, a load generator.
// work is a send-only channel, once set, Banlancer can start to dispatch
func requester(work chan<- Request) {
	c := make(chan int)
	for {
		// Kill some time (fake load).
		time.Sleep(time.Duration(rand.Intn(1e3)) * time.Microsecond)
		fmt.Println("Will create a new request, waiting for a worker ...")
		work <- Request{workFn, c} // send request
		result := <-c              // wait for answer
		fmt.Println("Request has been processed, will send to furtherProcess()")
		furtherProcess(result)
		fmt.Println("furtherProcess has finished too. Full life cycle of a request is done.")
	}
}

type Pool []*Worker

func (p Pool) Len() int { return len(p) }

func (p Pool) Less(i, j int) bool {
	return p[i].pending < p[j].pending
}

func (p *Pool) Pop() interface{} {
	old := *p
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
	// fmt.Println("popped out", n-1, "its index was", item.index)
	item.index = -1 // for safety
	*p = old[0 : n-1]
	return item
}

func (p *Pool) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	// fmt.Println("Pushing in", x.(*Worker).index)
	n := len(*p)
	item := x.(*Worker)
	item.index = n
	*p = append(*p, item)
}

func (p Pool) Swap(i, j int) {
	// fmt.Println("Swapping", i, j)
	p[i], p[j] = p[j], p[i]
	p[i].index = i
	p[j].index = j
}

type Worker struct {
	requests chan Request // work to do (buffered channel)
	pending  int          // count of pending tasks
	index    int          // index in the heap
}

func (w *Worker) work(done chan *Worker) {
	for {
		fmt.Println("The worker ready for requests")
		req := <-w.requests // get Request from balancer
		fmt.Println("Request is being received from the channel for the worker. Request.fn starts.")
		req.c <- req.fn() // call fn and send result
		fmt.Println("Worker has sent result to Request's channel. Next, tell balancer it is done.")
		done <- w // we've finished this request
		fmt.Println("Balancer has been notified from a worker.")
	}
}

type Balancer struct {
	pool Pool
	done chan *Worker
}

func (b *Balancer) balance(work chan Request) {
	n := 1
	for {
		select {
		case req := <-work: // received a Request...
			fmt.Println(n, "Balancer received request. Start to dispatch ...")
			b.dispatch(req) // ...so send it to a Worker
		case w := <-b.done: // a worker has finished ...
			fmt.Println(n, "Balancer received the signal of Done. Cleaning up ...")
			b.completed(w) // ...so update its info
			n++
		}
	}
}

// Send Request to worker
func (b *Balancer) dispatch(req Request) {
	fmt.Println("Getting a worker from the pool.")
	// Grab the least loaded worker...``
	w := heap.Pop(&b.pool).(*Worker)
	// ...send it the task.
	w.requests <- req
	// One more in its work queue.
	w.pending++
	// Put it into its place on the heap.
	heap.Push(&b.pool, w)
	fmt.Println("Worker", w.index, "has been popped, dispatched")
}

// Job is complete; update heap
func (b *Balancer) completed(w *Worker) {
	// One fewer in the queue.
	w.pending--
	// Remove it from heap.
	heap.Remove(&b.pool, w.index)
	// Put it into its place on the heap.
	heap.Push(&b.pool, w)
	fmt.Println("Cleanup done, and push the worker", w.index, "back to the pool for new requests.\n\n")
}

// if there is a send only chanel, how not to block
// in := make(chan int)
// go dummy(in)
// fmt.Println(<-in)
func dummy(i chan<- int) {
	i <- 1
}

func main() {
	// // Just a demonstration if there is no Worker and Balancer, how a Request
	// // which was generated from a load generator is processed.
	// r := make(chan Request)
	// go requester(r)
	// req := <-r
	// req.c <- req.fn()
	// // there is no wait for furtherProcess, so sleep a bit to let furtherProcess to finish
	// time.Sleep(1 * time.Second)
	// End of the simple dome

	workers := 5
	wp := make(Pool, workers)

	for i := 0; i < workers; i++ {
		wp[i] = &Worker{
			requests: make(chan Request), // this is an unbuffered channel
			index:    i,
		}
	}

	heap.Init(&wp)

	b := Balancer{
		wp,
		make(chan *Worker),
	}
	for i := 0; i < workers; i++ {
		go wp[i].work(b.done)
	}

	r := make(chan Request)
	// set up channel, it has to be done through goroutine
	go requester(r)
	go b.balance(r)

	boom := time.After(1 * time.Second)
	<-boom
	fmt.Println("Too much, going home. BOOM!")
}
