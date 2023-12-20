package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/pluswing/datasync/data"
	"github.com/pluswing/datasync/file"
	"github.com/pluswing/datasync/storage"
	"github.com/spf13/cobra"
)

var pushCmd = &cobra.Command{
	Use:   "push [flags] [version_id]",
	Short: "upload version",
	Long:  `upload version`,
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

		dir, err := file.DataDir()
		cobra.CheckErr(err)

		data.DispatchStorage(setting.Storage, data.StorageFuncTable{
			Gcs: func(conf data.StorageGcsType) {
				storage.Upload(filepath.Join(dir, version.Id+".zip"), version.Id+".zip", conf)
				// FIXME .datasyncを同期したほうが良い。
				file.MoveVersionToRemote(version)
				storage.Upload(filepath.Join(dir, ".datasync"), ".datasync", conf)
			},
		})
		fmt.Printf("push Succeeded. version_id = %s\n", version.Id)
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
}
