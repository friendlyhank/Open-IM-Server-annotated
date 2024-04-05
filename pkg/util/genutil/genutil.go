package genutil

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/OpenIMSDK/tools/errs"
)

// OutDir creates the absolute path name from path and checks path exists.
// Returns absolute path including trailing '/' or error if path does not exist.
func OutDir(path string) (string, error) {
	outDir, err := filepath.Abs(path)
	if err != nil {
		return "", errs.Wrap(err, "output directory %s does not exist", path)
	}

	stat, err := os.Stat(outDir)
	if err != nil {
		return "", errs.Wrap(err, "output directory %s does not exist", outDir)
	}

	if !stat.IsDir() {
		return "", errs.Wrap(err, "output directory %s is not a directory", outDir)
	}
	outDir += "/"
	return outDir, nil
}

func ExitWithError(err error) {
	progName := filepath.Base(os.Args[0])
	fmt.Fprintf(os.Stderr, "%s exit -1: %+v\n", progName, err)
	os.Exit(-1)
}
