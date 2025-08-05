package logger

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/sirupsen/logrus"
)

var globalLogger *logrus.Logger

func InitLogger() {
	globalLogger = logrus.New()
	formatter := &logrus.JSONFormatter{
		TimestampFormat:   "2006-01-02 15:04:05",
		DisableHTMLEscape: true,
		PrettyPrint:       true,
	}
	globalLogger.SetFormatter(formatter)
	globalLogger.SetLevel(logrus.InfoLevel)
}

func getCallerInfo() (string, int) {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return "", 0
	}

	return filepath.Base(file), line
}

func LogInfo(message string) {
	file, line := getCallerInfo()
	globalLogger.WithFields(logrus.Fields{
		"file": fmt.Sprintf("%s:%d", file, line),
	}).Info(message)
}

func LogWarn(message string) {
	file, line := getCallerInfo()
	globalLogger.WithFields(logrus.Fields{
		"file": fmt.Sprintf("%s:%d", file, line),
	}).Warn(message)
}

func LogError(message string) {
	file, line := getCallerInfo()
	globalLogger.WithFields(logrus.Fields{
		"file": fmt.Sprintf("%s:%d", file, line),
	}).Error(message)
}

func LogFatal(message string) {
	file, line := getCallerInfo()
	globalLogger.WithFields(logrus.Fields{
		"file": fmt.Sprintf("%s:%d", file, line),
	}).Fatal(message)
}

func LogPanic(message string) {
	file, line := getCallerInfo()
	globalLogger.WithFields(logrus.Fields{
		"file": fmt.Sprintf("%s:%d", file, line),
	}).Panic(message)
}

func LogWithFields(fields logrus.Fields, message string) {
	file, line := getCallerInfo()
	fields["file"] = fmt.Sprintf("%s:%d", file, line)
	globalLogger.WithFields(fields).Info(message)
}
