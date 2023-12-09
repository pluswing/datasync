package data

type TargetMysqlType struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

type TargetFileType struct {
	Path string
}

type TargetType struct {
	Kind   string
	Config interface{}
}

type StorageGcsType struct {
	Bucket string
	Dir    string
}

type StorageType struct {
	Kind   string
	Config interface{}
}

type SettingType struct {
	Targets []TargetType
	Storage StorageType
}

// ---------------------
type VersionType struct {
	Id      string `json:"id"`
	Time    int64  `json:"time"`
	Message string `json:"message"`
}

// ---------------------
// target funcs
type TargetMysqlFunc func(config TargetMysqlType)
type TargetFileFunc func(config TargetFileType)

type TargetFuncTable struct {
	Mysql TargetMysqlFunc
	File  TargetFileFunc
}

// ---------------------
// storage funcs
type StorageGcsFunc func(config StorageGcsType)

type StorageFuncTable struct {
	Gcs StorageGcsFunc
}
