package pathutil

import (
	"path/filepath"
	"strings"
)

func LocalFromReal(root, realPath string) (string, error) {
	localPath, err := filepath.Rel(root, realPath)
	if err != nil {
		return "", err
	}

	// WARN: may break with ../stuff -> /../stuff -> /stuff if path is outside of root (prob wont happen)
	localPath = filepath.Clean("/" + localPath)

	return localPath, nil
}

func RealFromLocal(root, localPath string) string {
	realPath := strings.TrimLeft(localPath, "/")
	realPath = filepath.Join(root, localPath)

	return realPath
}

func TargetPath(sourcePath, targetPath string) string {
	if strings.HasPrefix(targetPath, "/") {
		return targetPath
	}
	return filepath.Join(filepath.Dir(sourcePath), targetPath)
}

func Language(path string) string {
	if ext := filepath.Ext(path); len(ext) >= 1 {
		return filepath.Ext(path)[1:]
	}
	return ""
}
