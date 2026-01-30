package pathutil

import (
	"path/filepath"
	"strings"
)

// LocalFromReal converts a real path to a local path relative to the root.
func LocalFromReal(root, realPath string) (string, error) {
	localPath, err := filepath.Rel(root, realPath)
	if err != nil {
		return "", err
	}

	// WARN: may break with ../stuff -> /../stuff -> /stuff if path is outside of root (prob wont happen)
	localPath = filepath.Clean("/" + localPath)

	return localPath, nil
}

// RealFromLocal converts a local path relative to the root to a real path.
func RealFromLocal(root, localPath string) string {
	realPath := strings.TrimLeft(localPath, "/")
	realPath = filepath.Join(root, localPath)

	return realPath
}

// TargetPath converts a targetPath assumed to be relative to the sourcePath into local path.
func TargetPath(sourcePath, targetPath string) string {
	if strings.HasPrefix(targetPath, "/") {
		return targetPath
	}
	return filepath.Join(filepath.Dir(sourcePath), targetPath)
}

// LanguageFromPath returns the language of a file path.
func LanguageFromPath(path string) string {
	if ext := filepath.Ext(path); len(ext) >= 1 {
		return filepath.Ext(path)[1:]
	}
	return ""
}
