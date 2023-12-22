package cmd

import (
	"fmt"
	"os"

	"github.com/pluswing/datasync/compress"
	"github.com/pluswing/datasync/data"
	"github.com/pluswing/datasync/dump/dump_file"
	"github.com/pluswing/datasync/dump/dump_mysql"
	"github.com/pluswing/datasync/file"
	"github.com/spf13/cobra"
)

var applyCmd = &cobra.Command{
	Use:   "apply [flags] [version_id]",
	Short: "apply version",
	Long:  `apply version`,
	Args:  cobra.MatchAll(cobra.RangeArgs(0, 1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {

		versionId, err := file.GetCurrentVersion(args)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		version, err := file.FindVersion(versionId)
		if err != nil {
			fmt.Println("version not found.")
			return
		}

		dir, err := file.DataDir()
		cobra.CheckErr(err)
		tmpFile := version.FileNameWithDir(dir)

		s, err := os.Stat(tmpFile)
		if err != nil || s.IsDir() {
			fmt.Printf("file not found. \nplease run `datasync pull %s`\n", version.Id)
		}

		tmpDir, err := file.MakeTempDir()
		cobra.CheckErr(err)
		defer os.RemoveAll(tmpDir)

		compress.Decompress(tmpDir, tmpFile)

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
}
