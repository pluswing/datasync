/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/pluswing/datasync/data"
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

		// 読み込む
		data, err := os.ReadFile(tmpFile)
		cobra.CheckErr(err)
		fmt.Println(string(data))
		lines := strings.Split(string(data), "\n")
		var ver data.VersionType
		for _, line := range lines {
			err := json.Unmarshal([]byte(line), &ver)
			cobra.CheckErr(err)
			// いい感じに出力。
			fmt.Printf("")
		}
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}
