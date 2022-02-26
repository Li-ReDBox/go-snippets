package main

import (
	"fmt"
	"time"
)

func main() {
	testDates := [][]time.Time{
		{time.Date(1983, 6, 2, 0, 0, 0, 0, time.UTC), time.Date(1983, 6, 22, 0, 0, 0, 0, time.UTC)},
		{time.Date(1984, 7, 4, 0, 0, 0, 0, time.UTC), time.Date(1984, 12, 25, 0, 0, 0, 0, time.UTC)},
		{time.Date(1983, 8, 3, 0, 0, 0, 0, time.UTC), time.Date(1989, 3, 1, 0, 0, 0, 0, time.UTC)},
		{time.Date(1983, 8, 3, 0, 0, 0, 0, time.UTC), time.Date(1989, 1, 3, 0, 0, 0, 0, time.UTC)},
	}

	for _, ds := range testDates {
		fmt.Printf("%v vs %v\n", ds[0], ds[1])
		difference := ds[1].Sub(ds[0])
		fmt.Printf("Days between them is: %.0f\n\n", difference.Hours()/24-1)
	}
}
