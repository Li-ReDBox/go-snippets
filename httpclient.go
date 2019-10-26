package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type responseError struct {
}

func main() {
	// without authentication, it only returns error
	// {
	// 	"error": {
	// 	 "errors": [
	// 	  {
	// 	   "domain": "usageLimits",
	// 	   "reason": "dailyLimitExceededUnreg",
	// 	   "message": "Daily Limit for Unauthenticated Use Exceeded. Continued use requires signup.",
	// 	   "extendedHelp": "https://code.google.com/apis/console"
	// 	  }
	// 	 ],
	// 	 "code": 403,
	// 	 "message": "Daily Limit for Unauthenticated Use Exceeded. Continued use requires signup."
	// 	}
	//  }
	uri := "https://www.googleapis.com/customsearch/v1"

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
	fmt.Printf("%v\n", f)

	m := f.(map[string]interface{})
	for k, v := range m {
		switch vv := v.(type) {
		case string:
			fmt.Println(k, "is string", vv)
		case float64:
			fmt.Println(k, "is float64", vv)
		case []interface{}:
			fmt.Println(k, "is an array:")
			for i, u := range vv {
				fmt.Println(i, u)
			}
		case map[string]interface{}:
			fmt.Println(k, "is map", vv)
		default:
			fmt.Println(k, "is of a type I don't know how to handle")
		}
	}
}
