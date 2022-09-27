// Package commands
// Time    : 2022/9/25 23:05
// Author  : xushiyin
// contact : yuqingxushiyin@gmail.com
package commands

import (
	"context"
	"fmt"
	"github.com/Achillesxu/ndr/internal"
	"github.com/go-cmd/cmd"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"strings"
)

func Mount(ctx context.Context, isMount bool) error {
	var name string
	var args []string
	name, args, err := cmdBuild(internal.SysType(), isMount)
	if err != nil {
		return err
	}
	c := cmd.NewCmd(name, args...)
	s := <-c.Start()
	if s.Exit != 0 {
		return errors.Wrapf(s.Error, "failed to run %s %#v, err: %#v", name, args, s.Stderr)
	} else {
		return nil
	}
}

func cmdBuild(sys string, isMount bool) (string, []string, error) {
	var name string
	var args []string

	switch sys {
	case "windows":
		if isMount {
			sharePath := fmt.Sprintf(`\\%s\%s`,
				viper.GetString("smb.host"),
				viper.GetString("smb.path"),
			)
			user := fmt.Sprintf("/user:%s", viper.GetString("smb.username"))
			name, args = "net", []string{"use", viper.GetString("smb.target"), sharePath, user, viper.GetString("smb.password")}
		} else {
			name, args = "net", []string{"use", viper.GetString("smb.target"), "/delete"}
		}
		return name, args, nil
	case "darwin":
		if isMount {
			password := viper.GetString("smb.password")
			// TODO maybe other chars also need to replace
			password = strings.ReplaceAll(password, "@", "%40")

			sharePath := fmt.Sprintf(`//%s:%s@%s/%s`,
				viper.GetString("smb.username"),
				password,
				viper.GetString("smb.host"),
				viper.GetString("smb.path"),
			)
			name, args = "mount", []string{"-t", "smbfs", sharePath, viper.GetString("smb.target")}
		} else {
			name, args = "umount", []string{viper.GetString("smb.target")}
		}
		return name, args, nil
	default:
		return name, args, fmt.Errorf("unsupported system type: %s", sys)
	}
}
