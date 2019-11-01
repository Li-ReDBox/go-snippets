package main

// import "parent/child"

import (
	"fmt"
    "parentlocal/child"
    "parentlocal/sibling"
)

func main() {
	fmt.Println("Parent say")
    child.Say()
    sibling.SSay()
}
