/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/pluswing/datasync/compress"
	"github.com/pluswing/datasync/data"
	"github.com/pluswing/datasync/dump/dump_file"
	"github.com/pluswing/datasync/dump/dump_mysql"
	"github.com/pluswing/datasync/file"
	"github.com/pluswing/datasync/storage"
	"github.com/spf13/cobra"
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
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

		data.DispatchStorage(setting.Storage, data.StorageFuncTable{
			Gcs: func(conf data.StorageGcsType) {
				// アップロード
				storage.Upload(zipFile, fmt.Sprintf("%s.zip", versionId), conf)
			},
		})
		// TODO ..datasync_version書き換え
		fmt.Printf("push Succeeded. version_id = %s\n", versionId)
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
}
