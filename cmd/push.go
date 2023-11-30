/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pluswing/datasync/compress"
	"github.com/pluswing/datasync/data"
	"github.com/pluswing/datasync/dump"
	"github.com/pluswing/datasync/storage"
	"github.com/spf13/cobra"
)

var message string

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

		dumpDir, err := os.MkdirTemp("", ".datasync")
		cobra.CheckErr(err)

		data.DispatchTarget(setting.Target, data.TargetFuncTable{
			Mysql: func(conf data.TargetMysqlType) {
				dump.Dump(dumpDir, conf)
			},
		})

		// zip圧縮
		zipFile := compress.Compress(dumpDir)

		_uuid, err := uuid.NewRandom()
		cobra.CheckErr(err)
		uuidStr := _uuid.String()
		uuidStr = strings.Replace(uuidStr, "-", "", -1)

		data.DispatchStorage(setting.Storage, data.StorageFuncTable{
			Gcs: func(conf data.StorageGcsType) {
				// アップロード
				storage.Upload(zipFile, fmt.Sprintf("%s.zip", uuidStr), conf)

				fmt.Println("DONE upload")

				now := time.Now()

				v := data.VersionType{
					Id:      uuidStr,
					Time:    now.Unix(),
					Message: message,
				}
				b, err := json.Marshal(v)
				cobra.CheckErr(err)
				version := string(b)

				// TODO 抽象化
				if storage.Exists(".datasync", conf) {
					filePath := storage.Download(".datasync", conf)
					f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
					cobra.CheckErr(err)
					_, err = f.WriteString(fmt.Sprintf("%s\n", version))
					cobra.CheckErr(err)
					err = f.Close()
					cobra.CheckErr(err)
					storage.Upload(filePath, ".datasync", conf)
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
					storage.Upload(tmpFile, ".datasync", conf)
				}
			},
		})
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)

	pushCmd.Flags().StringVarP(&message, "message", "m", "", "commit mesasge")
	pushCmd.MarkPersistentFlagRequired("message")
}
