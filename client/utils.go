package client

import (
	"errors"
	"os"
)

var (
	NotRegularFile = errors.New("specify regular file")
)

func GetRegularFile(filePath string) (file *os.File, err error) {
	file, err = os.Open(filePath)
	if err != nil {
		return nil, err
	}

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	if !stat.Mode().IsRegular() {
		return nil, NotRegularFile
	}
	return file, nil
}
