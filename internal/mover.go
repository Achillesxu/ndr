// Package internal
// Time    : 2022/11/7 21:53
// Author  : xushiyin
// contact : yuqingxushiyin@gmail.com
package internal

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Mover struct {
	srcPath  string
	dstPath  string
	areaCode string
	version  string
	Logger   *log.Entry
}

func NewMover(ctx context.Context, areaCode, ver string, logger *log.Entry) *Mover {
	return &Mover{
		srcPath:  filepath.Clean(viper.GetString("watcher.src_path")),
		dstPath:  filepath.Clean(viper.GetString("watcher.dst_path")),
		areaCode: areaCode,
		version:  ver,
		Logger:   logger,
	}
}

func (m *Mover) LogErr(err error, message string, args ...interface{}) error {
	err = errors.Wrapf(err, message, args...)
	m.Logger.Error(err)
	return err
}

func (m *Mover) CheckConf() error {
	if err := IsFileDirExist(m.srcPath); err != nil {
		return err
	}
	if err := IsFileDirExist(m.dstPath); err != nil {
		return err
	}
	targetFilePath := filepath.Join(m.srcPath, m.areaCode, m.version)
	if err := IsFileDirExist(targetFilePath); err != nil {
		return err
	}
	return nil
}

func (m *Mover) MoveThem() (e error) {
	tmpDir := filepath.Join(m.srcPath, strconv.FormatInt(time.Now().Unix(), 10))
	err := os.Mkdir(tmpDir, 0766)
	if err != nil {
		return m.LogErr(err, "create directory %v failed", tmpDir)
	}
	srcFilePath := filepath.Join(m.srcPath, m.areaCode, m.version, fmt.Sprintf("%s.exe", m.areaCode))
	tmpFilePath := filepath.Join(tmpDir, fmt.Sprintf("%s.exe", m.areaCode))

	err = CopyFile(srcFilePath, tmpFilePath, 0755)
	if err != nil {
		return m.LogErr(err, "")
	}

	// unzip

	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			e = m.LogErr(err, "remove %s", path)
		}
	}(tmpDir)
	return nil
}
