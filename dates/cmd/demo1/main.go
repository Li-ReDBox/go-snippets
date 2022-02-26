package main

import (
	"flag"
	"fmt"
	"os"

	"funmech/dates"
	"funmech/dates/demo1"
)

func main() {
	flag.Usage = dates.Usage
	flag.Parse()

	inputs := flag.Args()
	if len(inputs) != 2 {
		dates.Usage()
		os.Exit(1)
	}

	var parsed [][]int

	for i, v := range inputs {
		d, err := dates.ParseDateNumbers(v)

		if err != nil {
			fmt.Printf("Input %d is not valid.\n%s\n", i, err)
		} else {
			parsed = append(parsed, d)

		}
	}

	if len(parsed) != 2 {
		os.Exit(1)
	}

	fmt.Println("Days between", inputs[0], "and", inputs[1], "is", demo1.Interval(parsed[0], parsed[1]))
}
