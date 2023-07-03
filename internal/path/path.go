package path

import (
	"fmt"
	"path/filepath"
	"runtime"
)

func trimPackagePath(pkgPath string) (string, error) {
	_, thisFile, _, _ := runtime.Caller(1)
	baseImportPath := filepath.Dir(thisFile)

	relPath, err := filepath.Rel(baseImportPath, pkgPath)
	if err != nil {
		return "", fmt.Errorf("target path cannot be made relative to basepath")
	}

	trimmedPath := filepath.ToSlash(relPath)

	return trimmedPath, nil
}
