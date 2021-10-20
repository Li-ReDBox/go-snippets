package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func findAll(pattern string) {
	fs, err := filepath.Glob(pattern)
	if err != nil {
		log.Fatalf("Cannot list with pattern of *.txt. Detail %+v\n", err)
	}

	// We use sync.WaitGroup to have multiple goroutine run and only finishes when they all finish.
	var wg sync.WaitGroup

	for _, f := range fs {
		wg.Add(1)
		go func(n string) {
			// We can chain a few steps here
			defer wg.Done()
			s := readFile(n)
			processContent(s)
		}(f)
	}
	wg.Wait()
}

func readFile(fn string) string {
	fmt.Println("Reading file", fn)
	content, err := os.ReadFile(fn)
	if err != nil {
		fmt.Printf("Cannot read %s, %+v\n", fn, err)
		return ""
	}

	fmt.Printf("From %s, content:\n%s\n", fn, content)
	// prepare for next step:
	parts := strings.Split(fn, ".")
	fmt.Println(parts[0])
	return string(content)
}

func processContent(content string) {
	// Only process if there is content
	if len(content) > 0 {
		s := strings.Replace(content, "{}", "%s", 1)
		fmt.Println("New content", s)
	}
}

func main() {
	findAll("*.txt")
}
