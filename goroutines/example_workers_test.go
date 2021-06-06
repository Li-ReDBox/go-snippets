// This example demonstrates a WorkerPool using the heap interface.
package workerpool_test

import (
	"container/heap"
	"fmt"
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

// This example creates a WorkerPool with some items, adds and manipulates an item,
// and then removes the items with less pending pops out first.
func Example_WorkerPool() {
	pendings := []int{1, 30, 29, 15, 27}
	wp := make(Pool, len(pendings))

	for i, p := range pendings {
		wp[i] = &Worker{
			pending: p,
			index:   i,
		}
	}

	heap.Init(&wp)

	heap.Push(&wp, &Worker{pending: 3})

	for wp.Len() > 0 {
		w := heap.Pop(&wp).(*Worker)
		fmt.Printf("%d ", w.pending)
	}

	// Output:
	// 1 3 15 27 29 30
}
