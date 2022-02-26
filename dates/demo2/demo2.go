package demo2

import (
	"fmt"

	"funmech/dates"
)

// A Date reprents a date
type Date struct {
	Year, Month, Day int
}

// CreateDate returns a Date from the given string
func CreateDate(ds string) (Date, error) {
	d, err := dates.ParseDateNumbers(ds)
	if err != nil {
		return Date{}, err
	}
	return Date{d[2], d[1], d[0]}, nil
}

// String is Stringer of Date
func (d Date) String() string {
	return fmt.Sprintf("%02d/%02d/%4d", d.Day, d.Month, d.Year)
}

// toInts returns an int array of the current Date
func (d Date) toInts() []int {
	return []int{d.Day, d.Month, d.Year}
}

// Interval returns the number of days between to Dates
func (d Date) Interval(o Date) int {
	small := d
	big := o
	if d.Year > o.Year || (d.Year == o.Year && d.Month > o.Month) ||
		(d.Year == o.Year && d.Month == o.Month && d.Day > o.Day) {
		small = o
		big = d
	}

	return dates.Diff(small.toInts(), big.toInts())
}
