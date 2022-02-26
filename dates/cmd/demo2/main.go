package main

import (
	"flag"
	"fmt"
	"os"

	"funmech/dates"
	"funmech/dates/demo2"
)

func main() {
	flag.Usage = dates.Usage
	flag.Parse()

	inputs := flag.Args()
	if len(inputs) != 2 {
		dates.Usage()
		os.Exit(1)
	}

	var parsed []demo2.Date

	for i, v := range inputs {
		d, err := demo2.CreateDate(v)

		if err != nil {
			fmt.Printf("Input %d is not valid.\n%s\n", i, err)
		} else {
			parsed = append(parsed, d)

		}
	}

	if len(parsed) != 2 {
		os.Exit(1)
	}

	fmt.Println("Days between", parsed[0], "and", parsed[1], "is", parsed[0].Interval(parsed[1]))
}
