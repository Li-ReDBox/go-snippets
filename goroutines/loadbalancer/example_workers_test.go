// This example demonstrates a WorkerPool using the heap interface.
// go test example_workers_test.go pool.go
package main

import (
	"container/heap"
	"fmt"
)

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
