/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"database/sql"

	"github.com/JamesStewy/go-mysqldump"
	"github.com/mitchellh/mapstructure"
	"github.com/pluswing/datasync/data"
	"github.com/spf13/cobra"
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("push called", config)

		if config.Target.Kind == "mysql" {
			// config.Target.Config // TargetMysqlConfigTypeに変換する
			var conf data.TargetMysqlConfigType
			err := mapstructure.Decode(config.Target.Config, &conf)
			cobra.CheckErr(err)
			fmt.Println(conf)
			// dbダンプ ...
		}
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pushCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pushCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func processMysqldump(cfg data.TargetMysqlConfigType) string {

	dumpDir, err := os.MkdirTemp("", ".datasync")
	cobra.CheckErr(err)

	dumpFilenameFormat := fmt.Sprintf("%s-20060102T150405", cfg.Database)

	dns := fmt.Sprintf("%s:%s@%s:%d/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
	db, err := sql.Open("mysql", dns)
	cobra.CheckErr(err)

	// Register database with mysqldump.
	dumper, err := mysqldump.Register(db, dumpDir, dumpFilenameFormat)
	cobra.CheckErr(err)

	// Dump database to file.
	resultFilename, err := dumper.Dump()
	cobra.CheckErr(err)

	fmt.Printf("Successiflly mysql dump. %s", resultFilename)

	dumper.Close()

	return resultFilename
}
