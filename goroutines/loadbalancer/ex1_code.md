# INTRO

Concurrency is a domain I have wanted to explore for a long time because the locks and the race conditions have always intimidated me. I recall somebody suggesting concurrency patterns in golang because they said "you share the data and not the variables".

Amused by that, I searched for "concurrency in golang" and bumped into this awesome slide by Rob Pike: https://talks.golang.org/2012/waza.slide#1 which does a great job of explaining channels, concurrency patterns and a mini-architecture of load-balancer (also explains the above one-liner).

Let's dig in:


# Goroutines

Slide #32: 

- Goroutines're a bit like threads, but they're much cheaper.
- When a goroutine blocks, that thread blocks but no other goroutine blocks.
- You spawn a go-routine using `go func()`.

Cool let's try it out.

``` go
package main

import (
   "fmt"
)

func func1() {
  fmt.Println("boom")
}

func main() {
  go func1()
}
```
```bash
go build main.go
```

``` bash
./main
```

Huh, it outputs nothing.
So we look at point 2 of go-routines: 

- "When a goroutine blocks, that thread blocks but no other goroutine blocks."

Umm..so the go-routine of func1() was spawned but it didn't get time to execute?

Let's add some delay.

``` go
package main

import (
   "fmt"
   "time"
)

func func1() {
  fmt.Println("boom")
}

func main() {
  go func1()
  time.Sleep(time.Millisecond)
}
```

```bash
boom
```

Awesome! Now we know if we give go-routine enough time to execute it can perform some of our concurrently.
It'd be really great if the main program could communicate with func1() and share some data-structures.

# Channels

They are the link through which we can achieve the above.
- You write in this link using `link<-data`.
- You read from this link using `data<-link`.
- You can read as long as you writing. Or else it gets blocked (and waits for the next write).

So is it a shared queue?

Let's see

```go
package main

import (
   "fmt"
)

func func1(channel chan string) {
  fmt.Println("haha")
  channel <- "boom"
}

func main() {
  channel := make(chan string)
  go func1(channel)
  fmt.Println(<-channel)
}
```

```bash
haha
boom
```
Indeed!

Yeap, the `fmt.Println(<-channel)` blocked the main program and waited for func1().
Our go-routine wrote "boom" to the shared queue and main program read from it.

# Multiple channels
Why not! We can call this function again with second channel. 
We'll then have 2 links in which main can communicate with 2 go-routines.

```go
package main

import (
   "fmt"
)

func func1(channel chan string) {
  fmt.Println("haha")
  channel <- "boom"
}

func main() {
  channel1 := make(chan string)
  channel2 := make(chan string)
  go func1(channel1)
  go func1(channel2)
  fmt.Println("channel one sends ", <-channel1)
  fmt.Println("channel two sends ", <-channel2)
}
```


```bash
haha
haha
channel one sends  boom
channel two sends  boom
```

# Switch
Slide #34:
Wouldn't it be great if we had some kind of a **switch** for go routines?
We do!
Its called select!

`select` queries each channel and the channel which is ready to be read, gets selected and we print the data.
Now, the interesting part is undecidability of the order.

Since `fmt.Println()` is an (slower) I/O operation, its almost certain that neither channel would be ready for i=0.
However, we can't really guarantee which channel will have data available first (for i=1). No, that is something the kernel decides.

The print order depends on which go-routine gets executed first.

```go
package main

import (
   "fmt"
   "time"
)

func func1(channel chan string) {
  channel <- "boom"
}

func main() {
  channel1 := make(chan string)
  channel2 := make(chan string)
  go func1(channel1)
  go func1(channel2)
  count := 0
  for ; count < 3; {
    select {
      case v := <-channel1:
        fmt.Println("channel 1 sends", v)
      case v := <-channel2:
        fmt.Println("channel 2 sends", v)
      default: // optional
        fmt.Println("neither channel was ready")
    }
    time.Sleep(time.Millisecond)
    count++;
  }
}
```

# Load balancer architecture:
- Slide #45
```
|Client1|         |Load|  <-DONE-- |Worker1| processing R3 coming from Client1
|Client2| --REQ-> |Blncr| --WOK->  |Worker2| processing R8 coming from Client2
                                   |Worker3| processing R5 coming from Client1
                          <-RESP-- response of R4 to Client2
```

# Data Flow 

- k Clients pack the value x in Request object and sends it to REQ channel.
- Load balancer blocks on REQ channel listening to Request(s).
- Load balancer chooses a worker and sends Request to one of the channels of worker WOK(i).
- Worker receives Request and processes x (say calculates sin(x) lol).
- Worker updates load balancer using DONE channel. LB uses this for load-balancing.
- Worker writes the sin(x) value in the RESP response channel (enclosed in Request object).

# Channels in play

- central REQ channel (Type: Request)
- n WOK channels (n sized worker pool, Type: Work)
- k RESP channels (k clients, Type: Float)
- n DONE channels (Type: Work)

# Client and Request

- Each client is a forever-running loop go-routine.
- In that loop, it is spawning requests that are sent to the central REQ channel linked to LB.
- For response, requests use a common channel (RESP) per client.

```go
type Request struct {
	data int
	resp chan float64
}

func createAndRequest(req chan Request) {
	resp := make(chan float64)
        // spawn requests indefinitely
	for {
        	// wait before next request
		time.Sleep(time.Duration(rand.Int63n(int64(time.Millisecond))))
		req <- Request{int(rand.Int31n(90)), resp}
		// read value from RESP channel
		<-resp
	}
}
```

# Worker and processing

- Each worker is a forever-running loop go-routine.
- In that loop, each worker is blocked on its channel trying to get Request object and then later process it.
- Worker can take multiple requests. # of pending keeps track of number of requests being executed.
- `pending` in other words means how many requests are present/being executed in the channel of each worker.

