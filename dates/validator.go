package dates

import (
	"errors"
	"regexp"
	"strconv"
)

var ErrInvalidateDate = errors.New("invlidate date, it has to be in the format of dd/mm/yyyy between 01/01/1900 - 31/12/2999")

// ParseDateNumbers parses a stirng using dd/mm/yyyy format into an int array.
func ParseDateNumbers(dateStr string) ([]int, error) {
	parts := getParts(dateStr)
	if parts == nil {
		return nil, ErrInvalidateDate
	}
	return validate(parts)
}

// validate checks if digits are valid days, months and years
func validate(parts []string) ([]int, error) {
	// parts are digits
	if len(parts) != 3 {
		return nil, ErrInvalidateDate
	}

	var d, m, y int

	if d, _ = strconv.Atoi(parts[0]); d < 1 || d > 31 {
		return nil, ErrInvalidateDate
	}
	if m, _ = strconv.Atoi(parts[1]); m < 1 || m > 12 {
		return nil, ErrInvalidateDate
	}
	if y, _ = strconv.Atoi(parts[2]); y < 1900 || y > 2999 {
		return nil, ErrInvalidateDate
	}
	return []int{d, m, y}, nil
}

// getParts generates a string array from a string in dd/mm/yyyy format
func getParts(dateStr string) []string {
	// if day or month is single digit, they cannot be zero
	re := regexp.MustCompile(`^(\d{1,2}?)/(\d{1,2}?)/(\d{4})$`)
	parts := re.FindStringSubmatch(dateStr)
	if parts != nil {
		return parts[1:]
	}
	return nil
}
