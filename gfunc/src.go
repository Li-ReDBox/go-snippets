package src

import (
	"fmt"
	"net/http"
)

func Demo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, this is main demo")
}
