package cmd

import (
	"fmt"

	"github.com/pluswing/datasync/data"
	"github.com/pluswing/datasync/dump/dump_mysql"
	"github.com/pluswing/datasync/file"
	"github.com/spf13/cobra"
)

var database string

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "import mysql dump file",
	Long:  `import mysql dump file`,
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {

		_, err := file.FindCurrentDir()
		if err != nil {
			fmt.Println("datasync.yaml not found.\nPlease run `datasync init`")
			return
		}

		dumpFile := args[0]

		for _, target := range setting.Targets {
			data.DispatchTarget(target, data.TargetFuncTable{
				Mysql: func(config data.TargetMysqlType) {
					if database == config.Database {
						dump_mysql.Import(dumpFile, config)
					}
				},
				File: func(config data.TargetFileType) {
					// no support
				},
			})
		}
	},
}

func init() {
	rootCmd.AddCommand(importCmd)

	importCmd.Flags().StringVarP(&database, "database", "d", "", "database name")
}
