// Package excels
// Package internal
// Time    : 2022/9/3 15:26
// Author  : xushiyin
// contact : yuqingxushiyin@gmail.com
package excels

import (
	"fmt"
	"github.com/Achillesxu/ndr/internal"
	_ "github.com/Achillesxu/ndr/internal"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/xuri/excelize/v2"
	"os"
	"path/filepath"
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
	Logger       *log.Entry
}

func NewExcels(filePath, filePassword, sheet string, logger *log.Entry) *Excels {
	return &Excels{
		FilePath:     filePath,
		Sheet:        sheet,
		FilePassword: filePassword,
		Logger:       logger.WithFields(log.Fields{"module": "excels"}),
	}
}

func (e *Excels) LogErr(err error, message string, args ...interface{}) error {
	err = errors.Wrapf(err, message, args...)
	e.Logger.Error(err)
	return err
}

func (e *Excels) IsExcelExists() (bool, error) {
	if _, err := os.Stat(e.FilePath); errors.Is(err, os.ErrNotExist) {
		return false, e.LogErr(err, "file: %s not found", e.FilePath)
	} else if err != nil {
		return false, e.LogErr(err, "stat file: %s, err: %v", e.FilePath, err)
	} else {
		return true, nil
	}
}

func (e *Excels) OpenFile() error {
	ops := excelize.Options{}
	if len(e.FilePassword) > 0 {
		ops.Password = e.FilePassword
	}
	var err error

	e.File, err = excelize.OpenFile(e.FilePath, ops)
	if err != nil {
		return e.LogErr(err, "Could not open %s, err: %v", e.FilePath, err)
	}
	return nil
}

func (e *Excels) WriteDailyReport2Excel(dr *DailyReport) error {
	e.Logger = e.Logger.WithFields(log.Fields{
		"file":  e.FilePath,
		"sheet": e.Sheet,
	})
	err := e.OpenFile()
	if err != nil {
		return err
	}
	logger := e.Logger
	defer func() {
		if err := e.File.Close(); err != nil {
			logger.Fatal("close failed", err)
		}
	}()

	nowMonth := internal.GetNowMonthNumber()
	if !strings.HasPrefix(e.Sheet, nowMonth) {
		// TODO: fix set .ndr.toml
		return e.LogErr(fmt.Errorf("now month is %s, you should fix xls.sheet of .ndr.toml", nowMonth), "")
	}

	idx := e.File.GetSheetIndex(e.Sheet)
	if idx == -1 {
		logger.Fatal("sheet not found")
	}
	e.Logger.Debugf("sheet idx: %d", idx)
	e.File.SetActiveSheet(idx)
	rCnt := e.FindValidRowNumber()
	logger.Debugf("valid row number: %d", rCnt)
	err = e.WriteDailyReport(rCnt, dr)
	if err != nil {
		return err
	}
	err = e.MergeDateCell(dr.DateStr)
	if err != nil {
		return err
	}
	if err := e.File.Save(); err != nil {
		return e.LogErr(err, "save daily report failed")
	}
	return nil
}

// FindValidRowNumber return the number of valid rows in the sheet
func (e *Excels) FindValidRowNumber() int {
	logger := e.Logger
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
			logger.Debugf("col iterator %d, content: %v", rowCnt, row)
			rowCnt += 1
			continue
		} else {
			break
		}
	}
	return rowCnt
}

// WriteDailyReport writes the daily report to the valid row of Excel file
func (e *Excels) WriteDailyReport(rowNum int, dr *DailyReport) error {
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
		return e.LogErr(err, "set sheet row date: %#v failed", dr)
	}
	return nil
}

func (e *Excels) FindDateCell(date string) ([]string, error) {
	result, err := e.File.SearchSheet(e.Sheet, date)
	if err != nil {
		return nil, e.LogErr(err, "find date cell %s failed", date)
	}
	return result, nil
}

func (e *Excels) MergeDateCell(date string) error {
	result, err := e.FindDateCell(date)
	if err != nil {
		return err
	}
	if len(result) > 1 {
		err := e.File.MergeCell(e.Sheet, result[0], result[len(result)-1])
		if err != nil {
			return e.LogErr(err, "merge date cell %s failed", date)
		}
	}
	return nil
}

func (e *Excels) GetRowAxis(axis string, line int) ([]string, error) {
	axisArr := make([]string, 0)
	colName := axis[0]
	rowNum, err := strconv.Atoi(axis[1:])
	if err != nil {
		return nil, e.LogErr(err, "get row axis %s failed", axis)
	}

	rowNum += line

	for i := 0; i < 6; i++ {
		axisArr = append(axisArr, fmt.Sprintf("%c%v", colName+uint8(i), rowNum))
	}
	return axisArr, nil
}

