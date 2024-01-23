package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pluswing/datasync/data"
	"github.com/pluswing/datasync/file"
	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "remove dump",
	Long:  `remove dump`,
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		_, err := file.FindCurrentDir()
		if err != nil {
			fmt.Println("datasync.yaml not found.\nPlease run `datasync init`")
			return
		}

		versionId := args[0]

		ds := file.ReadLocalDataSyncFile()

		newHistories := []data.VersionType{}
		for _, v := range ds.Histories {
			if v.Id != versionId {
				newHistories = append(newHistories, v)
			}
		}
		if len(ds.Histories) == len(newHistories) {
			fmt.Printf("version not found. %s\n", versionId)
			return
		}
		ds.Histories = newHistories
		file.WriteLocalDataSyncFile(ds)

		dir, err := file.DataDir()
		cobra.CheckErr(err)
		os.Remove(filepath.Join(dir, versionId))

		fmt.Printf("removed. %s\n", versionId)
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
}
