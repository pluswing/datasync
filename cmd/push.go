/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"cloud.google.com/go/storage"
	"github.com/aliakseiz/go-mysqldump"
	"github.com/go-sql-driver/mysql"
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
			dumpfile := processMysqldump(conf)
			fmt.Println(dumpfile)

			// アップロード
			uploadGoogleStorage()

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

	config := mysql.NewConfig()
	config.User = cfg.User
	config.Passwd = cfg.Password
	config.DBName = cfg.Database
	config.Net = "tcp"
	config.Addr = fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	dumpDir, err := os.MkdirTemp("", ".datasync")
	cobra.CheckErr(err)

	dumpFilenameFormat := fmt.Sprintf("%s-20060102T150405", cfg.Database)

	db, err := sql.Open("mysql", config.FormatDSN())
	cobra.CheckErr(err)

	// Register database with mysqldump.
	dumper, err := mysqldump.Register(db, dumpDir, dumpFilenameFormat, config.DBName)
	cobra.CheckErr(err)

	// FIXME ダンプファイル名がとれないよ。
	// Dump database to file.
	err = dumper.Dump()
	cobra.CheckErr(err)

	fmt.Printf("Successiflly mysql dump. %s\n", dumpFilenameFormat)

	dumper.Close()

	return filepath.Join(dumpDir, dumpFilenameFormat+".sql")
}

func uploadGoogleStorage() {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	cobra.CheckErr(err)
	defer client.Close()

	// Open local file.
	f, err := os.Open("README.md")
	cobra.CheckErr(err)
	defer f.Close()

	// ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	// defer cancel()

	o := client.Bucket("datasync000001").Object("README.md")

	// o = o.If(storage.Conditions{DoesNotExist: true})

	wc := o.NewWriter(ctx)
	_, err = io.Copy(wc, f)
	cobra.CheckErr(err)

	err = wc.Close()
	cobra.CheckErr(err)
}
