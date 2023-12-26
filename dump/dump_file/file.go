package dump_file

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/briandowns/spinner"
	cp "github.com/otiai10/copy"
	"github.com/pluswing/datasync/data"
	"github.com/pluswing/datasync/file"
	"github.com/spf13/cobra"
)

func Dump(dumpDir string, cfg data.TargetFileType) {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = fmt.Sprintf(" dump file(s) ... (path: %s)", cfg.Path)
	s.Start()

	cwd, err := file.FindCurrentDir()
	cobra.CheckErr(err)
	src := filepath.Join(cwd, cfg.Path)
	stat, err := os.Stat(src)
	cobra.CheckErr(err)
	dest := filepath.Join(dumpDir, cfg.Path)
	destDir := dest
	if !stat.IsDir() {
		destDir = filepath.Dir(dest)
	}
	err = os.MkdirAll(destDir, os.ModePerm)
	cobra.CheckErr(err)
	err = cp.Copy(src, dest)
	cobra.CheckErr(err)

	s.Stop()
	fmt.Printf("✔︎ dump file(s) completed. (path: %s)", cfg.Path)
}

func Expand(dumpDir string, cfg data.TargetFileType) {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = fmt.Sprintf(" restore file(s) ... (path: %s)", cfg.Path)
	s.Start()

	cwd, err := file.FindCurrentDir()
	cobra.CheckErr(err)
	src := filepath.Join(dumpDir, cfg.Path)
	stat, err := os.Stat(src)
	cobra.CheckErr(err)
	dest := filepath.Join(cwd, cfg.Path)
	destDir := dest
	if !stat.IsDir() {
		destDir = filepath.Dir(dest)
	}
	err = os.MkdirAll(destDir, os.ModePerm)
	cobra.CheckErr(err)
	err = cp.Copy(src, dest)
	cobra.CheckErr(err)

	s.Stop()
	fmt.Printf("✔︎ restore file(s) completed. (path: %s)", cfg.Path)
}
