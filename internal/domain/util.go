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
	return getRangePreviousDaysFrom(time.Now().UTC(), day)
}

// GetLastWeekDates return range of dates in UTC between actual day and the previous week
//
//	example: if today is 2023-03-03; then returns (2023-02-24, 2023-03-03)
func GetLastWeekDates() (time.Time, time.Time) {
	return getRangePreviousDaysFrom(time.Now().UTC(), week)
}

func getRangePreviousDaysFrom(date time.Time, duration time.Duration) (time.Time, time.Time) {
	dateLastDuration := date.Add(-1 * duration)
	return time.Date(dateLastDuration.Year(), dateLastDuration.Month(), dateLastDuration.Day(), 3, 0, 0, 0, time.UTC),
		time.Date(date.Year(), date.Month(), date.Day(), 3, 0, 0, 0, time.UTC)
}

// GetFollowingDay returns the following day at 03:00hs - UTC (00:00hs - GMT-3)
//
//	example 1: if date is "2023-02-24T22:30:59"; then returns "2023-02-25T03:00:00"
//	example 2: if date is "2023-02-25T02:30:59"; then returns "2023-02-25T03:00:00"
func GetFollowingDay(date time.Time) time.Time {
	dateUTC := date.UTC()
	followingDate := dateUTC
	// if time is in 00:00hs and 02:59hs its same day
	if dateUTC.Hour() >= 3 {
		followingDate = dateUTC.Add(day)
	}
	return time.Date(followingDate.Year(), followingDate.Month(), followingDate.Day(), 3, 0, 0, 0, time.UTC)
}

// GetNextHour returns the date with next hour at 0 minute and second in UTC
//
//	example: if date is "2023-02-24T22:30:59"; then returns "2023-02-24T23:00:00"
func GetNextHour(date time.Time) time.Time {
	dateUTC := date.UTC()
	nextHourFromDateNow := dateUTC.Add(time.Hour)
	return time.Date(nextHourFromDateNow.Year(), nextHourFromDateNow.Month(), nextHourFromDateNow.Day(), nextHourFromDateNow.Hour(), 0, 0, 0, time.UTC)
}
