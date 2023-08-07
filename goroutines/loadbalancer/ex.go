package main

import (
	"fmt"
	"testing"
)

func testMixed(t *testing.T) {
	// c := make(chan bool)
	m := make(map[string]string)
	go func() {
		m["1"] = "a" // First conflicting access.
		// c <- true
	}()
	m["2"] = "b" // Second conflicting access.
	// <-c
	for k, v := range m {
		fmt.Println(k, v)
	}
	t.Log("All good")
}

func main() {
	n := 0
	for i := 0; i < 100; i++ {
		if i%3 != 0 {
			n++
			fmt.Println(i)
		}
	}
	fmt.Println("Total count", n)
}
