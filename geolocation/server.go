package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func longRunTask(c chan bool) {
	// Demo a background processor
	for {
		select {
		case <-c:
			fmt.Println("The long task starts")
			time.Sleep(1 * time.Second)
			fmt.Println("The long task finished")
		case <-time.After(5 * time.Second):
			fmt.Println(time.Now(), "pong")
		}
	}
}

func main() {
	c := make(chan bool)
	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("Set c to true to trigger background task runner.")
		c <- true
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"key":"value"}`)
	}

	go longRunTask(c)
	http.HandleFunc("/hello", helloHandler)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
