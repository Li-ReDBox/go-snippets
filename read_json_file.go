package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// read a json file from saved file
func readFromFile(fn string) []shortener.LinkMap {
	content, err := ioutil.ReadFile(fn)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	fmt.Printf("File contents: %s", content)

	var entries []shortener.LinkMap
	err = json.Unmarshal(content, &entries)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return entries
}

func main() {
	backup := "../springboard_runner/uploads_penfolds.json"
	entries := readFromFile(backup)
	if entries == nil {
		fmt.Println("There is no entries read from file %s", backup)
		return
	}

	for _, entry := range entries {
		fmt.Println(entry)
    }
}
