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
