package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

var servers = [...]string{"https://ipinfo.io/%s", "https://ipgeolocation.com/%s", "https://ipapi.co/%s/json"}

type Coordinate struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

func splitValues(s string) Coordinate {
	values := strings.Split(s, ",")
	return Coordinate{Latitude: values[0], Longitude: values[1]}
}

func readCoodinate(body io.Reader) Coordinate {
	var result interface{}
	err := json.NewDecoder(body).Decode(&result)
	if err != nil {
		fmt.Println(err)
		return Coordinate{}
	}
	m := result.(map[string]interface{})
	if coords, ok := m["coords"]; ok {
		return splitValues(coords.(string))
	} else if loc, ok := m["loc"]; ok {
		return splitValues(loc.(string))
	} else if _, ok := m["latitude"]; ok {
		return Coordinate{Latitude: fmt.Sprintf("%f", m["latitude"].(float64)), Longitude: fmt.Sprintf("%f", m["longitude"].(float64))}
	}
	return Coordinate{}
}

func query(url string, c chan Coordinate) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Returning from", url)
	c <- readCoodinate(resp.Body)
}

// Get geolocation of an IP
func GetGeolocation(ip string) Coordinate {
	c := make(chan Coordinate, len(servers))
	for _, svr := range servers {
		url := fmt.Sprintf(svr, ip)
		go query(url, c)
	}
	return <-c
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Missing positional argument ip.")
		os.Exit(1)
	}
	fmt.Println(GetGeolocation(os.Args[1]))
}
