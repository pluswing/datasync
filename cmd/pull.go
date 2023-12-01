/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/pluswing/datasync/compress"
	"github.com/pluswing/datasync/data"
	"github.com/pluswing/datasync/dump"
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
			// TODO .datasync_versionがない場合の考慮
			f, err := os.Open(".datasync_version")
			cobra.CheckErr(err)
			data := make([]byte, 1024)
			_, err = f.Read(data)
			cobra.CheckErr(err)
			versionId = strings.Replace(string(data), "\n", "", -1)
			err = f.Close()
			cobra.CheckErr(err)
		}

		// 指定のバージョンをダウンロード
		var tmpFile string
		data.DispatchStorage(setting.Storage, data.StorageFuncTable{
			Gcs: func(config data.StorageGcsType) {
				tmpFile = storage.Download(fmt.Sprintf("%s.zip", versionId), config)
			},
		})

		tmpDir, err := os.MkdirTemp("", ".datasync")
		cobra.CheckErr(err)

		// 展開する => tmp
		compress.Decompress(tmpDir, tmpFile)

		// 展開したものを適用する
		data.DispatchTarget(setting.Target, data.TargetFuncTable{
			Mysql: func(config data.TargetMysqlType) {
				dump.Import(tmpDir, config)
			},
		})

		// .datasync_versionを書き換える。
		f, err := os.Open(".datasync_version")
		cobra.CheckErr(err)
		defer f.Close()
		f.WriteString(versionId)

		fmt.Printf("pull Succeeded. version_id = %s\n", versionId)
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)
}
