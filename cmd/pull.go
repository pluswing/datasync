/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/pluswing/datasync/data"
	"github.com/pluswing/datasync/storage"
	"github.com/spf13/cobra"
)

var id string

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
		fmt.Println("pull called")
		fmt.Println(args)

		// idがあるかどうか。
		var versionId = ""
		if len(args) == 1 {
			versionId = args[0]
		} else {
			versionId = "" // TODO read .datasync_version
		}

		// 指定のバージョンをダウンロード
		if setting.Upload.Kind == "gcs" {
			var gcsConf data.UploadGcsType
			err := mapstructure.Decode(setting.Upload.Config, &gcsConf)
			cobra.CheckErr(err)
			tmpFile := storage.Download(fmt.Sprintf("%s.zip", versionId), gcsConf)
		}
		// storage.Download()

		// 展開する => tmp
		// compress.Decompress(tmpFile)

		// 展開したものを適用する
		// mysql  -> mysql.Import()
		// file   -> copy

		// .datasync_versionを書き換える。
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)
}
