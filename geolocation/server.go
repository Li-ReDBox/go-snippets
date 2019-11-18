package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

var c chan bool

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

func redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://google.com.au", 307)
}

// Below code intend to response to caller as soon as possible, and notify caller once the long run
// background task finishes using Flush().
// In reality, Cloud Function will not send response after Flush().
// Also, browsers behave differently to Flush(). there is no guarantee the content will be rendered
// when it is received in pieces.
func balanced(w http.ResponseWriter, r *http.Request) {
	c := make(chan bool)
	go func(c chan bool) {
		time.Sleep(3 * time.Second)
		c <- true
	}(c)
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintln(w, "Hi, I am speaking to you, but will redirect once things have been done.")
	f, _ := w.(http.Flusher)
	f.Flush()
	fmt.Fprintln(w, <-c)
}

func main() {
	c = make(chan bool, 10)
	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("Set c to true to trigger background task runner.")
		c <- true
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"key":"value"}`)
	}

	go longRunTask(c)
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/red", redirect)
	http.HandleFunc("/sleep", balanced)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
