// Package cmd
/*
Copyright © 2022 Achilles Xu  <yuqingxushiyin@gmail.com>
*/
package cmd

import (
	"github.com/Achillesxu/ndr/internal"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/thediveo/enumflag/v2"
	"time"

	"github.com/Achillesxu/ndr/internal/excels"
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
	reportFlag   string
	categoryFlag = OtherC
	completeFlag = YeahI
	progressFlag = HundredP
	remarkFlag   string
	rawFlag      bool
	rangeFlag    int
)

func init() {
	rootCmd.AddCommand(xlsCmd)
	xlsCmd.AddCommand(writeCmd)
	xlsCmd.AddCommand(readCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// xlsCmd.PersistentFlags().String("foo", "", "A help for foo")
	writeCmd.Flags().StringVarP(&dateFlag, "date", "d", time.Now().Format("2006/1/2"),
		"default today, date with date format year/month/day，default: 当天, 指定日期2022/9/1")
	writeCmd.Flags().StringVarP(&reportFlag, "report", "r", "1 ",
		`daily report content may contain multiple lines, default：空，例如：
-r "1 天"
`)
	_ = writeCmd.MarkFlagRequired("report")
	writeCmd.Flags().VarP(enumflag.New(&categoryFlag, "category", categoryModeIds, enumflag.EnumCaseSensitive),
		"category", "c",
		"select one 其他，调休，请假、出差，会议，学习提升，技术调研，协助他人，代码优化，运维问题，任务需求, default: 其他",
	)
	// writeCmd.Flags().Lookup("category").NoOptDefVal = "其他"

	writeCmd.Flags().VarP(enumflag.New(&completeFlag, "complete", isCompletedModeIds, enumflag.EnumCaseSensitive),
		"complete", "o",
		"select one 是，否, default: 是",
	)
	// writeCmd.Flags().Lookup("complete").NoOptDefVal = "是"

	writeCmd.Flags().VarP(enumflag.New(&progressFlag, "progress", progressModeIds, enumflag.EnumCaseSensitive),
		"progress", "p",
		"select one 100%，90%，80%，70%，60%，50%，40%，30%，20%，10%，0%，default: 100%",
	)
	// writeCmd.Flags().Lookup("progress").NoOptDefVal = "100%"

	writeCmd.Flags().StringVarP(&remarkFlag, "remark", "m", "", "remark 备注，default: 空")

	readCmd.Flags().StringVarP(&dateFlag, "date", "d", time.Now().Format("2006/1/2"),
		"default today, date with date format yyyy/m/d，default: 某天, 指定日期2022/9/1,or 2022/10/14")

	readCmd.Flags().IntVarP(&rangeFlag, "range", "r", 1,
		"default: 1, only today daily report if range is 1, if range is 5, will contains 5 days daily report")
}

// xlsCmd represents the xls command
var xlsCmd = &cobra.Command{
	Use:   "xls",
	Short: "write/read daily report in excel file",
	Long: `write daily report to target excel file,
read daily report make weekly report`,
}

var writeCmd = &cobra.Command{
	Use:   "write",
	Short: "write one daily report",
	Long: `write one daily report, 
for instance: 
ndr xls write -d 2022/9/5 -r "这是第一条日志" -c 其他 -o 是 -p 100% -m "标记"
ndr xls write -r "这是第二条日志"
ndr xls write -r "这是第二条日志" -o 否 -p 60%
ndr xls write -r "开会" -c 会议
`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := log.WithFields(log.Fields{
			"subCommand": "xls write",
		})
		logger.Debug("--date: ", dateFlag)
		logger.Debug("--report: ", reportFlag)
		logger.Debug("--category: ", categoryModeIds[categoryFlag][0])
		logger.Debug("--complete: ", isCompletedModeIds[completeFlag][0])
		logger.Debug("--progress: ", progressModeIds[progressFlag][0])
		logger.Debug("--remark: ", remarkFlag)

		_, err := internal.CheckDateFormat(dateFlag)
		if err != nil {
			logger.Errorf("dateFlag format must yyyy/m/d, err: %v", err)
			return
		}

		dr := excels.DailyReport{
			DateStr:     dateFlag,
			ReportStr:   reportFlag,
			CategoryStr: categoryModeIds[categoryFlag][0],
			CompleteStr: isCompletedModeIds[completeFlag][0],
			ProgressStr: progressModeIds[progressFlag][0],
			Remarks:     remarkFlag,
		}
		logger.Infof("input flags: %#v", dr)
		xls := excels.NewExcels(
			viper.GetString("xls.path"),
			viper.GetString("xls.password"),
			viper.GetString("xls.sheet"),
			logger,
		)
		if _, err := xls.IsExcelExists(); err != nil {
			return
		}
		_ = xls.WriteDailyReport2Excel(&dr)
	},
}

var readCmd = &cobra.Command{
	Use:   "read",
	Short: "read one day's daily reports, or one week of daily reports",
	Long: `write one daily report,
for instance:
ndr xls read -d 2022/9/5 # one day reports
`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := log.WithFields(log.Fields{
			"subCommand": "xls read",
		})
		logger.Debug("--date: ", dateFlag)
		logger.Debug("--range: ", rangeFlag)

		_, err := internal.CheckDateFormat(dateFlag)
		if err != nil {
			logger.Errorf("dateFlag format must yyyy/m/d, err: %v", err)
			return
		}

		err = validation.Validate(rangeFlag,
			validation.Required, // not empty
			validation.Min(1),   // length between 5 and 100
		)
		if err != nil {
			logger.Errorf("rangeFlag must >= 1, err: %v", err)
			return
		}

		xls := excels.NewExcels(
			viper.GetString("xls.path"),
			viper.GetString("xls.password"),
			viper.GetString("xls.sheet"),
			logger,
		)

		if _, err := xls.IsExcelExists(); err != nil {
			return
		}
		_, _ = xls.ReadOneDayDailyReportFromExcel(dateFlag, rangeFlag, true)

	},
}
