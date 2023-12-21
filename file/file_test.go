package file

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

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

func TestUpdateVersionFile(t *testing.T) {
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
