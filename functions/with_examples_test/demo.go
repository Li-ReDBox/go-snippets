package demo

import (
    "fmt"
    "net/http"
)

func Display(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "This is good")
}
