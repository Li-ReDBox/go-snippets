package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
)

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
		log.Printf("This is in midware - %s\t%s", req.RemoteAddr, req.URL)
		// continue on to server the request
		h.ServeHTTP(w, req)
		// once the wrapped Handler done, we can do other things aftwards.
		// But we should not write to w
	})
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

	httpAddr := ":8000"
	fmt.Printf("Serving on port %s\n", httpAddr)
	if err := http.ListenAndServe(httpAddr, handler); err != nil {
		log.Fatalf("ListenAndServe %s: %v", httpAddr, err)
	}
}
