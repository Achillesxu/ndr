// Package cmd
/*
Copyright © 2022 Achilles Xu  <yuqingxushiyin@gmail.com>
*/
package cmd

import (
	"github.com/Achillesxu/ndr/internal"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	areaCode string
	version  string
)

func init() {
	rootCmd.AddCommand(moveCmd)
	moveCmd.Flags().StringVarP(
		&areaCode, "area_code", "a", "000000",
		"area code length is six, for instance: 230000")
	moveCmd.Flags().StringVarP(
		&version, "version", "v",
		"2300.1.0.0", "taxTool.exe version")
}

// moveCmd represents the move command
var moveCmd = &cobra.Command{
	Use:   "move",
	Short: "move src file to dst directory",
	Long: `mkdir tmp directory, move it to temp directory, unzip it,
and move files to dst directory, finally remove tmp directory;
src_path and dst_path in .ndr.toml in home directory;
for instance:
ndr move --area_code 230000 --version 2300.0.0
`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := log.WithFields(log.Fields{
			"subCommand": "move",
		})
		mover := internal.NewMover(cmd.Context(), areaCode, version, logger)
		if err := mover.CheckConf(); err != nil {
			return
		}
		if err := mover.MoveThem(); err != nil {
			return
		}
		logger.Info("move completed")
	},
}
