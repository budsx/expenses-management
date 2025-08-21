package util

import (
	"fmt"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
)

func NewLogger(level int) *logrus.Logger {
	log := logrus.New()

	log.SetLevel(logrus.Level(level))
	log.SetReportCaller(true)
	log.SetFormatter(&logrus.JSONFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			caller := fmt.Sprintf("%s:%d", filename, f.Line)
			return caller, ""
		},
	})

	return log
}
