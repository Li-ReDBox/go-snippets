package main

import (
	"fmt"
	"gfunc/src"
	"net/http"
)

func main() {
	http.Handle("/", http.HandlerFunc(src.Demo))
	err := http.ListenAndServe(":1718", nil)
	if err != nil {
		fmt.Println("ListenAndServe:", err)
	}
}
