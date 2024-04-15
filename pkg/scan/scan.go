package scan

import (
	"os"
	"path/filepath"
)

type Rule struct {
	Excludes []string
	Includes []string
}

func Dirs(dirs []string, rule *Rule) (files []string, err error) {
	for _, dir := range dirs {
		filesScanned, err := subDir(dir, rule)
		if err != nil {
			return []string{}, err
		}
		files = append(files, filesScanned...)
	}

	return files, nil
}

func subDir(dir string, rule *Rule) (files []string, err error) {
	dir = filepath.Clean(dir)
	err = filepath.Walk(dir, func(file string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		} else if file == dir {
			files = append(files, filepath.ToSlash(file))
			return filepath.SkipDir
		}

		// Skip if
		// we have includes patterns, and we don't match it
		// we match excludes patterns
		if (len(rule.Includes) > 0 && !fileMatch(info.Name(), rule.Includes)) || fileMatch(info.Name(), rule.Excludes) {
			return nil
		}

		files = append(files, filepath.ToSlash(file))
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

func fileMatch(filename string, patterns []string) bool {
	for _, pattern := range patterns {
		match, err := filepath.Match(pattern, filename)
		if err != nil {
			return false
		}
		if match {
			return true
		}
	}

	return false
}
