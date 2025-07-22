package resolver

import (
	"path/filepath"
)

func FindFilesWithExclude(pattern string, excludes []string) ([]string, error) {
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	var result []string
FILES:
	for _, file := range files {
		for _, exclude := range excludes {
			match, err := filepath.Match(exclude, filepath.Base(file))
			if err != nil {
				return nil, err
			}
			if match {
				continue FILES
			}
		}
		result = append(result, file)
	}

	return result, nil
}
