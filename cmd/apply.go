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

		_, err := file.FindCurrentDir()
		if err != nil {
			fmt.Println("datasync.yaml not found.\nPlease run `datasync init`")
			return
		}

		version, err := file.GetCurrentVersion(args)
		if err != nil {
			fmt.Println(err.Error())
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

		fmt.Println("decompressing...")
		compress.Decompress(tmpDir, tmpFile)

		for _, target := range setting.Targets {
			data.DispatchTarget(target, data.TargetFuncTable{
				Mysql: func(config data.TargetMysqlType) {
					dump_mysql.Import(dump_mysql.MysqlDumpFile(tmpDir, config), config)
				},
				File: func(config data.TargetFileType) {
					dump_file.Expand(tmpDir, config)
				},
			})
		}

		fmt.Println("finalizing...")
		err = file.UpdateVersionFile(version.Id)
		cobra.CheckErr(err)

		fmt.Printf("apply Succeeded. version_id = %s\n", version.Id)
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)
}
