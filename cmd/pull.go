package cmd

import (
	"fmt"

	"github.com/pluswing/datasync/data"
	"github.com/pluswing/datasync/file"
	"github.com/pluswing/datasync/storage"
	"github.com/spf13/cobra"
)

var pullCmd = &cobra.Command{
	Use:   "pull [flags] [version_id]",
	Short: "pull remote version",
	Long:  `pull remote version`,
	Args:  cobra.MatchAll(cobra.RangeArgs(0, 1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {

		version, err := file.GetCurrentVersion(args)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		var tmpFile string
		data.DispatchStorage(setting.Storage, data.StorageFuncTable{
			Gcs: func(config data.StorageGcsType) {
				tmpFile = storage.Download(version.FileName(), config)
			},
		})

		dir, err := file.DataDir()
		cobra.CheckErr(err)

		err = file.MoveFile(tmpFile, version.FileNameWithDir(dir))
		cobra.CheckErr(err)

		fmt.Printf("pull Succeeded. version_id = %s\n", version.Id)
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)
}
