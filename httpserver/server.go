package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
)

type Demo struct {
	data string
}

func (d *Demo) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "We are serving from %s with %s", req.URL, d.data)
}

// MiddleWare: stacked closures
// A middleware to provide logging of requests
func loggingHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Printf("%s\t%s", req.RemoteAddr, req.URL)
		// continue on to server the request
		h.ServeHTTP(w, req)
	})
}

// Another middleware to demonstrate chaining middleware
func midwareHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Printf("This is in midware - %s\t%s, method=%s", req.RemoteAddr, req.URL, req.Method)
		// continue on to server the request
		h.ServeHTTP(w, req)
		// once the wrapped Handler done, we can do other things aftwards.
		// But we should not write to w
	})
}

// // http.HandlerFunc type is an adapter to allow the use of
// // ordinary functions as HTTP handlers. If f is a function
// // with the appropriate signature, HandlerFunc(f) is a
// // Handler that calls f.
// type HandlerFunc func(ResponseWriter, *Request)

// // HandlerFunc.ServeHTTP calls f(w, r).
// func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
// 	f(w, r)
// }

// A Rester responds http Get, Post, Put and Delete methods
type Rester interface {
	Get(http.ResponseWriter, *http.Request) error
	Post(http.ResponseWriter, *http.Request) error
	Put(http.ResponseWriter, *http.Request) error
	Delete(http.ResponseWriter, *http.Request) error
}

// Create wraps a Rester into an http.HandlerFunc, in turn
// responds from a correct defined method responder.
func Create(rest Rester) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			fmt.Fprintln(w, "GET")
			rest.Get(w, r)
		default:
			fmt.Fprintln(w, "Not supported")
		}
	})
}

type Book struct {
	name string
}

func (b Book) Get(w http.ResponseWriter, r *http.Request) error {
	fmt.Fprintf(w, "Your book is %s", b.name)
	return nil
}

func (b Book) Post(w http.ResponseWriter, r *http.Request) error {
	fmt.Fprintf(w, "Your book  %s is being created", b.name)
	return nil
}

func (b Book) Put(w http.ResponseWriter, r *http.Request) error {
	fmt.Fprintf(w, "Your book %s is being updating", b.name)
	return nil
}

func (b Book) Delete(w http.ResponseWriter, r *http.Request) error {
	fmt.Fprintf(w, "Your book %s is being delete", b.name)
	return nil
}

func main() {
	// http.Handler: an interface contains ServeHTTP(ResponseWriter, *Request)
	// ServeMux is a router, default Handler
	var handler http.Handler = http.DefaultServeMux

	// Hook middlewares to http.DefaultServeMux

	// These are used for the router, can be hooked for a particular pattern
	// The order is important: the last run first
	// handler = loggingHandler(handler)
	// handler = midwareHandler(handler)

	// can be chained:
	handler = midwareHandler(loggingHandler(handler))

	// this is unnecessary for DefaultServeMux, just to demonstrate set up a router
	mux := http.DefaultServeMux
	mux.Handle("/resources", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(w, "Hi, this is resources")
	}))

	// This is the normal way adding to DefaultServeMux
	http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	http.Handle("/demo", &Demo{data: "Demo Data"})

	b := Book{name: "Go is great"}
	http.Handle("/rest", Create(b))

	httpAddr := ":8000"
	fmt.Printf("Serving on port %s\n", httpAddr)
	if err := http.ListenAndServe(httpAddr, handler); err != nil {
		log.Fatalf("ListenAndServe %s: %v", httpAddr, err)
	}
}
