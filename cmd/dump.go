package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/pluswing/datasync/compress"
	"github.com/pluswing/datasync/data"
	"github.com/pluswing/datasync/dump/dump_file"
	"github.com/pluswing/datasync/dump/dump_mysql"
	"github.com/pluswing/datasync/file"
	"github.com/spf13/cobra"
)

var message string

var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "dump current data",
	Long:  "dump current data",
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

		zipFile := compress.Compress(dumpDir)

		versionId, err := file.NewUUID()
		cobra.CheckErr(err)

		newVersion := data.VersionType{
			Id:      versionId,
			Time:    time.Now().Unix(),
			Message: message,
		}

		dir, err := file.DataDir()
		cobra.CheckErr(err)
		dest := newVersion.FileNameWithDir(dir)
		err = os.Rename(zipFile, dest)
		cobra.CheckErr(err)

		local := file.ReadLocalDataSyncFile()
		local.Histories = append(local.Histories, newVersion)
		file.WriteLocalDataSyncFile(local)
		file.UpdateVersionFile(versionId)

		fmt.Printf("dump Succeeded. version_id = %s\n", versionId)
	},
}

func init() {
	rootCmd.AddCommand(dumpCmd)

	dumpCmd.Flags().StringVarP(&message, "message", "m", "", "commit mesasge")
	dumpCmd.MarkFlagRequired("message")
}
