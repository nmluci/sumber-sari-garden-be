package timeutil

import (
	"time"
)

func ParseLocalTime(t string, format string) (time.Time, error) {
	loc, err := time.LoadLocation("Asia/Makassar")
	if err != nil {
		return time.Time{}, err
	}
	return time.ParseInLocation(format, t, loc)
}

func FormatLocalTime(t time.Time, format string) string {
	loc, _ := time.LoadLocation("Asia/Makassar")
	return t.In(loc).Format(format)
}

func Localize(t time.Time) time.Time {
	loc, _ := time.LoadLocation("Asia/Makassar")
	return t.In(loc)
}

func FirstDayOfMonth(date time.Time) time.Time {
	return date.AddDate(0, 0, -date.Day()+1)
}

func LastDayOfMonth(date time.Time) time.Time {
	return date.AddDate(0, 1, -date.Day())
}
