package util

import (
	"io"
	"os"
	"path/filepath"
)

func CopyFile(srcName string, destName string) {
	src, err := os.Open(srcName)
	if err != nil {
		panic(err)
	}
	defer src.Close()

	os.MkdirAll(filepath.Dir(destName), 0777)

	dest, err := os.Create(destName)
	if err != nil {
		panic(err)
	}
	defer dest.Close()

	_, err = io.Copy(dest, src)
	if err != nil {
		panic(err)
	}
}
