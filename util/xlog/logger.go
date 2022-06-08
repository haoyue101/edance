package xlog

import "io"

var logger *Log

func Init(logFile string) {
	defaultLogger, err := DefaultLog(logFile)
	if err != nil {
		return
	}
	logger = defaultLogger
}

func LogFile() *io.Writer {
	return &logger.file
}

func Debug(format string, v ...interface{}) {
	logger.Debug(format, v...)
}

func Info(format string, v ...interface{}) {
	logger.Info(format, v...)
}

func Warn(format string, v ...interface{}) {
	logger.Warn(format, v...)
}

func Error(format string, v ...interface{}) {
	logger.Error(format, v...)
}
