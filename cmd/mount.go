// Package cmd
/*
Copyright © 2022 Achilles Xu  <yuqingxushiyin@gmail.com>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// mountCmd represents the mount command
var mountCmd = &cobra.Command{
	Use:   "mount",
	Short: "mount samba into local path",
	Long: `mount samba share path into local path, or umount samba share path
support Mac, Windows.
Mac command line:
	mount -t smbfs '//user:password@ip/share' /target/share
	umount /target/share
Windows command line:
	net use X: \\ip\share /user:smb password
	net use X: /delete

for instance:
	ndr mount       // mount samba share path into local path
	ndr mount -d    //umount samba share path
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("mount called")
	},
}

func init() {
	rootCmd.AddCommand(mountCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mountCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// mountCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
