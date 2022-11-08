// Package files_watcher
// Time    : 2022/11/2 20:18
// Author  : xushiyin
// contact : yuqingxushiyin@gmail.com
package files_watcher

import (
	"context"
	"github.com/Achillesxu/ndr/internal"
	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

type FilesWatcher struct {
	isDaemon       bool
	watcher        *fsnotify.Watcher
	srcPath        string
	srcAreaEntries []os.DirEntry
	dstPath        string
	Logger         *log.Entry
}

func NewFilesWatcher(ctx context.Context, isDaemon bool, logger *log.Entry) (*FilesWatcher, error) {
	watch, err := fsnotify.NewWatcher()
	if err != nil {
		logger.Errorf("error creating watcher: %v", err)
		return nil, err
	}

	return &FilesWatcher{
		isDaemon: isDaemon,
		Logger:   logger,
		watcher:  watch,
		srcPath:  filepath.Clean(viper.GetString("watcher.src_path")),
		dstPath:  filepath.Clean(viper.GetString("watcher.dst_path")),
	}, nil
}

func (fw *FilesWatcher) LogErr(err error, message string, args ...interface{}) error {
	err = errors.Wrapf(err, message, args...)
	fw.Logger.Error(err)
	return err
}

func (fw *FilesWatcher) IsDaemon() bool {
	return fw.isDaemon
}

func (fw *FilesWatcher) CheckDirStructure() error {
	if err := internal.IsFileDirExist(fw.srcPath); err != nil {
		return err
	}
	if err := internal.IsFileDirExist(fw.dstPath); err != nil {
		return err
	}
	return nil
}

func (fw *FilesWatcher) GetAreaSrcDirEntries() error {
	dEntries, err := os.ReadDir(fw.srcPath)
	if err != nil {
		return fw.LogErr(err, "read %s failed", fw.srcPath)
	}

	fw.srcAreaEntries = dEntries
	return nil
}

// AddFiles adds files which are from config to the watcher
func (fw *FilesWatcher) AddFiles() error {
	fw.Logger.Infof("watch dir: %s, move files to %s, if they changed", fw.srcPath, fw.dstPath)
	err := fw.watcher.Add(fw.srcPath)
	if err != nil {
		return fw.LogErr(err, "add files watcher failed")
	}
	return nil
}

func (fw *FilesWatcher) Close() error {
	err := fw.watcher.Close()
	if err != nil {
		return fw.LogErr(err, "close files watcher failed")
	}
	return nil
}
