package data

type TargetMysqlType struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
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
	Target  TargetType
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

type TargetFuncTable struct {
	Mysql TargetMysqlFunc
}

// ---------------------
// storage funcs
type StorageGcsFunc func(config StorageGcsType)

type StorageFuncTable struct {
	Gcs StorageGcsFunc
}
