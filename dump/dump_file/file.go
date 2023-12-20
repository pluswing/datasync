package dump_file

import (
	"os"
	"path/filepath"

	cp "github.com/otiai10/copy"
	"github.com/pluswing/datasync/data"
	"github.com/pluswing/datasync/file"
	"github.com/spf13/cobra"
)

func Dump(dumpDir string, cfg data.TargetFileType) {
	cwd, err := file.FindCurrentDir()
	cobra.CheckErr(err)
	src := filepath.Join(cwd, cfg.Path)
	s, err := os.Stat(src)
	cobra.CheckErr(err)
	dest := filepath.Join(dumpDir, cfg.Path)
	destDir := dest
	if !s.IsDir() {
		destDir = filepath.Dir(dest)
	}
	err = os.MkdirAll(destDir, os.ModePerm)
	cobra.CheckErr(err)
	err = cp.Copy(src, dest)
	cobra.CheckErr(err)
}

func Expand(dumpDir string, cfg data.TargetFileType) {
	cwd, err := file.FindCurrentDir()
	cobra.CheckErr(err)
	src := filepath.Join(dumpDir, cfg.Path)
	s, err := os.Stat(src)
	cobra.CheckErr(err)
	dest := filepath.Join(cwd, cfg.Path)
	destDir := dest
	if !s.IsDir() {
		destDir = filepath.Dir(dest)
	}
	err = os.MkdirAll(destDir, os.ModePerm)
	cobra.CheckErr(err)
	err = cp.Copy(src, dest)
	cobra.CheckErr(err)
}
