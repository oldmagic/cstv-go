package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Init initializes the logrus logger
func Init(logLevel string) {
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		level = logrus.InfoLevel
	}

	logrus.SetLevel(level)
	logrus.SetFormatter(&logrus.JSONFormatter{}) // JSON format for structured logging
	logrus.SetOutput(os.Stdout)
}
