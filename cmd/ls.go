/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/pluswing/datasync/data"
	"github.com/pluswing/datasync/file"
	"github.com/pluswing/datasync/storage"
	"github.com/spf13/cobra"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// .datasync をsotrageから持ってくる
		var tmpFile string
		data.DispatchStorage(setting.Storage, data.StorageFuncTable{
			Gcs: func(config data.StorageGcsType) {
				tmpFile = storage.Download(".datasync", config)
			},
		})

		dir, err := file.DataDir()
		cobra.CheckErr(err)

		err = os.Rename(tmpFile, filepath.Join(dir, ".datasync"))
		cobra.CheckErr(err)

		// 読み込む
		list := file.ListHistory("")

		for _, ver := range list {
			// TODO 出力方法を工夫する
			//  --oneline
			//  デフォルトは git log 的な出力。
			d := time.Unix(ver.Time, 0).Format("2006-01-02 15:04:05")
			fmt.Printf("%s %s %s\n", ver.Id, d, ver.Message)
		}
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}
