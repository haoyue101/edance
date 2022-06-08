package util

import (
	"fmt"
	"testing"
)

func TestGetFileDir(t *testing.T) {
	fmt.Println(GetFileDir("/opt/bin/test.sh"))
}

func TestCreateIfNotExist(t *testing.T) {
	_, err := CreateIfNotExist("/opt/bin/test.sh")
	fmt.Println(err)
}
