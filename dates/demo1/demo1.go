package demo1

import (
	"funmech/dates"
)

// Interval returns the number of days between to Dates
func Interval(dateLeft, dateRight []int) int {
	small := dateLeft
	big := dateRight
	if dateLeft[2] > dateRight[2] || (dateLeft[2] == dateRight[2] && dateLeft[1] > dateRight[1]) ||
		(dateLeft[2] == dateRight[2] && dateLeft[1] == dateRight[1] && dateLeft[0] > dateRight[1]) {
		small = dateRight
		big = dateLeft
	}
	return dates.Diff(small, big)
}
