// heap
// This example demonstrates a priority queue built using the heap interface.
package main

import (
	"container/heap"
	"fmt"
	"math/rand"
)

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
	pending int // count of pending tasks
	index   int // index in the heap
}

func main() {
	workers := 5
	wp := make(Pool, workers)

	i := 0

	for i = 0; i < workers; i++ {
		wp[i] = &Worker{
			pending: rand.Intn(300),
			index:   i,
		}
		fmt.Println(wp[i].index, ", pending: ", wp[i].pending)
	}

	heap.Init(&wp)
	// for i = 0; i < workers; i++ {
	// 	fmt.Println("Worker's pending", wp[i].pending, wp[i].index)
	// }

	// heap.Push(&wp, &Worker{pending: 3})

	for i = 0; i < workers; i++ {
		w := heap.Pop(&wp).(*Worker)
		fmt.Println("Worker's pending", w.pending)
	}

	// fmt.Println("Final length:", len(wp))
}
