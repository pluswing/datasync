/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/aliakseiz/go-mysqldump"
	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/pluswing/datasync/compress"
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
		fmt.Println("push called", setting)

		if setting.Target.Kind == "mysql" {
			var conf data.TargetMysqlType
			err := mapstructure.Decode(setting.Target.Config, &conf)
			cobra.CheckErr(err)

			dumpDir, err := os.MkdirTemp("", ".datasync")
			cobra.CheckErr(err)

			// dbダンプ ...
			processMysqldump(dumpDir, conf)

			// zip圧縮
			zipFile := compress.Compress(dumpDir)

			// アップロード
			var gcsConf data.UploadGcsType
			err = mapstructure.Decode(setting.Upload.Config, &gcsConf)
			cobra.CheckErr(err)

			_uuid, err := uuid.NewRandom()
			cobra.CheckErr(err)
			uuidStr := _uuid.String()
			uuidStr = strings.Replace(uuidStr, "-", "", -1)

			uploadGcs(zipFile, fmt.Sprintf("%s.zip", uuidStr), gcsConf)

			fmt.Println("DONE upload")

			now := time.Now()

			v := data.VersionType{
				Hash:    uuidStr,
				Time:    now.Unix(),
				Comment: "comment",
			}
			b, err := json.Marshal(v)
			cobra.CheckErr(err)
			version := string(b)

			if existsGcs(".datasync", gcsConf) {
				filePath := downloadGcs(".datasync", gcsConf)
				f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
				cobra.CheckErr(err)
				_, err = f.WriteString(fmt.Sprintf("%s\n", version))
				cobra.CheckErr(err)
				err = f.Close()
				cobra.CheckErr(err)
				uploadGcs(filePath, ".datasync", gcsConf)
			} else {
				tmpDir, err := os.MkdirTemp("", ".datasync")
				cobra.CheckErr(err)
				tmpFile := filepath.Join(tmpDir, ".datasync")
				f, err := os.OpenFile(tmpFile, os.O_CREATE|os.O_WRONLY, 0644)
				cobra.CheckErr(err)
				_, err = f.WriteString(fmt.Sprintf("%s\n", version))
				cobra.CheckErr(err)
				err = f.Close()
				cobra.CheckErr(err)
				uploadGcs(tmpFile, ".datasync", gcsConf)
			}
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

func processMysqldump(dumpDir string, cfg data.TargetMysqlType) {

	config := mysql.NewConfig()
	config.User = cfg.User
	config.Passwd = cfg.Password
	config.DBName = cfg.Database
	config.Net = "tcp"
	config.Addr = fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	dumpFilenameFormat := fmt.Sprintf("%s-%s", "mysql", cfg.Database)

	db, err := sql.Open("mysql", config.FormatDSN())
	cobra.CheckErr(err)

	// Register database with mysqldump.
	dumper, err := mysqldump.Register(db, dumpDir, dumpFilenameFormat, config.DBName)
	cobra.CheckErr(err)

	err = dumper.Dump()
	cobra.CheckErr(err)

	fmt.Println("Successfully mysql dump.")

	dumper.Close()
}

func uploadGcs(target string, fileName string, conf data.UploadGcsType) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	cobra.CheckErr(err)
	defer client.Close()

	f, err := os.Open(target)
	cobra.CheckErr(err)
	defer f.Close()

	// ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	// defer cancel()
	var uploadPath = ""
	if conf.Dir == "" {
		uploadPath = fileName
	} else {
		uploadPath = filepath.Join(conf.Dir, fileName)
	}

	o := client.Bucket(conf.Bucket).Object(uploadPath)

	// o = o.If(storage.Conditions{DoesNotExist: true})

	wc := o.NewWriter(ctx)
	_, err = io.Copy(wc, f)
	cobra.CheckErr(err)

	err = wc.Close()
	cobra.CheckErr(err)
}

func downloadGcs(target string, conf data.UploadGcsType) string {
	// TODO client は一度作ったものを使い回す。
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	cobra.CheckErr(err)
	defer client.Close()

	var filePath = ""
	if conf.Dir == "" {
		filePath = target
	} else {
		filePath = filepath.Join(conf.Dir, target)
	}

	o := client.Bucket(conf.Bucket).Object(filePath)

	tmpDir, err := os.MkdirTemp("", ".datasync")
	cobra.CheckErr(err)

	tmpFile := filepath.Join(tmpDir, target)
	f, err := os.Open(tmpFile)
	cobra.CheckErr(err)
	defer f.Close()

	rc, err := o.NewReader(ctx)
	cobra.CheckErr(err)
	defer rc.Close()

	_, err = io.Copy(f, rc)
	cobra.CheckErr(err)

	return tmpFile
}

func existsGcs(target string, conf data.UploadGcsType) bool {
	// TODO client は一度作ったものを使い回す。
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	cobra.CheckErr(err)
	defer client.Close()

	var filePath = ""
	if conf.Dir == "" {
		filePath = target
	} else {
		filePath = filepath.Join(conf.Dir, target)
	}

	o := client.Bucket(conf.Bucket).Object(filePath)

	attrs, _ := o.Attrs(ctx)

	return attrs != nil
}
