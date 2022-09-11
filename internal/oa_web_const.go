// Package internal
// Time    : 2022/9/11 16:24
// Author  : xushiyin
// contact : yuqingxushiyin@gmail.com
package internal

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	log "github.com/sirupsen/logrus"
)

type OaWeb struct {
	IsHeadless  bool
	BrowserPath string
	Logger      *log.Entry
	Browser     *rod.Browser
	Launcher    *launcher.Launcher
	Page        *rod.Page
}

type ReportType uint

const (
	Daily ReportType = iota
	Weekly
)

var (
	CopyToClickXPath = map[ReportType]string{
		Daily:  `//*[@id="z-day-form"]//div[@_search_]`,
		Weekly: `//*[@id="z-week-form"]//div[@_search_]`,
	}
	ReportTypeBtnXPath = map[ReportType]string{
		Daily:  `//a[@href="javascript:void(0);" and string()="日报"]`,
		Weekly: `//a[@href="javascript:void(0);" and string()="周报"]`,
	}
	ReportTextareaXPath = map[ReportType]string{
		Daily:  `//label[contains(text(),'今日完成工作')]/following-sibling::textarea`,
		Weekly: `//label[contains(text(),'本周完成工作')]/following-sibling::textarea`,
	}
	ReportSubmitBtnXPath = map[ReportType]string{
		Daily:  `//*[@id="z-day-form"]/a[string()="提 交"]`,
		Weekly: `//*[@id="z-week-form"]/a[string()="提 交"]`,
	}
)

type captchaReq struct {
	Base64Str string `json:"base64str"`
}

type captchaRes struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}
