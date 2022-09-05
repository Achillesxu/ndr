// Package internal
// Time    : 2022/9/5 20:51
// Author  : xushiyin
// contact : yuqingxushiyin@gmail.com
package internal

import (
	"fmt"
	"time"
)

func GetNowMonthNumber() string {
	now := time.Now()
	month := now.Month()
	return fmt.Sprintf("%d", month)
}
