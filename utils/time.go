package utils

import "time"

func GetDiffDays(t1, t2 time.Time) int {
	timeDay1 := time.Date(t1.Year(), t1.Month(), t1.Day(), 0, 0, 0, 0, time.Local)
	timeDay2 := time.Date(t2.Year(), t2.Month(), t2.Day(), 0, 0, 0, 0, time.Local)

	return int(timeDay2.Sub(timeDay1).Hours() / 24)
}
