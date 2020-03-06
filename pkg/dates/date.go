package dates

import (
	"time"
)

const EndOfDayHour = 17

func LatestBusinessEndOfDay(date time.Time) time.Time {
	switch date.Weekday() {
	case time.Monday:
		if date.Hour() < EndOfDayHour {
			date = date.AddDate(0, 0, -3)
		}
	case time.Sunday:
		date = date.AddDate(0, 0, -2)
	case time.Saturday:
		date = date.AddDate(0, 0, -1)
	default:
		if date.Hour() < EndOfDayHour {
			date = date.AddDate(0, 0, -1)
		}
	}
	year, month, day := date.Date()
	return time.Date(year, month, day, EndOfDayHour, 0, 0, 0, date.Location())
}
