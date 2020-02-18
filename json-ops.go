package main

import (
	"encoding/json"
	"fmt"
	"log"
)


func main() {
    fmt.Println("Unmarshal a list of map of string:string")
    blob := `[{"small":"regular"},{"large":"unrecognized"},{"small":"normal"},{"small":"large"}]`
	var inventory []map[string]string
	if err := json.Unmarshal([]byte(blob), &inventory); err != nil {
		log.Fatal(err)
	}
	fmt.Println(inventory)
	
    fmt.Println("\nMarshal a list of map of string:string to json bytes")
	x := []map[string]string{{"regular":"good"},{"unrecognized":"bad"}}
	
	js, err := json.Marshal(x)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(js)
	
	var inventory2 []map[string]string
	if err := json.Unmarshal(js, &inventory2); err != nil {
		log.Fatal(err)
	}
	fmt.Println(inventory2)

	fmt.Println("\nUnmarshal to an existing data")
	if err := json.Unmarshal(js, &inventory); err != nil {
		log.Fatal(err)
	}
	fmt.Println("How did this happen?")
	fmt.Println(inventory)
}
