package file

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/pluswing/datasync/data"
	"github.com/stretchr/testify/assert"
)

func setup(t *testing.T) string {
	dir, err := os.MkdirTemp("", ".datasync_for_test")
	assert.NoError(t, err)
	os.Chdir(dir)
	cwd, err := os.Getwd()
	assert.NoError(t, err)
	return cwd
}

func teardown(dir string) {
	os.RemoveAll(dir)
}

func createFile(home string, file string, content string) {
	os.WriteFile(filepath.Join(home, file), []byte(content), os.ModePerm)
}

func TestFindCurrentDir(t *testing.T) {
	home := setup(t)
	defer teardown(home)

	dir, err := FindCurrentDir()
	assert.Error(t, err)
	assert.Equal(t, fmt.Errorf("config file not found"), err)
	assert.Equal(t, "", dir)

	createFile(home, "datasync.yaml", "")
	dir, err = FindCurrentDir()
	assert.NoError(t, err)
	assert.Equal(t, home, dir)
}

func TestReadVersionFile(t *testing.T) {
	home := setup(t)
	defer teardown(home)

	assert.Equal(t, "", ReadVersionFile())

	createFile(home, ".datasync_version", "")
	assert.Equal(t, "", ReadVersionFile())

	createFile(home, "datasync.yaml", "")
	assert.Equal(t, "", ReadVersionFile())

	createFile(home, ".datasync_version", "123")
	assert.Equal(t, "123", ReadVersionFile())

	createFile(home, ".datasync_version", "456\n")
	assert.Equal(t, "456", ReadVersionFile())
}

func TestVersionFile(t *testing.T) {
	home := setup(t)
	defer teardown(home)

	createFile(home, "datasync.yaml", "")

	err := UpdateVersionFile("345")
	assert.NoError(t, err)
	assert.Equal(t, "345", ReadVersionFile())

	err = UpdateVersionFile("456")
	assert.NoError(t, err)
	assert.Equal(t, "456", ReadVersionFile())
}

func TestDataDir(t *testing.T) {
	home := setup(t)
	defer teardown(home)

	dir, err := DataDir()
	assert.Error(t, err)
	assert.Equal(t, fmt.Errorf("config file not found"), err)
	assert.Equal(t, "", dir)

	createFile(home, "datasync.yaml", "")
	dir, err = DataDir()
	assert.NoError(t, err)
	assert.Equal(t, filepath.Join(home, ".datasync"), dir)
}

func TestDataSyncFile(t *testing.T) {
	home := setup(t)
	defer teardown(home)

	empty := data.DataSyncType{
		Version:   "1",
		Histories: []data.VersionType{},
	}

	createFile(home, "datasync.yaml", "")

	r := ReadRemoteDataSyncFile()
	assert.Equal(t, empty, r)

	r = ReadLocalDataSyncFile()
	assert.Equal(t, empty, r)

	localData := data.DataSyncType{
		Version: "1",
		Histories: []data.VersionType{
			{
				Id:      "1",
				Time:    1,
				Message: "a",
			}, {
				Id:      "2",
				Time:    2,
				Message: "b",
			},
		},
	}

	WriteLocalDataSyncFile(localData)

	r = ReadLocalDataSyncFile()
	assert.Equal(t, localData, r)

	moveData := localData.Histories[0]
	restData := localData.Histories[1:]
	MoveVersionToRemote(moveData)

	r = ReadRemoteDataSyncFile()
	assert.Equal(t, data.DataSyncType{
		Version:   "1",
		Histories: []data.VersionType{moveData},
	}, r)

	r = ReadLocalDataSyncFile()
	assert.Equal(t, data.DataSyncType{
		Version:   "1",
		Histories: restData,
	}, r)
}

func TestGetCurrentVersion(t *testing.T) {
	home := setup(t)
	defer teardown(home)

	createFile(home, "datasync.yaml", "")

	v, err := GetCurrentVersion([]string{})
	assert.Error(t, err)
	assert.Equal(t, fmt.Errorf("version not found."), err)
	assert.Equal(t, data.VersionType{}, v)

	UpdateVersionFile("123")

	v, err = GetCurrentVersion([]string{})
	assert.Error(t, err)
	assert.Equal(t, fmt.Errorf("version not found."), err)
	assert.Equal(t, data.VersionType{}, v)

	localData := data.DataSyncType{
		Version: "1",
		Histories: []data.VersionType{
			{
				Id:      "123",
				Time:    1,
				Message: "a",
			},
			{
				Id:      "456",
				Time:    2,
				Message: "b",
			},
		},
	}

	WriteLocalDataSyncFile(localData)

	v, err = GetCurrentVersion([]string{})
	assert.NoError(t, err)
	assert.Equal(t, localData.Histories[0], v)

	v, err = GetCurrentVersion([]string{"456"})
	assert.NoError(t, err)
	assert.Equal(t, localData.Histories[1], v)

	v, err = GetCurrentVersion([]string{"789"})
	assert.Error(t, err)
	assert.Equal(t, fmt.Errorf("version not found."), err)
	assert.Equal(t, data.VersionType{}, v)

	MoveVersionToRemote(localData.Histories[0])

	v, err = GetCurrentVersion([]string{})
	assert.NoError(t, err)
	assert.Equal(t, localData.Histories[0], v)

	v, err = GetCurrentVersion([]string{"456"})
	assert.NoError(t, err)
	assert.Equal(t, localData.Histories[1], v)

	v, err = GetCurrentVersion([]string{"789"})
	assert.Error(t, err)
	assert.Equal(t, fmt.Errorf("version not found."), err)
	assert.Equal(t, data.VersionType{}, v)

	MoveVersionToRemote(localData.Histories[1])

	v, err = GetCurrentVersion([]string{})
	assert.NoError(t, err)
	assert.Equal(t, localData.Histories[0], v)

	v, err = GetCurrentVersion([]string{"456"})
	assert.NoError(t, err)
	assert.Equal(t, localData.Histories[1], v)

	v, err = GetCurrentVersion([]string{"789"})
	assert.Error(t, err)
	assert.Equal(t, fmt.Errorf("version not found."), err)
	assert.Equal(t, data.VersionType{}, v)
}

func TestNewUUID(t *testing.T) {
	uuid, err := NewUUID()
	assert.NoError(t, err)
	assert.Equal(t, 32, len(uuid))

	uuid2, err := NewUUID()
	assert.NoError(t, err)
	assert.Equal(t, 32, len(uuid2))

	assert.NotEqual(t, uuid, uuid2)
}
