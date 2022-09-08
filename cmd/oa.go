// Package cmd
/*
Copyright © 2022 Achilles Xu  <yuqingxushiyin@gmail.com>
*/
package cmd

import (
	"fmt"
	"github.com/Achillesxu/ndr/internal"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	isHeadless bool
)

func init() {
	rootCmd.AddCommand(oaCmd)

	oaCmd.AddCommand(dayCmd)
	oaCmd.AddCommand(weekCmd)

	// Here you will define your flags and configuration settings.
	oaCmd.PersistentFlags().BoolVar(&isHeadless, "headless", true,
		"set headless mode, default value is true that means no gui")
}

// oaCmd represents the oa command
var oaCmd = &cobra.Command{
	Use:   "oa",
	Short: "oa open https://oa.jss.com.cn in chrome, and write daily reports or weekly reports",
	Long: `oa open https://oa.jss.com.cn in chrome, chrome can run headless mode, or not,
oa use rod to control chrome via devtools protocol, login your nuoyan account and go to work report page,
submit your daily reports or weekly reports, these reports is from your 每日汇总 xlsx file`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("oa called")
	},
}

var dayCmd = &cobra.Command{
	Use:   "day",
	Short: "submit day reports to the work report page of oa",
	Long: `submit day reports to the work report page of oa,
for instance:
`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := log.WithFields(log.Fields{
			"subCommand": "oa day",
		})
		_, err := internal.NewChromeLogin(cmd.Context(), isHeadless, logger)
		if err != nil {
			logger.Fatal(err)
		}

	},
}

var weekCmd = &cobra.Command{
	Use:   "week",
	Short: "submit week reports to the work report page of oa",
	Long: `submit week reports to the work report page of oa,
for instance:
`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := log.WithFields(log.Fields{
			"subCommand": "oa week",
		})
		_, err := internal.NewChromeLogin(cmd.Context(), isHeadless, logger)
		if err != nil {
			logger.Fatal(err)
		}
	},
}