func (e *Excels) GetRowData(axis []string) ([]string, error) {
	colStr := make([]string, 0)
	for _, a := range axis {
		v, err := e.File.GetCellValue(e.Sheet, a)
		if err != nil {
			return nil, e.LogErr(err, "get daily report row data failed, axis: %v", axis)
		}
		colStr = append(colStr, v)
	}
	return colStr, nil
}

func (e *Excels) GetOneDayDailyReport(month string, date string) ([][]string, error) {
	e.Sheet = fmt.Sprintf("%s月份", month)
	result, err := e.FindDateCell(date)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, e.LogErr(fmt.Errorf("no find date %s cell", date), "")
	}
	drData := make([][]string, 0)
	line := 0
	axisArr, err := e.GetRowAxis(result[0], line)
	if err != nil {
		return nil, err
	}
	rowData, err := e.GetRowData(axisArr)
	if err != nil {
		return nil, err
	}
	drData = append(drData, rowData)
	for {
		line += 1
		axisArr, err := e.GetRowAxis(result[0], line)
		if err != nil {
			return nil, err
		}
		rowData, err := e.GetRowData(axisArr)
		if err != nil {
			return nil, err
		}
		if rowData[0] != date || (len(rowData[0]) == 0 && len(rowData[1]) == 0) {
			break
		}
		drData = append(drData, rowData)
	}
	return drData, nil
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

func (e *Excels) RenderData(ds [][][]string, stdOutFlag, onlyContent bool) (string, error) {
	tabStr := &strings.Builder{}
	var tab *tablewriter.Table
	var header []string

	if stdOutFlag {
		tab = tablewriter.NewWriter(os.Stdout)
	} else {
		tab = tablewriter.NewWriter(tabStr)
	}
	if onlyContent {
		header = []string{"工作描述"}
	} else {
		header = []string{"日期", "工作描述", "工作分类", "今日是否能完成", "工作进度", "备注"}
	}

	tab.SetHeader(header)
	tab.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	tab.SetAlignment(tablewriter.ALIGN_CENTER)

	conFilter := func(rows [][]string) [][]string {
		var res [][]string
		for _, v := range rows {
			res = append(res, []string{v[1]})
		}
		return res
	}

	if onlyContent {
		for _, data := range ds {
			tab.AppendBulk(conFilter(data))
		}
	} else {
		for _, data := range ds {
			tab.AppendBulk(data)
		}
	}
	tab.Render()
	if stdOutFlag {
		return "", nil
	} else {
		return tabStr.String(), nil
	}
}

func (e *Excels) ReadOneDayDailyReportFromExcel(dateFlag string, rangeFlag int, stdOutFlag, onlyContent bool) (string, error) {
	e.Logger = e.Logger.WithFields(log.Fields{
		"file":  e.FilePath,
		"sheet": e.Sheet,
	})
	err := e.OpenFile()
	if err != nil {
		return "", err
	}
	logger := e.Logger
	defer func() {
		if err := e.File.Close(); err != nil {
			logger.Error("close failed", err)
		}
	}()

	dates, err := internal.GetDateList(dateFlag, rangeFlag)
	if err != nil {
		return "", e.LogErr(err, "get dates list failed")
	}
	months := internal.GetMonthList(dates)

	monthNums := e.GetSheetMonth()

	monthNumMap := make(map[string]int)
	for i := 0; i < len(monthNums); i++ {
		monthNumMap[monthNums[i]] = 1
	}
	for i := 0; i < len(months); i++ {
		if _, ok := monthNumMap[months[i]]; !ok {
			return "", e.LogErr(fmt.Errorf("not found sheet for month: %s ", months[i]), "")
		}
	}

	datas := make([][][]string, 0)
	for i := len(months) - 1; i >= 0; i-- {
		data, err := e.GetOneDayDailyReport(months[i], dates[i])
		if err != nil {
			return "", e.LogErr(err, "no find sheet %s月份 date %s, err: %v", months[i], dates[i])
		} else {
			datas = append(datas, data)
		}
	}
	return e.RenderData(datas, stdOutFlag, onlyContent)
}

func GetDaysReports(startDate string, rangeCnt int, onlyContent bool, logger *log.Entry) (string, error) {
	xls := NewExcels(
		filepath.Join(viper.GetString("smb.mount_dir"), viper.GetString("xls.path")),
		viper.GetString("xls.password"),
		viper.GetString("xls.sheet"),
		logger,
	)
	if _, err := xls.IsExcelExists(); err != nil {
		return "", err
	}

	reports, err := xls.ReadOneDayDailyReportFromExcel(startDate, rangeCnt, false, onlyContent)
	if err != nil {
		return "", xls.LogErr(err, "")
	}
	if len(reports) <= 0 {
		return "", xls.LogErr(fmt.Errorf("daily report length must > 0"), "")
	}
	return reports, nil
}
