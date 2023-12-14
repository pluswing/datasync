/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pluswing/datasync/data"
	"github.com/pluswing/datasync/file"
	"github.com/pluswing/datasync/storage"
	"github.com/spf13/cobra"
)

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pull [flags] [version_id]",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MatchAll(cobra.RangeArgs(0, 1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {

		// idがあるかどうか。
		var versionId = ""
		if len(args) == 1 {
			// TODO 先頭6文字くらいでもいけるようにする(git like)
			versionId = args[0]
		} else {
			var err error
			versionId, err = file.ReadVersionFile()
			cobra.CheckErr(err)
		}

		// 指定のバージョンをダウンロード
		var tmpFile string
		data.DispatchStorage(setting.Storage, data.StorageFuncTable{
			Gcs: func(config data.StorageGcsType) {
				tmpFile = storage.Download(versionId+".zip", config)
			},
		})

		dir, err := file.DataDir()
		cobra.CheckErr(err)

		err = os.Rename(tmpFile, filepath.Join(dir, versionId+".zip"))
		cobra.CheckErr(err)

		// TODO リモートの.datasyncを持ってくる

		fmt.Printf("pull Succeeded. version_id = %s\n", versionId)
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)
}
