// Package excels
// Package internal
// Time    : 2022/9/3 15:26
// Author  : xushiyin
// contact : yuqingxushiyin@gmail.com
package excels

import (
	"errors"
	"fmt"
	"github.com/Achillesxu/ndr/internal"
	_ "github.com/Achillesxu/ndr/internal"
	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
	"os"
	"strconv"
	"strings"
)

type DailyReport struct {
	DateStr     string
	ReportStr   string
	CategoryStr string
	CompleteStr string
	ProgressStr string
	Remarks     string
}

type Excels struct {
	FilePath     string
	Sheet        string
	FilePassword string
	File         *excelize.File
	logger       *log.Entry
}

func NewExcels(filePath, filePassword, sheet string, logger *log.Entry) *Excels {
	return &Excels{
		FilePath:     filePath,
		Sheet:        sheet,
		FilePassword: filePassword,
		logger:       logger.WithFields(log.Fields{"module": "excels"}),
	}
}

func (e *Excels) IsExcelExists() bool {
	if _, err := os.Stat(e.FilePath); errors.Is(err, os.ErrNotExist) {
		e.logger.Errorf("file: %s not found", e.FilePath)
		return false
	} else if err != nil {
		e.logger.Errorf("stat file: %s, err: %v", e.FilePath, err)
		return false
	} else {
		return true
	}
}

func (e *Excels) OpenFile() {
	ops := excelize.Options{}
	if len(e.FilePassword) > 0 {
		ops.Password = e.FilePassword
	}
	var err error

	e.File, err = excelize.OpenFile(e.FilePath, ops)
	if err != nil {
		e.logger.Fatal("Could not open", err)
	}
}

func (e *Excels) WriteDailyReport2Excel(dr *DailyReport) {
	e.logger = e.logger.WithFields(log.Fields{
		"file":  e.FilePath,
		"sheet": e.Sheet,
	})
	e.OpenFile()
	logger := e.logger
	defer func() {
		if err := e.File.Close(); err != nil {
			logger.Fatal("close failed", err)
		}
	}()

	nowMonth := internal.GetNowMonthNumber()
	if !strings.HasPrefix(e.Sheet, nowMonth) {
		logger.Fatalf("now month is %s, you should fix xls.sheet of .ndr.toml", nowMonth)
		return
	}

	idx := e.File.GetSheetIndex(e.Sheet)
	if idx == -1 {
		logger.Fatal("sheet not found")
	}
	e.logger.Infof("sheet idx: %d", idx)
	e.File.SetActiveSheet(idx)
	rCnt := e.FindValidRowNumber()
	logger.Infof("valid row number: %d", rCnt)
	e.WriteDailyReport(rCnt, dr)
	e.MergeDateCell(dr.DateStr)
	if err := e.File.Save(); err != nil {
		logger.Fatalf("save daily report failed: %v", err)
	}
}

// FindValidRowNumber return the number of valid rows in the sheet
func (e *Excels) FindValidRowNumber() int {
	logger := e.logger
	rowCnt := 1
	rows, err := e.File.Rows(e.Sheet)
	if err != nil {
		logger.Fatal("row iterator error", err)
	}

	for rows.Next() {
		row, err := rows.Columns()
		if err != nil {
			logger.Fatal("col iterator error", err)
		}
		if len(row) > 0 && len(row[:4]) > 0 || rowCnt == 2 {
			logger.Infof("col iterator %d, content: %v", rowCnt, row)
			rowCnt += 1
			continue
		} else {
			break
		}
	}
	return rowCnt
}

// WriteDailyReport writes the daily report to the valid row of Excel file
func (e *Excels) WriteDailyReport(rowNum int, dr *DailyReport) {
	rowData := []interface{}{
		dr.DateStr,
		dr.ReportStr,
		dr.CategoryStr,
		dr.CompleteStr,
		dr.ProgressStr,
		dr.Remarks,
	}
	axis := fmt.Sprintf("A%d", rowNum)

	err := e.File.SetSheetRow(e.Sheet, axis, &rowData)
	if err != nil {
		e.logger.Fatal("set sheet row date failed, ", err)
	}
}

