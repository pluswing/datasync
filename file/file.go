package file

import "os"

func MakeTempFile() (string, error) {
	return os.MkdirTemp("", ".datasync")
}
