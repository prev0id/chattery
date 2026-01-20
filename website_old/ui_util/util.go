package ui_util

import "time"

const (
	oneDay = 24 * time.Hour

	timeFormatHourAndMinite             = "15:04"
	timeFormatMonthDayHourAndMinite     = "Mon, Jan _2 15:04"
	timeFormatMonthDayYearHourAndMinite = "Mon, Jan _2 2006 15:04"
)

func GetProp[T any](props []T) T {
	if len(props) > 0 {
		return props[0]
	}
	var defaultValue T
	return defaultValue
}

func PrettyTime(t time.Time) string {
	now := time.Now()
	if t.Add(oneDay).After(now) {
		return t.Format(timeFormatHourAndMinite)
	}
	if t.Year() != now.Year() {
		return t.Format(timeFormatMonthDayYearHourAndMinite)
	}
	return t.Format(timeFormatMonthDayHourAndMinite)
}
