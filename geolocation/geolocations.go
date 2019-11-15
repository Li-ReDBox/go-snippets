package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var servers = [...]string{"https://ipinfo.io/%s", "https://ipgeolocation.com/%s", "https://ipapi.co/%s/json"}

type Coordinate struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

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

func main_main() {
	if len(os.Args) == 1 {
		fmt.Println("Missing positional argument ip.")
		os.Exit(1)
	}
	fmt.Println(retrieve(os.Args[1]))
}

func splitValues(s string) Coordinate {
	fmt.Println(s)
	values := strings.Split(s, ",")
	return Coordinate{Latitude: values[0], Longitude: values[1]}
}

func readCoodinate(jsonStr string) Coordinate {
	var result interface{}
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		fmt.Println(err)
		return Coordinate{}
	}
	m := result.(map[string]interface{})
	if coords, ok := m["coords"]; ok {
		return splitValues(coords.(string))
	} else if loc, ok := m["loc"]; ok {
		return splitValues(loc.(string))
	} else if latitude, ok := m["latitude"]; ok {
		fmt.Println(latitude, m["longitude"])
		return Coordinate{Latitude: fmt.Sprintf("%f", m["latitude"].(float64)), Longitude: fmt.Sprintf("%f", m["longitude"].(float64))}
	}
	return Coordinate{}
}

func main() {
	ex1 := `{"ip":"203.11.249.223","city":"Barangaroo","region":"New South Wales","country":"Australia","coords":"-33.859100,151.200200","asn":"AS17906, PricewaterhouseCoopers","postal":"2000","timezone":"Australia/Sydney"}`
	ex2 := `{"ip":"203.11.249.223","city":"Adelaide","region":"South Australia","country":"AU","loc":"-34.9287,138.5986","org":"AS17906 PricewaterhouseCoopers","postal":"5000","timezone":"Australia/Adelaide","readme":"https://ipinfo.io/missingauth"}`
	ex3 := `{"ip":"203.11.249.223","city":"Barangaroo","region":"New South Wales","region_code":"NSW","country":"AU","country_name":"Australia","continent_code":"OC","in_eu":false,"postal":"2000","latitude":-33.8591,"longitude":151.2002,"timezone":"Australia/Sydney","utc_offset":"+1100","country_calling_code":"+61","currency":"AUD","languages":"en-AU","asn":"AS17906","org":"PricewaterhouseCoopers"}`
	fmt.Println(readCoodinate(ex1))
	fmt.Println(readCoodinate(ex2))
	fmt.Println(readCoodinate(ex3))
}
