package utils

import (
	"fmt"
	"os"
)

func MustExist(dirpath ...string) string {
	var path string
	for _, path = range dirpath {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			panic(fmt.Sprintf("create directory %q: %v", path, err))
		}
	}
	return path
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	panic(err)
}
