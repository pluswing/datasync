package file

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigFiles(t *testing.T) {
	cwd, _ := os.Getwd()
	dir, err := FindCurrentDir()
	assert.NoError(t, err)
	assert.Equal(t, filepath.Dir(cwd), dir)
}

func TestReadVersionFile(t *testing.T) {
	assert.Equal(t, "f56e90715043400ab8adf5d18f984105", ReadVersionFile())
}
