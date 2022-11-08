// Package cmd
/*
Copyright © 2022 Achilles Xu  <yuqingxushiyin@gmail.com>
*/
package cmd

import (
	"github.com/Achillesxu/ndr/internal/files_watcher"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var (
	isDaemon bool
)

func init() {
	rootCmd.AddCommand(watcherCmd)

	// Here you will define your flags and configuration settings.

	watcherCmd.Flags().BoolVar(&isDaemon, "daemon", false,
		"ndr watcher --daemon=true, command will run in daemon mode")
}

// watcherCmd represents the watcher command
var watcherCmd = &cobra.Command{
	Use:   "watcher",
	Short: "watch files and move them to target location",
	Long:  `watch files in config, move them to target location if them changed`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := log.WithFields(log.Fields{
			"subCommand": "watcher",
		})
		logger.Info("start files watcher")
		fw, err := files_watcher.NewFilesWatcher(cmd.Context(), false, logger)
		if err != nil {
			logger.Error(err)
		}
		if err := fw.CheckDirStructure(); err != nil {
			return
		}

		_ = fw.Close()
	},
}
