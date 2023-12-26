package cmd

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/pluswing/datasync/data"
	"github.com/pluswing/datasync/file"
	"github.com/pluswing/datasync/storage"
	"github.com/spf13/cobra"
)

var all bool

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "list history",
	Long:  `list history`,
	Run: func(cmd *cobra.Command, args []string) {

		dir, err := file.DataDir()
		cobra.CheckErr(err)

		var remote data.DataSyncType
		if all {
			var tmpFile string
			data.DispatchStorage(setting.Storage, data.StorageFuncTable{
				Gcs: func(config data.StorageGcsType) {
					tmpFile = storage.Download(".datasync", config)
				},
			})
			file.MoveFile(tmpFile, filepath.Join(dir, ".datasync"))
			remote = file.ReadRemoteDataSyncFile()
		}

		local := file.ReadLocalDataSyncFile()

		if all {
			fmt.Println("-- remote versions --")
			for _, ver := range remote.Histories {
				d := time.Unix(ver.Time, 0).Format("2006-01-02 15:04:05")
				fmt.Printf("%s %s %s\n", ver.Id, d, ver.Message)
			}
		}
		fmt.Println("-- local versions --")
		for _, ver := range local.Histories {
			d := time.Unix(ver.Time, 0).Format("2006-01-02 15:04:05")
			fmt.Printf("%s %s %s\n", ver.Id, d, ver.Message)
		}
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)

	lsCmd.Flags().BoolVarP(&all, "all", "a", false, "show with remote history")
}
