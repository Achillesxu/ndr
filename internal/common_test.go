// Package internal
// Time    : 2022/9/5 20:55
// Author  : xushiyin
// contact : yuqingxushiyin@gmail.com
package internal

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetNowMonth(t *testing.T) {
	monthStr := GetNowMonthNumber()
	require.Equal(t, "9", monthStr)
}

func TestGetDateList(t *testing.T) {
	dates, err := GetDateList("2022/9/5", 5)
	fmt.Println(dates)
	require.NoError(t, err)
}
