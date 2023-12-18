package file

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pluswing/datasync/data"
	"github.com/spf13/cobra"
)

const VERSION_FILE = ".datasync_version"
const HISTORY_FILE = ".datasync"
const DATADIR = ".datasync"

func configFiles() []string {
	return []string{"datasync.yaml", "datasync.yml"}
}

func MakeTempDir() (string, error) {
	return os.MkdirTemp("", ".datasync")
}

func FindCurrentDir() (string, error) {
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
	if err != nil {
		return searchFile(filepath.Dir(dir), filename)
	}
	return dir, nil
}

func ReadVersionFile() string {
	dir, err := FindCurrentDir()
	if err != nil {
		return ""
	}
	file := filepath.Join(dir, VERSION_FILE)
	data, err := readFile(file)
	if err != nil {
		return ""
	}
	return strings.Replace(data, "\n", "", -1)
}

func UpdateVersionFile(versionId string) error {
	dir, err := FindCurrentDir()
	file := filepath.Join(dir, VERSION_FILE)
	if err != nil {
		return err
	}
	return writeFile(file, versionId)
}

func DataDir() (string, error) {
	dir, err := FindCurrentDir()
	if err != nil {
		return "", err
	}
	d := filepath.Join(dir, DATADIR)
	s, err := os.Stat(d)
	if err != nil {
		os.Mkdir(d, os.ModePerm)
	} else if !s.IsDir() {
		return "", fmt.Errorf("datadir is file")
	}
	return d, nil
}

func AddHistoryFile(dir string, suffix string, newVersion data.VersionType) error {
	newLine, err := versionToString(newVersion)
	cobra.CheckErr(err)
	file := filepath.Join(dir, HISTORY_FILE+suffix)
	_, err = os.Stat(file)
	if err != nil {
		writeFile(file, newLine)
		return nil
	}
	return appendFile(file, newLine)
}

func versionToString(version data.VersionType) (string, error) {
	b, err := json.Marshal(version)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s\n", string(b)), nil
}

func ListHistory(suffix string) []data.VersionType {
	dir, err := DataDir()
	cobra.CheckErr(err)
	file := filepath.Join(dir, HISTORY_FILE+suffix)
	content, err := readFile(file)
	if err != nil {
		return []data.VersionType{}
	}
	lines := strings.Split(content, "\n")
	var list = make([]data.VersionType, 0)
	var ver data.VersionType
	for _, line := range lines {
		if line == "" {
			continue
		}
		err := json.Unmarshal([]byte(line), &ver)
		cobra.CheckErr(err)
		list = append(list, ver)
	}
	return list
}

func findVersion(versionId string, suffix string) (data.VersionType, error) {
	list := ListHistory(suffix)
	for _, ver := range list {
		if strings.HasPrefix(ver.Id, versionId) {
			return ver, nil
		}
	}
	return data.VersionType{}, fmt.Errorf("version not found")
}

func FindVersion(versionId string) (data.VersionType, error) {
	remoteVersion, err := findVersion(versionId, "")
	if err == nil {
		return remoteVersion, nil
	}
	localVersion, err := findVersion(versionId, "-local")
	if err != nil {
		return localVersion, nil
	}
	return data.VersionType{}, fmt.Errorf("version not found")
}

func MoveVersion(target data.VersionType) {
	localList := ListHistory("-local")
	remoteList := ListHistory("")

	newLocalList := make([]data.VersionType, 0)
	for _, ver := range localList {
		if ver.Id == target.Id {
			continue
		}
		newLocalList = append(newLocalList, ver)
	}

	remoteList = append(remoteList, target)
	sort.Slice(remoteList, func(i, j int) bool {
		return remoteList[i].Time < remoteList[j].Time
	})

	writeFile(filepath.Join(DATADIR, HISTORY_FILE), versionListToString(remoteList))
	writeFile(filepath.Join(DATADIR, HISTORY_FILE+"-local"), versionListToString(newLocalList))
}

func versionListToString(list []data.VersionType) string {
	var str = ""
	for _, ver := range list {
		line, err := versionToString(ver)
		cobra.CheckErr(err)
		str += line + "\n"
	}
	return str
}

func readFile(file string) (string, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func writeFile(file string, data string) error {
	return os.WriteFile(file, []byte(data), os.ModePerm)
}

func appendFile(file string, data string) error {
	f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY, os.ModePerm)
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
