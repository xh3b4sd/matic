package generator

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func skipFile(path string, info os.FileInfo) bool {
	// Skip directories.
	if info.IsDir() {
		return true
	}

	// Skip none go files.
	if filepath.Ext(path) != ".go" {
		return true
	}

	return false
}

func readFile(path string) (string, error) {
	reader, err := os.Open(path)
	if err != nil {
		return "", Mask(err)
	}

	byteSlice, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", Mask(err)
	}

	return string(byteSlice), nil
}
