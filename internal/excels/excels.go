// Package excels
// Package internal
// Time    : 2022/9/3 15:26
// Author  : xushiyin
// contact : yuqingxushiyin@gmail.com
package excels

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/xuri/excelize/v2"
	"os"
)

type DailyReport struct {
	DateStr     string
	Reports     string
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

func (e *Excels) WriteDailyReport2Excel() {
	ops := excelize.Options{}
	if len(e.FilePassword) > 0 {
		ops.Password = e.FilePassword
	}
	logger := e.logger.WithFields(log.Fields{
		"file":  e.FilePath,
		"sheet": e.Sheet,
	})
	var err error

	e.File, err = excelize.OpenFile(viper.GetString("xls.path"), ops)
	if err != nil {
		logger.Error("Could not open", err)
		return
	}
	defer func() {
		if err := e.File.Close(); err != nil {
			logger.Error("close failed", err)
		}
	}()
	idx := e.File.GetSheetIndex(viper.GetString("xls.sheet"))
	if idx == -1 {
		logger.Error("sheet not found")
	}
	e.logger.Infof("sheet idx: %d", idx)
	e.File.SetActiveSheet(idx)

}
