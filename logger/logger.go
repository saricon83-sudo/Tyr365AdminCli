package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

var logger = logrus.New()
var file *os.File // Global file object for logging
func SetupLogging(logFile string, useJSON bool) {
	if useJSON {
		logger.Formatter = &logrus.JSONFormatter{}
	} else {
		logger.Formatter = &logrus.TextFormatter{
			ForceColors:   true, // Force enabling colors even if non-tty output is detected
			FullTimestamp: true, // Show full timestamp
		}
	}
	logger.SetLevel(logrus.InfoLevel)

	if logFile != "" {
		file, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			logger.Warn("Failed to open log file, logging to standard output")
			logger.SetOutput(os.Stdout)
			return
		}
		logger.SetOutput(file)
		logger.AddHook(&commaHook{File: file}) // Add hook to append commas
	} else {
		logger.SetOutput(os.Stdout)
	}
}

type commaHook struct {
	File *os.File
}

func (hook *commaHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *commaHook) Fire(entry *logrus.Entry) error {
	_, err := hook.File.WriteString(",")
	return err
}

func CloseLogging() {
	if file != nil {
		file.Close()
	}
}

func GetLogger() *logrus.Logger {
	return logger
}
