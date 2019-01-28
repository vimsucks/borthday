package util

import (
	"fmt"
	"github.com/lidaobing/chinese_calendar"
	"strconv"
	"strings"
	"time"
)

func DateStr(date *time.Time) string {
	layout := "2006-01-02"
	return date.Format(layout)
}

func MonthDay(date *time.Time) string {
	layout := "01-02"
	return date.Format(layout)
}

func CCStr(cc *chinese_calendar.ChineseCalendar) string {
	return fmt.Sprintf("%d-%02d-%02d", cc.Year, cc.Month, cc.Day)
}

func CCTime(cc *chinese_calendar.ChineseCalendar) time.Time {
	return ParseDate(CCStr(cc))
}

func ParseDate(str string) time.Time {
	layout := "2006-01-02"
	day, _ := time.Parse(layout, str)
	return day
}

// format: 1996-10-03
func SolarToLunar(solarStr string) (*chinese_calendar.ChineseCalendar, error) {
	layout := "2006-01-02"
	day, err := time.Parse(layout, solarStr)
	if err != nil {
		return nil, err
	}
	birthday := chinese_calendar.MustFromTime(day)
	return &birthday, nil
}

// format: 1996-10-03
func LunarToSolar(lunarStr string) (*time.Time, error) {
	digits := strings.Split(lunarStr, "-")
	year, err := strconv.Atoi(digits[0])
	if err != nil {
		return nil, err
	}
	month, err := strconv.Atoi(digits[1])
	if err != nil {
		return nil, err
	}
	day, err := strconv.Atoi(digits[2])
	if err != nil {
		return nil, err
	}
	lunar := chinese_calendar.ChineseCalendar{
		Year: year,
		Month: month,
		Day: day,
	}
	solar := lunar.MustToTime()
	return &solar, nil
}
