package utils

import "time"

func TimeTruncate(date time.Time) time.Time {
	y, m, d := date.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

func FirstLastTimeOfMonth(year, month int) (time.Time, time.Time) {
	first := TimeTruncate(time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC))
	last := TimeTruncate(first.AddDate(0, 1, -1))
	return first, last
}
