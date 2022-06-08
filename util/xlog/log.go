package xlog

import (
	"edance/util"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
)

type Log struct {
	filePath    string
	format      *LogFormat
	file        io.Writer
	debugLogger *log.Logger
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
}

type LogFormat struct {
	flag     int
	preDebug string
	preInfo  string
	preWarn  string
	preError string
}

func (logger *Log) Init() error {
	f, err := os.OpenFile(logger.filePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return Wrap(err)
	} else if f == nil {
		return Wrap(errors.New("file handler is nil"))
	}
	logger.file = f
	if logger.format != nil {
		logger.debugLogger = log.New(logger.file, logger.format.preDebug, logger.format.flag)
		logger.infoLogger = log.New(logger.file, logger.format.preInfo, logger.format.flag)
		logger.warnLogger = log.New(logger.file, logger.format.preWarn, logger.format.flag)
		logger.errorLogger = log.New(logger.file, logger.format.preError, logger.format.flag)
	}
	return nil
}

func (logger *Log) Debug(format string, v ...interface{}) {
	if !util.Exists(logger.filePath) {
		err := logger.Init()
		if err != nil {
			fmt.Println("Logging debug message failed, error: " + err.Error())
			return
		}
	}
	_ = logger.debugLogger.Output(2, fmt.Sprintf(format, v...))
}

func (logger *Log) Info(format string, v ...interface{}) {
	if !util.Exists(logger.filePath) {
		err := logger.Init()
		if err != nil {
			fmt.Println("Logging info message failed, error: " + err.Error())
			return
		}
	}
	_ = logger.infoLogger.Output(2, fmt.Sprintf(format, v...))
}

func (logger *Log) Warn(format string, v ...interface{}) {
	if !util.Exists(logger.filePath) {
		err := logger.Init()
		if err != nil {
			fmt.Println("Logging warn message failed, error: " + err.Error())
			return
		}
	}
	_ = logger.warnLogger.Output(2, fmt.Sprintf(format, v...))
}

func (logger *Log) Error(format string, v ...interface{}) {
	if !util.Exists(logger.filePath) {
		err := logger.Init()
		if err != nil {
			fmt.Println("Logging error message failed, error: " + err.Error())
			return
		}
	}
	_ = logger.errorLogger.Output(2, fmt.Sprintf(format, v...))
}

func DefaultLog(logFile string) (*Log, error) {
	defaultLogger := &Log{
		filePath: logFile,
		format:   defaultLogFormat(),
	}
	err := defaultLogger.Init()
	if err != nil {
		return nil, Wrap(err)
	}
	return defaultLogger, nil
}

func defaultLogFormat() *LogFormat {
	return &LogFormat{
		flag:     log.Ldate | log.Ltime | log.Lshortfile,
		preDebug: "[DEBUG] ",
		preInfo:  "[INFO ] ",
		preWarn:  "[WARN ] ",
		preError: "[ERROR] ",
	}
}

func Wrap(err error) error {
	return errors.Wrap(err, StackTrace(0))
}

func StackTrace(skip int) string {
	pc, file, line, ok := runtime.Caller(skip + 2)
	if !ok {
		return ""
	}
	funcName := runtime.FuncForPC(pc).Name()
	if strings.Contains(funcName, "/") {
		split := strings.Split(funcName, "/")
		funcName = split[len(split)-1]
	}
	if strings.Contains(funcName, ".") {
		split := strings.Split(funcName, ".")
		funcName = split[len(split)-1]
	}
	return fmt.Sprintf("(%s.%s:%d)", file, funcName, line)
}
