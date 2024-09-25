package server

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

var LOG *logrus.Logger

type CustomTextFormatter struct{}

func (f *CustomTextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	message := fmt.Sprintf(
		"%s %5.5s %s\n",
		entry.Time.Format("2006-01-02 15:04:05.000"), // Date-time
		strings.ToUpper(entry.Level.String()),        // Log level
		entry.Message,                                // Log message
	)

	return []byte(message), nil
}

const (
	DirPerm  = 0766
	FilePerm = 0666
)

func init() {
	file, err := os.OpenFile(
		"server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666,
	)
	if err != nil {
		panic(err)
	}
	LOG = &logrus.Logger{
		Out:       file,
		Level:     logrus.TraceLevel,
		Formatter: &CustomTextFormatter{},
	}
}

func createUploadsDir(dirPath string) error {
	return os.MkdirAll(dirPath, DirPerm)
}

func prepareFile(filePath string, fileSize int64) (file *os.File,
	fileExists bool, err error) {
	file, err = os.OpenFile(filePath, os.O_CREATE|os.O_EXCL|os.O_RDWR, FilePerm)
	if err != nil {
		if os.IsExist(err) {
			return nil, true, err
		}
		return nil, false, err
	}
	err = file.Truncate(0)
	if err != nil {
		_ = file.Close()
		return nil, false, err
	}
	err = file.Truncate(fileSize)
	if err != nil {
		_ = file.Close()
		return nil, false, err
	}
	return file, false, nil
}
