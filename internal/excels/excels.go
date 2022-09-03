// Package excels
// Package internal
// Time    : 2022/9/3 15:26
// Author  : xushiyin
// contact : yuqingxushiyin@gmail.com
package excels

import log "github.com/sirupsen/logrus"

type Options struct {
	DateStr     string
	Reports     []string
	CategoryStr string
	CompleteStr string
	ProgressStr string
	Remarks     string
}

type Excels struct {
	Options      *Options
	FilePath     string
	FilePassword string
	log          *log.Logger
}

func NewExcels(options *Options, log *log.Logger) *Excels {
	return &Excels{
		Options: options,
		log:     log,
	}
}
