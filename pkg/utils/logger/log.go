package logger

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var (
	entry *logrus.Entry
)

func init() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	entry = logger.WithField("sessionID", uuid.New().String())
	entry.Info("New session started")
}

func Info(msg ...interface{}) {
	entry.Info(msg...)
}

func Error(msg ...interface{}) {
	entry.Error(msg...)
}
