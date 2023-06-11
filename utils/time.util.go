package utils

import "time"

// TimeDiff calculates the difference between two time.Time values
// and returns the result in hours, minutes, and seconds.
func TimeDiff(t1, t2 time.Time) (int, int, int) {
	diff := t2.Sub(t1)
	hours := int(diff.Hours())
	minutes := int(diff.Minutes()) % 60
	seconds := int(diff.Seconds()) % 60
	return hours, minutes, seconds
}
