// Package excels
// Time    : 2022/9/6 01:22
// Author  : xushiyin
// contact : yuqingxushiyin@gmail.com
package excels

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"
)

func TestAxis(t *testing.T) {
	axis := "A9"
	axisArr := make([]string, 6)
	colName := axis[0]
	var rowNum int
	rowNum, _ = strconv.Atoi(axis[1:])
	rowNum += 1
	for i := 0; i < 6; i++ {
		axisArr = append(axisArr, fmt.Sprintf("%c%d", colName+uint8(i), rowNum))
	}
	fmt.Println(axisArr)
}

func TestMatcher(t *testing.T) {
	a := "2022/9/8"
	m, err := regexp.MatchString(`\d{4}/\d{1,2}/\d{1,2}`, a)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(m)
	}
}
