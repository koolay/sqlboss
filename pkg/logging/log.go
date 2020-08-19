package logging

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

func NewLogger(levelName string) *logrus.Logger {
	level, err := logrus.ParseLevel(levelName)
	if err != nil {
		level = logrus.ErrorLevel
	}

	logFormatter := &logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
		FullTimestamp:   true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			s := strings.Split(f.Function, ".")
			funcname := s[len(s)-1]
			_, filename := path.Split(f.File)
			return funcname, fmt.Sprintf("%s:%v", filename, f.Line)
		},
	}

	logger := logrus.New()
	logger.SetLevel(level)
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(logFormatter)
	logger.SetReportCaller(true)

	return logger
}
