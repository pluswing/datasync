package cmd

import (
	"fmt"
	"os"
	"path/filepath"

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

		var versionId = ""
		if len(args) == 1 {
			versionId = args[0]
		} else {
			versionId = file.ReadVersionFile()
		}
		if versionId == "" {
			fmt.Println("version not found.")
			return
		}

		version, err := file.FindVersion(versionId)
		if err != nil {
			fmt.Println("version not found.")
			return
		}

		var tmpFile string
		data.DispatchStorage(setting.Storage, data.StorageFuncTable{
			Gcs: func(config data.StorageGcsType) {
				tmpFile = storage.Download(version.Id+".zip", config)
			},
		})

		dir, err := file.DataDir()
		cobra.CheckErr(err)

		err = os.Rename(tmpFile, filepath.Join(dir, version.Id+".zip"))
		cobra.CheckErr(err)

		fmt.Printf("pull Succeeded. version_id = %s\n", version.Id)
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)
}
