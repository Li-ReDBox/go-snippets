package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var servers = [...]string{"https://ipinfo.io/%s", "https://ipgeolocation.com/%s", "https://ipapi.co/%s/json"}

func GetGeolocation(url string, c chan string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Returning from", url)
	c <- fmt.Sprintf("%s", b)
}

func retrieve(ip string) string {
	c := make(chan string, len(servers))
	for _, svr := range servers {
		url := fmt.Sprintf(svr, ip)
		go GetGeolocation(url, c)
	}
	return <-c
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Missing positional argument ip.")
		os.Exit(1)
	}
	fmt.Println(retrieve(os.Args[1]))
}
