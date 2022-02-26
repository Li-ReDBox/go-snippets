package dates

import (
	"fmt"
)

var DAYS = []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

// Usage prints help a message
func Usage() {
	fmt.Println("The application accepts exactly two arguments in the format of d[d]/m[m]/yyyy")
}

// DaysSoFar calculates the number of days since January 1 to the date
func DaysSoFar(d, m int) int {
	ds := d
	for i := 0; i < m-1; i++ {
		ds = ds + DAYS[i]
	}
	return ds
}

// LeapYear checks if a given year is a leap year
func LeapYear(y int) bool {
	if y%4 == 0 {
		if y%100 != 0 {
			return true
		}
		if y%400 == 0 {
			return true
		}
	}
	return false
}

// LeapYears calculates the total number of leap years between two years
func LeapYears(ys, ye int) int {
	total := 0
	for i := ys; i <= ye; i++ {
		if LeapYear(i) {
			total += 1
		}
	}
	return total
}

// Diff calculates the how many days are between two dates.
func Diff(dateSmall, dateBig []int) int {
	ds := DaysSoFar(dateSmall[0], dateSmall[1])
	de := DaysSoFar(dateBig[0], dateBig[1])
	var delta int
	// same year
	if dateSmall[2] == dateBig[2] {
		// same dates, the Count is zero
		delta = de - ds
		if delta == 0 {
			return 0
		}
		// check if this involve Feb in a leap year
		if dateSmall[1] <= 2 && dateBig[1] > 2 && LeapYear(dateBig[2]) {
			return delta
		}
		return delta - 1
	}
	lyears := 0
	if dateBig[1] > 2 {
		lyears += LeapYears(dateBig[2], dateBig[2])
	}
	if dateSmall[1] < 2 {
		lyears += LeapYears(dateSmall[2], dateSmall[2])
	}
	if dateBig[2]-dateSmall[2] > 1 {
		lyears += LeapYears(dateSmall[2]+1, dateBig[2]-1)
	}
	ys := dateBig[2] - dateSmall[2]
	return ys*365 + lyears + de - ds - 1
}
