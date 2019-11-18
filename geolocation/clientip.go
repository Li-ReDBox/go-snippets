// Packeg for getting IP of a HTTP client

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
)

const IPGeoLocation = "http://ipgeolocation.com"

// There may be more than one X-Forwarded-For values, the first is the client IP
func getFirstElm(content string) string {
	return strings.TrimSpace(strings.Split(content, ",")[0])
}

// validate an IP string really represents an IP
func parseIP(ip string) string {
	if len(ip) == 0 {
		fmt.Println("Cannot parse an empty string for IP.")
		return ""
	}
	parsed := net.ParseIP(ip)
	if parsed == nil {
		fmt.Printf("%s is not a proper value to be parsed as IP.\n", ip)
		return ""
	}
	return parsed.String()
}

// FromRequest extracts the user IP address from req, if present.
func FromRequest(req *http.Request) string {
	ip := parseIP(getFirstElm(req.Header.Get("X-Forwarded-For")))
	if len(ip) > 8 {
		fmt.Println("Get IP from header X-Forwarded-For.")
		return ip
	}

	ip = parseIP(req.Header.Get("X-Real-IP"))
	if len(ip) > 8 {
		fmt.Println("Get IP from header X-Real-Ip.")
		return ip
	}

	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		fmt.Printf("Cannot split host and port: %s\n", req.RemoteAddr)
		return ""
	}

	// There is always an appeared IP through host
	fmt.Println("Get IP from RemoteAddr")
	return parseIP(ip)
}

func GetGeolocaton(ip string) string {
	resp, err := http.Get(IPGeoLocation + "/" + ip)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return fmt.Sprintf("%s", b)
}

func main() {
	http.Handle("/", http.HandlerFunc(host))

	const addr = ":1718"
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func host(w http.ResponseWriter, r *http.Request) {
	ip := FromRequest(r)
	fmt.Fprintf(w, "Client IP: %v\n", ip)
	fmt.Fprintf(w, "Geo info: %s", GetGeolocaton(ip))
}
