package data

import (
	"fmt"

	"github.com/mcuadros/go-defaults"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
)

func DispatchTarget(target TargetType, table TargetFuncTable) {
	switch target.Kind {
	case "mysql":
		var conf TargetMysqlType
		err := mapstructure.Decode(target.Config, &conf)
		cobra.CheckErr(err)
		defaults.SetDefaults(&conf)
		table.Mysql(conf)
	case "file":
		var conf TargetFileType
		err := mapstructure.Decode(target.Config, &conf)
		cobra.CheckErr(err)
		defaults.SetDefaults(&conf)
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
		cobra.CheckErr(err)
		defaults.SetDefaults(&conf)
		table.Gcs(conf)
	default:
		panic(fmt.Sprintf("invalid storage.kind = %s\n", storage.Kind))
	}
}
