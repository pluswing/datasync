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
	Id      string `json:"id"`
	Time    int64  `json:"time"`
	Message string `json:"message"`
}
