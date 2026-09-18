package utils

import "time"

var WEEK_DAY_STR = []string{"日", "一", "二", "三", "四", "五", "六"}

func LastDayInMonth(year int, month int) int {
	if month > 12 && month <= 0 {
		return -1
	}

	now := time.Now()
	loc := now.Location()

	date := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, loc)
	date = date.AddDate(0, 1, -1)

	return date.Day()
}

func WeekDay(year int, month int, date int) string {
	now := time.Now()
	loc := now.Location()
	t := time.Date(year, time.Month(month), date, 0, 0, 0, 0, loc)
	return WEEK_DAY_STR[t.Weekday()]
}
