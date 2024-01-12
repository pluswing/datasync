package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "generate datasync.yaml",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("init called")
		// Bubble teaを使って、UIを作る。
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

/*
? How kind of dump target? …
❯ MySQL
  File(s)

-- mysql
? MySQL server hostname / port / username / password / databasename
>

-- file
? Select directory or file
> picker


? Add dump target?
  Yes
❯ No


? Setup remote server?
❯ Yes
  No

? Remote server type?
❯ Google Cloud Storage
  Amazon S3
	Samba

-- GCS
? GCS bucket / path
>

*/
