package dump_file

import (
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
	// TODO
	// /xxx/cmd/ccc/ssss
	//ssss

	err = cp.Copy(src, dumpDir)
	cobra.CheckErr(err)
}

func Expand(dumpDir string, cfg data.TargetFileType) {
	src := filepath.Join(dumpDir, cfg.Path)
	cwd, err := file.FindCurrentDir()
	cobra.CheckErr(err)
	err = cp.Copy(src, cwd)
	cobra.CheckErr(err)
}
