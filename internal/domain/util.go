package domain

import (
	"encoding/json"
	"time"
)

const day = time.Hour * 24
const week = day * 7

func ParseStruct(obj interface{}) string {
	jsonData, err := json.Marshal(obj)
	if err != nil {
		return ""
	}
	return string(jsonData)
}

// GetLastDayDates return range of dates in UTC between actual day and the previous day
//
//	example: if today is 2023-03-15; then returns (2023-03-14, 2023-03-15)
func GetLastDayDates() (time.Time, time.Time) {
	dateNow := time.Now()
	dateLastDay := dateNow.Add(-1 * day)
	return time.Date(dateLastDay.Year(), dateLastDay.Month(), dateLastDay.Day(), 3, 0, 0, 0, time.UTC),
		time.Date(dateNow.Year(), dateNow.Month(), dateNow.Day(), 3, 0, 0, 0, time.UTC)
}

// GetLastWeekDates return range of dates in UTC between actual day and the previous week
//
//	example: if today is 2023-03-03; then returns (2023-02-24, 2023-03-03)
func GetLastWeekDates() (time.Time, time.Time) {
	dateNow := time.Now()
	dateLastWeek := dateNow.Add(-1 * week)
	return time.Date(dateLastWeek.Year(), dateLastWeek.Month(), dateLastWeek.Day(), 3, 0, 0, 0, time.UTC),
		time.Date(dateNow.Year(), dateNow.Month(), dateNow.Day(), 3, 0, 0, 0, time.UTC)
}
