package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func longTask() {
	fmt.Println("The long task starts")
	time.Sleep(1 * time.Second)
	fmt.Println("The long task finished")
}

func main() {
	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		longTask()
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"key":"value"}`)
	}

	http.HandleFunc("/hello", helloHandler)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
