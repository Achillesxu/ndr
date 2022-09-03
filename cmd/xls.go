// Package cmd
/*
Copyright © 2022 Achilles Xu  <yuqingxushiyin@gmail.com>
*/
package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/thediveo/enumflag/v2"
	"time"

	"github.com/spf13/cobra"
)

type CategoryMode enumflag.Flag
type IsCompletedMode enumflag.Flag
type ProgressMode enumflag.Flag

// CategoryMode 其他，调休，请假、出差，会议，学习提升，技术调研，协助他人，代码优化，运维问题，任务需求
const (
	OtherC CategoryMode = iota // always colorize
	ChangeHolidayC
	LeaveAbsenceC
	BusinessTripC
	MeetingC
	LearningC
	TechnicalResearchC
	AssistingOthersC
	CodeOptimisationC
	MaintenanceIssuesC
	TaskRequirementC
)

// IsCompletedMode 是，否
const (
	YeahI IsCompletedMode = iota + 1
	NopeI
)

// ProgressMode 100%，，0%
const (
	HundredP ProgressMode = iota
	NinetyP
	EightyP
	SeventyP
	SixtyP
	FiftyP
	FortyP
	ThirtyP
	TwentyP
	TenP
	ZeroP
)

// categoryModeIds Defines the textual representations for the CategoryMode values.
var categoryModeIds = map[CategoryMode][]string{
	OtherC:             {"其他"},
	ChangeHolidayC:     {"调休"},
	LeaveAbsenceC:      {"请假"},
	BusinessTripC:      {"出差"},
	MeetingC:           {"会议"},
	LearningC:          {"学习提升"},
	TechnicalResearchC: {"技术调研"},
	AssistingOthersC:   {"协助他人"},
	CodeOptimisationC:  {"代码优化"},
	MaintenanceIssuesC: {"运维问题"},
	TaskRequirementC:   {"任务需求"},
}

var isCompletedModeIds = map[IsCompletedMode][]string{
	YeahI: {"是"},
	NopeI: {"否"},
}

var progressModeIds = map[ProgressMode][]string{
	HundredP: {"100%"},
	NinetyP:  {"90%"},
	EightyP:  {"80%"},
	SeventyP: {"70%"},
	SixtyP:   {"60%"},
	FiftyP:   {"50%"},
	FortyP:   {"40%"},
	ThirtyP:  {"30%"},
	TwentyP:  {"20%"},
	TenP:     {"10%"},
	ZeroP:    {"0%"},
}

var (
	dateFlag     string
	reportFlag   []string
	categoryFlag = OtherC
	completeFlag = YeahI
	progressFlag = HundredP
	remarkFlag   string
)

func init() {
	rootCmd.AddCommand(xlsCmd)
	xlsCmd.AddCommand(inputCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// xlsCmd.PersistentFlags().String("foo", "", "A help for foo")
	inputCmd.Flags().StringVarP(&dateFlag, "date", "d", time.Now().Format("01/02/2006"),
		"default today, date with date format day/month/year，default: 当天, 指定日期29/09/2022")
	inputCmd.Flags().StringSliceVarP(&reportFlag, "report", "r", []string{},
		`daily report content may contain multiple lines, default：空，例如：
-r "1 天" "2 本" "3 本周"
`)

	inputCmd.Flags().VarP(enumflag.New(&categoryFlag, "category", categoryModeIds, enumflag.EnumCaseSensitive),
		"category", "c",
		"select one 其他，调休，请假、出差，会议，学习提升，技术调研，协助他人，代码优化，运维问题，任务需求, default: 其他",
	)
	// inputCmd.Flags().Lookup("category").NoOptDefVal = "其他"

	inputCmd.Flags().VarP(enumflag.New(&completeFlag, "complete", isCompletedModeIds, enumflag.EnumCaseSensitive),
		"complete", "o",
		"select one 是，否, default: 是",
	)
	// inputCmd.Flags().Lookup("complete").NoOptDefVal = "是"

	inputCmd.Flags().VarP(enumflag.New(&progressFlag, "progress", progressModeIds, enumflag.EnumCaseSensitive),
		"progress", "p",
		"select one 100%，90%，80%，70%，60%，50%，40%，30%，20%，10%，0%，default: 100%",
	)
	// inputCmd.Flags().Lookup("progress").NoOptDefVal = "100%"

	inputCmd.Flags().StringVarP(&remarkFlag, "remark", "m", "", "remark 备注，default: 空")
}

// xlsCmd represents the xls command
var xlsCmd = &cobra.Command{
	Use:   "xls",
	Short: "write/read daily report in excel file",
	Long: `write daily report to target excel file,
read daily report make weekly report`,
}

var inputCmd = &cobra.Command{
	Use:   "input",
	Short: "input daily report",
	Long:  "input daily report",
	Run: func(cmd *cobra.Command, args []string) {
		log := log.WithFields(log.Fields{
			"subCommand": "xls input",
		})
		log.Info("dateFlag: ", dateFlag)
		log.Info("reportFlag: ", reportFlag)
		log.Info("category: ", categoryModeIds[categoryFlag][0])
		log.Info("completeFlag: ", isCompletedModeIds[completeFlag][0])
		log.Info("progressFlag: ", progressModeIds[progressFlag][0])
		log.Info("remarkFlag: ", remarkFlag)
	},
}
