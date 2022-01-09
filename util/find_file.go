package util

import (
	"io/ioutil"
	"path/filepath"
)

// FindFiles returns all *.<ext> file pathes under `dir` (recursively).
func FindFiles(dir string, ext string) []string {
	var result []string
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		path := filepath.Join(dir, file.Name())
		if file.IsDir() {
			res := FindFiles(path, ext)
			for _, path := range res {
				result = append(result, path)
			}
		} else if filepath.Ext(path) == ext {
			result = append(result, path)
		}
	}
	return result
}
