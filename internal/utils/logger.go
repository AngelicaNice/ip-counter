package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

func InitLogger(logFileName string, logLevel logrus.Level) *logrus.Logger {
	logger := logrus.New()

	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		logrus.Fatalf("Failed to open log file: %v", err)
	}

	logger.Out = logFile
	logger.SetLevel(logLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	return logger
}
