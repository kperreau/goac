package utils

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func CleanPath(path string, filename string) string {
	return filepath.Clean(strings.TrimSuffix(path, filename))
}

func AddCurrentDirPrefix(path string) string {
	if !strings.HasPrefix(path, "./") && !strings.HasPrefix(path, "/") {
		return "./" + path
	}
	return path
}

func FileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		return false
	}
	return true
}

// AppendIfNotExist append new values to the existing slice only if the values are not already in.
// This avoids duplication.
func AppendIfNotExist(slice []string, elems ...string) []string {
	for _, elem := range elems {
		if !slices.Contains(slice, elem) {
			slice = append(slice, elem)
		}
	}
	return slice
}
