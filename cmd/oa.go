// Package cmd
/*
Copyright © 2022 Achilles Xu  <yuqingxushiyin@gmail.com>
*/
package cmd

import (
	"fmt"
	"github.com/Achillesxu/ndr/internal"
	"github.com/Achillesxu/ndr/internal/excels"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"time"
)

var (
	isHeadless bool
	remoteMode bool
)

func init() {
	rootCmd.AddCommand(oaCmd)

	oaCmd.AddCommand(dayCmd)
	oaCmd.AddCommand(weekCmd)

	// Here you will define your flags and configuration settings.
	oaCmd.PersistentFlags().BoolVar(&isHeadless, "headless", true,
		"set headless mode, default value is true that means no gui")

	oaCmd.PersistentFlags().BoolVar(&remoteMode, "remote", true,
		"using remote mode if it is true, or using local chrome mode, old version chrome maybe have wrongs")

	weekCmd.Flags().IntVarP(&rangeFlag, "range", "r", 5,
		"default: 5, range means that it contains the number of daily report")
}

// oaCmd represents the oa command
var oaCmd = &cobra.Command{
	Use:   "oa",
	Short: "oa open https://oa.jss.com.cn in chrome, and write daily reports or weekly reports",
	Long: `oa open https://oa.jss.com.cn in chrome, chrome can run in remote mode or local mode,
using local mode can select headless mode, oa use rod to control chrome via devtools protocol,
login your nuo yan account and go to work report page,
submit your daily reports, these reports is from your 每日汇总 xlsx file`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("oa called")
	},
}

var dayCmd = &cobra.Command{
	Use:   "day",
	Short: "submit day reports to the working report page of oa",
	Long: `submit day reports to the working report page of oa, date default is today, other date unsupported 
for instance:
ndr oa day // using remote mode
ndr oa day --remote=false --headless=true // using local mode, and headless mode

`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := log.WithFields(log.Fields{
			"subCommand": "oa day",
		})

		logger.Debug("--headless: ", isHeadless)
		logger.Debug("--remote: ", remoteMode)

		reports, err := excels.GetDaysReports(time.Now().Format("2006/1/2"), 1, true, logger)
		if err != nil {
			return
		}

		oa, err := internal.NewOaWebLogin(cmd.Context(), isHeadless, remoteMode, logger)
		if err != nil {
			return
		}
		if err := oa.StuffReport(0, reports, viper.GetStringSlice("oa.copy_to"), true); err != nil {
			return
		}
	},
}

var weekCmd = &cobra.Command{
	Use:   "week",
	Short: "submit week reports to the working report page of oa",
	Long: `submit week reports to the work reporting page of oa,
for instance:
`,
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		logger := log.WithFields(log.Fields{
			"subCommand": "oa week",
		})

		logger.Debug("--headless: ", isHeadless)
		logger.Debug("--remote: ", remoteMode)
		logger.Debug("--range: ", rangeFlag)

		if rangeFlag <= 1 {
			logger.Error("rangeFlag must >= 1")
			return
		}

		reports, err := excels.GetDaysReports(time.Now().Format("2006/1/2"), rangeFlag, true, logger)
		if err != nil {
			return
		}

		oa, err := internal.NewOaWebLogin(cmd.Context(), isHeadless, remoteMode, logger)
		if err != nil {
			logger.Fatal(err)
		}

		if err := oa.StuffReport(1, reports, viper.GetStringSlice("oa.copy_to"), false); err != nil {
			logger.Fatal(err)
		}
	},
}
