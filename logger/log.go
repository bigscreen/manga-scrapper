package logger

import (
	"log"
	"os"

	"github.com/bigscreen/manga-scrapper/config"
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func SetupLogger() *logrus.Logger {
	if logger != nil {
		return logger
	}
	level, err := logrus.ParseLevel(config.LogLevel())
	if err != nil {
		log.Fatalf(err.Error())
	}

	logger = &logrus.Logger{
		Out:       os.Stderr,
		Formatter: &logrus.JSONFormatter{},
		Hooks:     make(logrus.LevelHooks),
		Level:     level,
	}

	return logger
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func DebugF(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func ErrorF(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func InfoF(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

func Warn(args ...interface{}) {
	logger.Warn(args...)
}

func WarnF(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

func FatalF(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}
