// Package internal
// Time    : 2022/9/5 20:51
// Author  : xushiyin
// contact : yuqingxushiyin@gmail.com
package internal

import (
	"fmt"
	"strings"
	"time"
)

func GetNowMonthNumber() string {
	now := time.Now()
	month := now.Month()
	return fmt.Sprintf("%d", month)
}

func GetDateList(start string, r int) ([]string, error) {
	date, err := CheckDateFormat(start)
	if err != nil {
		return nil, err
	}
	dates := make([]string, 0)
	for i := 0; i < r; i++ {
		dates = append(dates, date.AddDate(0, 0, -i).Format("2006/1/2"))
	}
	return dates, nil
}

func GetMonthList(dates []string) []string {
	months := make([]string, 0, len(dates))
	for _, d := range dates {
		months = append(months, strings.Split(d, "/")[1])
	}
	return months
}

func CheckDateFormat(date string) (*time.Time, error) {
	d, err := time.Parse("2006/1/2", date)
	if err != nil {
		return nil, err
	}
	return &d, nil
}
