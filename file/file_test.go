package file

import (
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

func TestFindCurrentDir(t *testing.T) {
	home := setup(t)
	defer teardown(home)

	dir, err := FindCurrentDir()
	assert.Error(t, err)
	assert.Equal(t, "", dir)

	os.WriteFile(filepath.Join(home, "datasync.yaml"), []byte(""), os.ModePerm)

	dir, err = FindCurrentDir()
	assert.NoError(t, err)
	assert.Equal(t, home, dir)
}

// func TestReadVersionFile(t *testing.T) {
// 	assert.Equal(t, "d37be652686e4373bd01ca528b01f31d", ReadVersionFile())
// }
