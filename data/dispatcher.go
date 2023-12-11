package data

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
)

func DispatchTarget(target TargetType, table TargetFuncTable) {
	switch target.Kind {
	case "mysql":
		var conf TargetMysqlType
		err := mapstructure.Decode(target.Config, &conf)
		// TODO エラーハンドリング
		cobra.CheckErr(err)
		table.Mysql(conf)
	case "file":
		var conf TargetFileType
		err := mapstructure.Decode(target.Config, &conf)
		// TODO エラーハンドリング
		cobra.CheckErr(err)
		table.File(conf)
	default:
		panic(fmt.Sprintf("invalid target.kind = %s\n", target.Kind))
	}
}

func DispatchStorage(storage StorageType, table StorageFuncTable) {
	switch storage.Kind {
	case "gcs":
		var conf StorageGcsType
		err := mapstructure.Decode(storage.Config, &conf)
		// TODO エラーハンドリング
		cobra.CheckErr(err)
		table.Gcs(conf)
	default:
		panic(fmt.Sprintf("invalid storage.kind = %s\n", storage.Kind))
	}
}
