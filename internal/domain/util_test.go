package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_ParseStruct_WithCollectionsOrNested_Should_ReturnOk(t *testing.T) {
	anomStruct := struct {
		SomeArr    []string
		SomeMap    map[string]string
		SomeStruct struct{ SomeStr string }
	}{
		SomeArr:    []string{"a", "b"},
		SomeMap:    map[string]string{"a": "1", "b": "2", "c": "3"},
		SomeStruct: struct{ SomeStr string }{SomeStr: "sarasa"},
	}
	assert.Equal(t, `{"SomeArr":["a","b"],"SomeMap":{"a":"1","b":"2","c":"3"},"SomeStruct":{"SomeStr":"sarasa"}}`, ParseStruct(anomStruct))
}

func Test_ParseStruct_WithNoCommonEncodingJSON_Should_ReturnEmptyString(t *testing.T) {
	idString := func(x string) string { return x }
	assert.Equal(t, "", ParseStruct(idString))
}

func Test_GetRangePreviousDaysFromSimpleDateAndDayDuration(t *testing.T) {
	year, month, _day, hour, min, sec, nsec := 2023, time.February, 13, 3, 5, 2, 0
	paramDate := time.Date(year, month, _day, hour, min, sec, nsec, time.UTC)

	date1, date2 := getRangePreviousDaysFrom(paramDate, day)

	expDate1 := time.Date(date1.Year(), date1.Month(), date1.Day(), 3, 0, 0, 0, time.UTC)
	expDate2 := time.Date(paramDate.Year(), paramDate.Month(), paramDate.Day(), 3, 0, 0, 0, time.UTC)
	assert.Equal(t, expDate1, date1)
	assert.Equal(t, expDate2, date2)
	assert.Equal(t, date2.Sub(date1), day)
}

func Test_GetRangePreviousDaysFromConflictDateAndDayDuration(t *testing.T) {
	year, month, _day, hour, min, sec, nsec := 2023, time.March, 1, 0, 0, 0, 0
	paramDate := time.Date(year, month, _day, hour, min, sec, nsec, time.UTC)

	date1, date2 := getRangePreviousDaysFrom(paramDate, day)

	expDate1 := time.Date(date1.Year(), date1.Month(), date1.Day(), 3, 0, 0, 0, time.UTC)
	expDate2 := time.Date(paramDate.Year(), paramDate.Month(), paramDate.Day(), 3, 0, 0, 0, time.UTC)
	assert.Equal(t, expDate1, date1)
	assert.Equal(t, expDate2, date2)
	assert.Equal(t, date2.Sub(date1), day)
}

func Test_GetRangePreviousDaysFromWeekDate(t *testing.T) {
	year, month, _day, hour, min, sec, nsec := 2023, time.February, 28, 3, 5, 2, 0
	paramDate := time.Date(year, month, _day, hour, min, sec, nsec, time.UTC)

	date1, date2 := getRangePreviousDaysFrom(paramDate, week)

	expDate1 := time.Date(date1.Year(), date1.Month(), date1.Day(), 3, 0, 0, 0, time.UTC)
	expDate2 := time.Date(paramDate.Year(), paramDate.Month(), paramDate.Day(), 3, 0, 0, 0, time.UTC)
	assert.Equal(t, expDate1, date1)
	assert.Equal(t, expDate2, date2)
	assert.Equal(t, date2.Sub(date1), week)
}

func Test_GetFollowingDay(t *testing.T) {
	day0 := 15
	day1 := 5
	date0 := time.Date(2021, time.March, day0, 2, 59, 59, 59, time.UTC)
	date1 := time.Date(2023, time.February, day1, 3, 0, 0, 0, time.UTC)
	date2 := time.Date(2022, time.February, 28, 23, 57, 54, 32, time.UTC)

	expectedDate0 := time.Date(date0.Year(), date0.Month(), day0, 3, 0, 0, 0, time.UTC)
	expectedDate1 := time.Date(date1.Year(), date1.Month(), day1+1, 3, 0, 0, 0, time.UTC)
	expectedDate2 := time.Date(date2.Year(), time.March, 1, 3, 0, 0, 0, time.UTC)

	assert.Equal(t, expectedDate0, GetFollowingDay(date0))
	assert.Equal(t, expectedDate1, GetFollowingDay(date1))
	assert.Equal(t, expectedDate2, GetFollowingDay(date2))
}

func Test_GetNextHour(t *testing.T) {
	hour1 := 2
	date1 := time.Date(2023, time.February, 5, hour1, 57, 54, 32, time.UTC)
	date2 := time.Date(2023, time.February, 28, 23, 57, 54, 32, time.UTC)

	expectedDate1 := time.Date(date1.Year(), date1.Month(), date1.Day(), hour1+1, 0, 0, 0, time.UTC)
	expectedDate2 := time.Date(date1.Year(), time.March, 1, 0, 0, 0, 0, time.UTC)

	assert.Equal(t, expectedDate1, GetNextHour(date1))
	assert.Equal(t, expectedDate2, GetNextHour(date2))
}

func Test_ParseOrNilInterface(t *testing.T) {
	var elem1 interface{}
	var elem2 interface{}
	var elem3 interface{}

	timestamp := time.Now()
	elem2 = &Weather{City: "Quilmes", Temperature: 25.5, Timestamp: timestamp}
	elem3 = &Weather{City: "Madrid", Temperature: -1.5, Timestamp: timestamp}
	panicElem := Weather{City: "Quilmes", Temperature: 25.5, Timestamp: timestamp}

	assert.Nil(t, ParseOrNil[Weather](elem1))
	assert.Equal(t, ParseOrNil[Weather](elem2), &Weather{City: "Quilmes", Temperature: 25.5, Timestamp: timestamp})
	assert.Equal(t, ParseOrNil[Weather](elem3), &Weather{City: "Madrid", Temperature: -1.5, Timestamp: timestamp})
	assert.Panics(t, func() { ParseOrNil[Weather](panicElem) })
}
