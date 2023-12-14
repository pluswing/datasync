/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pluswing/datasync/compress"
	"github.com/pluswing/datasync/data"
	"github.com/pluswing/datasync/dump/dump_file"
	"github.com/pluswing/datasync/dump/dump_mysql"
	"github.com/pluswing/datasync/file"
	"github.com/spf13/cobra"
)

var message string

// dumpCmd represents the dump command
var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		dumpDir, err := file.MakeTempDir()
		cobra.CheckErr(err)
		defer os.RemoveAll(dumpDir)

		for _, target := range setting.Targets {
			data.DispatchTarget(target, data.TargetFuncTable{
				Mysql: func(conf data.TargetMysqlType) {
					dump_mysql.Dump(dumpDir, conf)
				},
				File: func(conf data.TargetFileType) {
					dump_file.Dump(dumpDir, conf)
				},
			})
		}

		// zip圧縮
		zipFile := compress.Compress(dumpDir)

		_uuid, err := uuid.NewRandom()
		cobra.CheckErr(err)
		versionId := _uuid.String()
		versionId = strings.Replace(versionId, "-", "", -1)

		// .datasyncに移動
		dir, err := file.DataDir()
		cobra.CheckErr(err)
		dest := filepath.Join(dir, fmt.Sprintf("%s.zip", versionId))
		err = os.Rename(zipFile, dest)
		cobra.CheckErr(err)

		// .datasync/.datasync-local この中がローカルの奴ら。
		// .datasync/.datasync(-remote) これがリモートのやつ。
		// messageを使う。

		newVersion := data.VersionType{
			Id:      versionId,
			Time:    time.Now().Unix(),
			Message: message,
		}

		file.AddHistoryFile(dir, "-local", newVersion)
		file.UpdateVersionFile(versionId)

		fmt.Printf("dump Succeeded. version_id = %s\n", versionId)
	},
}

func init() {
	rootCmd.AddCommand(dumpCmd)

	dumpCmd.Flags().StringVarP(&message, "message", "m", "", "commit mesasge")
	dumpCmd.MarkFlagRequired("message")
}
