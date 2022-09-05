// Package internal
// Time    : 2022/9/5 20:55
// Author  : xushiyin
// contact : yuqingxushiyin@gmail.com
package internal

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetNowMonth(t *testing.T) {
	monthStr := GetNowMonthNumber()
	require.Equal(t, "9", monthStr)
}
