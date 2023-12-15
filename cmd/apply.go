/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pluswing/datasync/compress"
	"github.com/pluswing/datasync/data"
	"github.com/pluswing/datasync/dump/dump_file"
	"github.com/pluswing/datasync/dump/dump_mysql"
	"github.com/pluswing/datasync/file"
	"github.com/spf13/cobra"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply [flags] [version_id]",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MatchAll(cobra.RangeArgs(0, 1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {

		var versionId = ""
		if len(args) == 1 {
			// TODO 先頭6文字くらいでもいけるようにする(git like)
			versionId = args[0]
		} else {
			var err error = nil
			versionId, err = file.ReadVersionFile()
			cobra.CheckErr(err)
		}

		dataDir, err := file.DataDir()
		cobra.CheckErr(err)
		tmpFile := filepath.Join(dataDir, versionId+".zip")

		remoteList := file.ListHistory("")
		localList := file.ListHistory("-local")
		list := append(remoteList, localList...)
		var found = false
		for _, ver := range list {
			if ver.Id == versionId {
				found = true
				break
			}
		}
		if !found {
			fmt.Println("invalid varsion id")
			return
		}
		// TODO tmpFileがあるかどうか

		tmpDir, err := file.MakeTempDir()
		cobra.CheckErr(err)
		defer os.RemoveAll(tmpDir)

		// 展開する => tmp
		compress.Decompress(tmpDir, tmpFile)

		// 展開したものを適用する
		for _, target := range setting.Targets {
			data.DispatchTarget(target, data.TargetFuncTable{
				Mysql: func(config data.TargetMysqlType) {
					dump_mysql.Import(tmpDir, config)
				},
				File: func(config data.TargetFileType) {
					dump_file.Expand(tmpDir, config)
				},
			})
		}

		err = file.UpdateVersionFile(versionId)
		cobra.CheckErr(err)

		fmt.Printf("apply Succeeded. version_id = %s\n", versionId)

	},
}

func init() {
	rootCmd.AddCommand(applyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// applyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// applyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
