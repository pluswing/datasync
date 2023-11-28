/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var id string

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("pull called")

		// idがあるかどうか。

		// ある場合  => id変数をみる
		// ない場合  => .datasync_version

		// 指定のバージョンをダウンロード
		// storage.Download()

		// 展開する => tmp
		// compress.Decompress()

		// 展開したものを適用する
		// mysql  -> mysql.Import()
		// file   -> copy

		// .datasync_versionを書き換える。
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)

	// TODO あとでなおす。
	pullCmd.Flags().StringVar(&id, "id", "", "version hash id")
}
