package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func parse(f interface{}, guessed string) {
	fmt.Println("Come into an interface, particularly: ", guessed)
	m := f.(map[string]interface{})
	for k, v := range m {
		switch vv := v.(type) {
		case string:
			fmt.Println(k, "is string, content: ", vv)
		case float64:
			fmt.Println(k, "is float64, content: ", vv)
		case []interface{}:
			fmt.Println(k, "is an array, parse its elements")
			for _, u := range vv {
				parse(u, "array")
			}
		case map[string]interface{}:
			fmt.Println(k, "is map, parse its fields")
			parse(vv, "map")
		default:
			fmt.Println(k, "is of a type I don't know how to handle")
		}
	}
	fmt.Printf("%s\n\n", "One interface is done")
}

func main() {
	// without authentication, it only returns error
	// {
	// 	"error": {
	// 		"errors": [
	// 			{
	// 				"domain": "usageLimits",
	// 				"reason": "dailyLimitExceededUnreg",
	// 				"message": "Daily Limit for Unauthenticated Use Exceeded. Continued use requires signup.",
	// 				"extendedHelp": "https://code.google.com/apis/console"
	// 			}
	// 		],
	// 		"code": 403,
	// 		"message": "Daily Limit for Unauthenticated Use Exceeded. Continued use requires signup."
	// 	}
	// }
	uri := "https://www.googleapis.com/customsearch/v1?q=hi"

	resp, err := http.Get(uri)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var f interface{}
	err = json.NewDecoder(resp.Body).Decode(&f)
	if err != nil {
		log.Fatal(err)
	}
	// This line below is very generic
	// fmt.Printf("%v\n", f)

	parse(f, "interface")
}
