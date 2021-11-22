package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"sync"
)

type Creator func(content, name string)

func createAll(root, pattern string, creator Creator) {
	fs, err := filepath.Glob(root + "/" + pattern)
	if err != nil {
		log.Fatalf("Cannot list with pattern of *.txt. Detail %+v\n", err)
	}

	if len(fs) == 0 {
		fmt.Println("No suitable files found for pattern", pattern)
		return
	}

	// We use sync.WaitGroup to have multiple goroutine run and only finishes when they all finish.
	var wg sync.WaitGroup

	for _, f := range fs {
		wg.Add(1)
		go func(n string) {
			// We can chain a few steps here
			defer wg.Done()
			content := readFile(n)
			fmt.Println(content)
			file := filepath.Base(n)
			name := strings.Split(file, ".")[0]
			creator(content, name)
		}(f)
	}
	wg.Wait()
}

func readFile(fn string) string {
	fmt.Println("Reading file", fn)
	content, err := ioutil.ReadFile(fn)
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
