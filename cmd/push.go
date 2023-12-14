/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"path/filepath"

	"github.com/pluswing/datasync/data"
	"github.com/pluswing/datasync/file"
	"github.com/pluswing/datasync/storage"
	"github.com/spf13/cobra"
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push pull [flags] [version_id ...]",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MatchAll(cobra.MinimumNArgs(0), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {

		dir, err := file.DataDir()
		cobra.CheckErr(err)

		// args == version_idの配列
		target := make([]string, len(args))

		for i, versionId := range args {
			// TODO 先頭6文字くらいでもいけるようにする(git like)
			target[i] = versionId
		}

		data.DispatchStorage(setting.Storage, data.StorageFuncTable{
			Gcs: func(conf data.StorageGcsType) {
				for _, versionId := range target {
					storage.Upload(filepath.Join(dir, versionId+".zip"), versionId+".zip", conf)
				}
				// TODO .datasync-localのデータをリモートの.datasyncに同期する
			},
		})
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
}