```go
type Work struct {
	// heap index
	idx        int
	// WOK channel
	wok chan Request
	// number of pending request this worker is working on
	pending  int
}

func (w *Work) doWork(done chan *Work) {
	// worker works indefinitely
	for {
		// extract request from WOK channel
		req := <-w.wok
		// write to RESP channel
		req.resp <- math.Sin(float64(req.data))
		// write to DONE channel
		done <- w
	}
}
```

# Balancer data structures

- The crux of Balancer is a heap (Pool) which balances based on number of pending requests.
- DONE channel, is used to notify heap that worker is finished and pending counter can be decremented.

```go
type Pool []*Work

type Balancer struct {
	// a pool of workers
	pool Pool
	done chan *Work
}

func InitBalancer() *Balancer {
	done := make(chan *Work, nWorker)
	// create nWorker WOK channels
	b := &Balancer{make(Pool, 0, nWorker), done}
	for i := 0; i < nWorker; i++ {
		w := &Work{wok: make(chan Request, nRequester)}
		// put them in heap
		heap.Push(&b.pool, w)
		go w.doWork(b.done)
	}
	return b
}
```

# Heap implementations

- It sucks but golang wants you to implement your own Len, Less, Push, Pop, Swap for Heap interface.
- Copied shamlessly from github (look at references below).

```go
func (p Pool) Len() int { return len(p) }

func (p Pool) Less(i, j int) bool {
	return p[i].pending < p[j].pending
}

func (p *Pool) Swap(i, j int) {
	a := *p
	a[i], a[j] = a[j], a[i]
	a[i].idx = i
	a[j].idx = j
}

func (p *Pool) Push(x interface{}) {
	n := len(*p)
	item := x.(*Work)
	item.idx = n
	*p = append(*p, item)
}

func (p *Pool) Pop() interface{} {
	old := *p
	n := len(old)
	item := old[n-1]
	item.idx = -1 // for safety
	*p = old[0 : n-1]
	return item
}
```

# Load balancing (channels)
- If the central REQ channel has a request coming in from clients, dispatch it to least loaded Worker and update the heap.
- If DONE channel reports back, the work assigned to WOK(i) has been finished.

```go
func (b *Balancer) balance(req chan Request) {
	for {
		select {
		// extract request from REQ channel
		case request := <-req:
			b.dispatch(request)
		// read from DONE channel
		case w := <-b.done:
			b.completed(w)
		}
		b.print()
	}
}

func (b *Balancer) dispatch(req Request) {
	// Grab least loaded worker
	w := heap.Pop(&b.pool).(*Work)
	w.wok <- req
	w.pending++
	// Put it back into heap while it is working
	heap.Push(&b.pool, w)
}

func (b *Balancer) completed(w *Work) {
	w.pending--
	// remove from heap
	heap.Remove(&b.pool, w.idx)
	// Put it back
	heap.Push(&b.pool, w)
}
```

# Glueing the code

- Imports and main
- Adding a print

```go
package main

import (
	"container/heap"
	"fmt"
	"math"
	"math/rand"
	"time"
)

const nRequester = 100
const nWorker = 10

func (b *Balancer) print() {
	sum := 0
	sumsq := 0
	// Print pending stats for each worker
	for _, w := range b.pool {
		fmt.Printf("%d ", w.pending)
		sum += w.pending
		sumsq += w.pending * w.pending
	}
	// Print avg for worker pool
	avg := float64(sum) / float64(len(b.pool))
	variance := float64(sumsq)/float64(len(b.pool)) - avg*avg
	fmt.Printf(" %.2f %.2f\n", avg, variance)
}

func main() {
	work := make(chan Request)
	for i := 0; i < nRequester; i++ {
		go createAndRequest(work)
	}
	InitBalancer().balance(work)
}
```

# Output

- Here you can see number of pending tasks per worker.
- Since work is just computing a sine value, I had to reduce sleep-time at the client level before they fire next request.

```bash
0 1 2 3 4 5 6 7 8 9  avg  variance

5 6 8 8 8 8 8 8 8 8  7.50 1.05
4 6 8 8 8 8 8 8 8 8  7.40 1.64
3 6 8 8 8 8 8 8 8 8  7.30 2.41
2 6 8 8 8 8 8 8 8 8  7.20 3.36
1 6 8 8 8 8 8 8 8 8  7.10 4.49
1 5 8 8 8 8 8 8 8 8  7.00 4.80
1 5 8 8 7 8 8 8 8 8  6.90 4.69
1 5 8 8 6 8 8 8 8 8  6.80 4.76
1 4 8 8 6 8 8 8 8 8  6.70 5.21
1 4 8 8 6 8 8 8 8 7  6.60 5.04
1 4 8 7 6 8 8 8 8 7  6.50 4.85
1 4 8 7 6 8 8 8 7 7  6.40 4.64
1 4 7 7 6 8 8 8 7 7  6.30 4.41
```

# Footnote

- Although this is still a single-process LB, it makes you appreciate the flexibility of asynchronous behavior.
- How channels communicate in form of a light-weight queue and offloading the tasks in form of goroutines is pretty amazing to me.
- Also, all the book-keeping of acquiring/releasing a lock is hidden from programmer and all you need to focus on "sharing data using channels and not the variables" ;).
- Event driven architecture is amazing concept. There's a nice writeup on event-driven architecture I read on hackernews the other day that tells you when and when not to use it: https://herbertograca.com/2017/10/05/event-driven-architecture/

# References
- Rob Pike's awesome slide: https://talks.golang.org/2012/waza.slide
- Full code written by @angeldm here: https://gist.github.com/angeldm/2421216