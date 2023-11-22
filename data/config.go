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

type UploadGcsType struct {
	Bucket string
	Dir    string
}

type UploadType struct {
	Kind   string
	Config interface{}
}

type SettingType struct {
	Target TargetType
	Upload UploadType
}

// ---------------------
type VersionType struct {
	Hash    string
	Time    int64
	Comment string
}
