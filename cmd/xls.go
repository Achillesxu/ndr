// Package cmd
/*
Copyright © 2022 Achilles Xu  <yuqingxushiyin@gmail.com>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	FilePath string
	Password string
	Date     string
	Content  string
)

func init() {
	rootCmd.AddCommand(xlsCmd)
	xlsCmd.AddCommand(inputCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// xlsCmd.PersistentFlags().String("foo", "", "A help for foo")

	xlsCmd.Flags().StringVarP(&FilePath, "file-path", "f", "", "excel file path, for instance: /file/to/i.xls")
	xlsCmd.Flags().StringVarP(&Password, "password", "p", "123", "excel password")
	xlsCmd.Flags().StringVarP(&Date, "date", "d", "today", "default today, date with date format")

	_ = xlsCmd.MarkFlagRequired("file-path")
	_ = xlsCmd.MarkFlagFilename("file-path", "xls", "xlsx")
}

// xlsCmd represents the xls command
var xlsCmd = &cobra.Command{
	Use:   "xls",
	Short: "add daily report to excel file",
	Long:  `add daily report to excel file in directory`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("xls called")
	},
}

var inputCmd = &cobra.Command{
	Use:   "input",
	Short: "input daily report",
	Long:  "input daily report",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("input called")
	},
}