func (e *Excels) FindDateCell(date string) []string {
	result, err := e.File.SearchSheet(e.Sheet, date)
	if err != nil {
		e.logger.Fatalf("find date cell %s failed, err: %v", date, err)
	}
	return result
}

func (e *Excels) MergeDateCell(date string) {
	result := e.FindDateCell(date)

	if len(result) > 1 {
		err := e.File.MergeCell(e.Sheet, result[0], result[len(result)-1])
		if err != nil {
			e.logger.Fatalf("merge date cell %s failed, err: %v", date, err)
		}
	}
}

func (e *Excels) GetRowAxis(axis string, line int) []string {
	axisArr := make([]string, 0)
	colName := axis[0]
	rowNum, err := strconv.Atoi(axis[1:])
	if err != nil {
		e.logger.Fatalf("get row axis %s failed, err: %v", axis, err)
	}

	rowNum += line

	for i := 0; i < 6; i++ {
		axisArr = append(axisArr, fmt.Sprintf("%c%v", colName+uint8(i), rowNum))
	}
	return axisArr
}

func (e *Excels) GetRowData(axis []string) []string {
	colStr := make([]string, 0)
	for _, a := range axis {
		v, err := e.File.GetCellValue(e.Sheet, a)
		if err != nil {
			e.logger.Fatalf("get daily report failed, err: %v", err)
		}
		colStr = append(colStr, v)
	}
	return colStr
}

func (e *Excels) GetOneDayDailyReport(axis string, date string) [][]string {
	drData := make([][]string, 0)
	line := 0
	axisArr := e.GetRowAxis(axis, line)
	rowData := e.GetRowData(axisArr)
	drData = append(drData, rowData)
	for {
		line += 1
		axisArr := e.GetRowAxis(axis, line)
		rowData := e.GetRowData(axisArr)
		if rowData[0] != date || (len(rowData[0]) == 0 && len(rowData[1]) == 0) {
			break
		}
		drData = append(drData, rowData)
	}
	return drData
}

func (e *Excels) GetSheetMonth() []string {
	sheetNames := e.File.GetSheetList()
	month := make([]string, 0)
	for _, name := range sheetNames {
		if strings.HasSuffix(name, "月份") {
			month = append(month, strings.Split(name, "月份")[0])
		}
	}
	return month
}

func (e *Excels) ReadOneDayDailyReportFromExcel(date string, raw bool) {
	e.logger = e.logger.WithFields(log.Fields{
		"file":  e.FilePath,
		"sheet": e.Sheet,
	})
	e.OpenFile()
	logger := e.logger
	defer func() {
		if err := e.File.Close(); err != nil {
			logger.Fatal("close failed", err)
		}
	}()

	monthStr := strings.Split(date, "/")[1]

	monthNums := e.GetSheetMonth()

	found := false
	for _, m := range monthNums {
		if m == monthStr {
			found = true
			break
		}
	}

	if !found {
		e.logger.Fatalf("not found sheet for month: %s ", monthStr)
	} else {
		e.Sheet = fmt.Sprintf("%s月份", monthStr)
	}
	result := e.FindDateCell(date)
	if len(result) == 0 {
		logger.Infof("cant find date %s cell", date)
		return
	}
	data := e.GetOneDayDailyReport(result[0], date)

	if raw {
		for _, v := range data {
			fmt.Println(v[1])
		}
	} else {
		table := tablewriter.NewWriter(os.Stdout)
		header := []string{"日期", "工作描述", "工作分类", "今日是否能完成", "工作进度", "备注"}
		table.SetHeader(header)
		for _, d := range data {
			table.Append(d)
		}
		table.Render()
	}
	return
}
