package util

import (
	"github.com/sirupsen/logrus"
)

func NewLogger(level int) *logrus.Logger {
	log := logrus.New()

	log.SetLevel(logrus.Level(level))
	log.SetFormatter(&logrus.JSONFormatter{})

	return log
}
