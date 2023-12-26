package file

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/google/uuid"
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
		return "", fmt.Errorf("config file not found")
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
	if err != nil {
		return err
	}
	file := filepath.Join(dir, VERSION_FILE)
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

func findVersion(versionId string, suffix string) (data.VersionType, error) {
	ds := readDataSyncFile(suffix)
	for _, ver := range ds.Histories {
		if strings.HasPrefix(ver.Id, versionId) {
			return ver, nil
		}
	}
	return data.VersionType{}, fmt.Errorf("version not found")
}

func findVersionLocalAndRemote(versionId string) (data.VersionType, error) {
	remoteVersion, err := findVersion(versionId, "")
	if err == nil {
		return remoteVersion, nil
	}
	localVersion, err := findVersion(versionId, "-local")
	if err == nil {
		return localVersion, nil
	}
	return data.VersionType{}, fmt.Errorf("version not found")
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

func ReadLocalDataSyncFile() data.DataSyncType {
	return readDataSyncFile("-local")
}

func ReadRemoteDataSyncFile() data.DataSyncType {
	return readDataSyncFile("")
}

func readDataSyncFile(suffix string) (ds data.DataSyncType) {
	ds = data.DataSyncType{
		Version:   "1",
		Histories: []data.VersionType{},
	}
	dir, err := DataDir()
	if err != nil {
		return
	}
	file := filepath.Join(dir, HISTORY_FILE+suffix)
	content, err := readFile(file)
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(content), &ds)
	if err != nil {
		return
	}
	return ds
}

func MoveVersionToRemote(version data.VersionType) {
	local := ReadLocalDataSyncFile()
	remote := ReadRemoteDataSyncFile()

	newLocalList := make([]data.VersionType, 0)
	for _, ver := range local.Histories {
		if ver.Id == version.Id {
			continue
		}
		newLocalList = append(newLocalList, ver)
	}
	local.Histories = newLocalList

	remote.Histories = append(remote.Histories, version)
	sort.Slice(remote.Histories, func(i, j int) bool {
		return remote.Histories[i].Time < remote.Histories[j].Time
	})

	err := WriteLocalDataSyncFile(local)
	cobra.CheckErr(err)
	err = WriteRemoteDataSyncFile(remote)
	cobra.CheckErr(err)
}

func WriteLocalDataSyncFile(d data.DataSyncType) error {
	return writeDataSyncFile(d, "-local")
}

func WriteRemoteDataSyncFile(d data.DataSyncType) error {
	return writeDataSyncFile(d, "")
}

func writeDataSyncFile(d data.DataSyncType, suffix string) error {
	b, err := json.MarshalIndent(d, "", "    ")
	if err != nil {
		return err
	}
	dir, err := DataDir()
	if err != nil {
		return err
	}
	err = writeFile(filepath.Join(dir, HISTORY_FILE+suffix), string(b))
	if err != nil {
		return err
	}
	return nil
}

func GetCurrentVersion(args []string) (data.VersionType, error) {
	var versionId = ""
	if len(args) == 1 {
		versionId = args[0]
	} else {
		versionId = ReadVersionFile()
	}
	if versionId == "" {
		return data.VersionType{}, fmt.Errorf("version not found.")
	}

	version, err := findVersionLocalAndRemote(versionId)
	if err != nil {
		return data.VersionType{}, fmt.Errorf("version not found.")
	}
	return version, nil
}

func NewUUID() (string, error) {
	_uuid, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	uuid := _uuid.String()
	uuid = strings.Replace(uuid, "-", "", -1)
	return uuid, nil
}

func MoveFile(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("Couldn't open source file: %s", err)
	}
	outputFile, err := os.Create(destPath)
	if err != nil {
		inputFile.Close()
		return fmt.Errorf("Couldn't open dest file: %s", err)
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return fmt.Errorf("Writing to output file failed: %s", err)
	}
	// The copy was successful, so now delete the original file
	err = os.Remove(sourcePath)
	if err != nil {
		return fmt.Errorf("Failed removing original file: %s", err)
	}
	return nil
}
