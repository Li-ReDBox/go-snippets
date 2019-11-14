package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func extractor(s string) {
	fmt.Println(s)
}

func get(url string, ex func(string)) {
	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	robots, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s", robots)
	ex("this is")
}

func main() {
	get("http://www.google.com/robots.txt")
}
