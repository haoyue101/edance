package util

import (
	"os"
	"strings"
)

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func IsFile(path string) bool {
	return !IsDir(path)
}

func GetFileDir(file string) string {
	if !strings.Contains(file, "/") {
		return "/"
	}
	lastIndex := strings.LastIndex(file, "/")
	return file[:lastIndex+1]
}

func CreateIfNotExist(file string) (*os.File, error) {
	if Exists(file) {
		return nil, nil
	}
	err := os.MkdirAll(GetFileDir(file), os.ModePerm)
	if err != nil {
		return nil, Wrap(err)
	}
	filePtr, err := os.Create(file)
	if err != nil {
		return nil, Wrap(err)
	}
	return filePtr, nil
}
