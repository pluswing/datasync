package main

import (
	"github.com/pluswing/datasync/data"
)

func DispatchTarget(target data.TargetType) {
	// var table = {
	// 	"mysql": dump.Dump,
	// 	// "postgres": dump.PgDump,
	// }
	switch target.Kind {
	// case "mysql":
	// 	var conf data.TargetMysqlType
	// 	err := mapstructure.Decode(target.Config, &conf)
	// 	// TODO エラーハンドリング
	// 	table[target.Kind](conf)
	}
}

// TODO 名前よくない
func DispatchUpload(upload data.UploadType) {

}
