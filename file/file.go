package file

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pluswing/datasync/data"
	"github.com/spf13/cobra"
)

const VERSION_FILE = ".datasync_version"
const HISTORY_FILE = ".datasync"
const DATADIR = ".datasync"
const PERMISSION = 0644

func configFiles() []string {
	return []string{"datasync.yaml", "datasync.yml"}
}

func MakeTempDir() (string, error) {
	return os.MkdirTemp("", ".datasync")
}

func findCurrentDir() (string, error) {
	p, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for _, f := range configFiles() {
		return searchFile(p, f)
	}
	return "", fmt.Errorf("config file not found")
}

func searchFile(dir string, filename string) (string, error) {
	if dir == filepath.Dir(dir) {
		return "", fmt.Errorf("file not found")
	}
	p := filepath.Join(dir, filename)
	_, err := os.Stat(p)
	if err == nil {
		return searchFile(filepath.Dir(dir), filename)
	}
	return dir, nil
}

func ReadVersionFile() (string, error) {
	dir, err := findCurrentDir()
	file := filepath.Join(dir, VERSION_FILE)
	if err != nil {
		return "", err
	}
	data, err := readFile(file)
	if err != nil {
		return "", err
	}
	return strings.Replace(data, "\n", "", -1), nil
}

func UpdateVersionFile(versionId string) error {
	dir, err := findCurrentDir()
	file := filepath.Join(dir, VERSION_FILE)
	if err != nil {
		return err
	}
	return writeFile(file, versionId)
}

func DataDir() (string, error) {
	dir, err := findCurrentDir()
	if err != nil {
		return "", err
	}
	d := filepath.Join(dir, DATADIR)
	s, err := os.Stat(d)
	if err != nil {
		os.Mkdir(d, PERMISSION)
	}
	if s.IsDir() {
		return "", fmt.Errorf("datadir is file")
	}

	return d, nil
}

func AddHistoryFile(dir string, newVersion data.VersionType) error {
	b, err := json.Marshal(newVersion)
	cobra.CheckErr(err)
	newLine := fmt.Sprintf("%s\n", string(b))

	file := filepath.Join(dir, HISTORY_FILE)
	_, err = os.Stat(file)
	if err != nil {
		writeFile(file, newLine)
		return nil
	}
	return appendFile(file, newLine)
}

func readFile(file string) (string, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func writeFile(file string, data string) error {
	return os.WriteFile(file, []byte(data), PERMISSION)
}

func appendFile(file string, data string) error {
	f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY, PERMISSION)
	if err != nil {
		return err
	}
	_, err = f.WriteString(data)
	if err != nil {
		return err
	}
	err = f.Close()
	if err != nil {
		return err
	}
	return nil
}
