package data

type TargetMysqlConfigType struct {
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

type UploadType struct {
	Kind   string
	Config interface{}
}

type ConfigType struct {
	Target TargetType
	Upload UploadType
}
