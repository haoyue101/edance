package util

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

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
	if !Exists(logger.filePath) {
		err := logger.Init()
		if err != nil {
			fmt.Println("Logging debug message failed, error: " + err.Error())
			return
		}
	}
	_ = logger.debugLogger.Output(2, fmt.Sprintf(format, v...))
}

func (logger *Log) Info(format string, v ...interface{}) {
	if !Exists(logger.filePath) {
		err := logger.Init()
		if err != nil {
			fmt.Println("Logging info message failed, error: " + err.Error())
			return
		}
	}
	_ = logger.infoLogger.Output(2, fmt.Sprintf(format, v...))
}

func (logger *Log) Warn(format string, v ...interface{}) {
	if !Exists(logger.filePath) {
		err := logger.Init()
		if err != nil {
			fmt.Println("Logging warn message failed, error: " + err.Error())
			return
		}
	}
	_ = logger.warnLogger.Output(2, fmt.Sprintf(format, v...))
}

func (logger *Log) Error(format string, v ...interface{}) {
	if !Exists(logger.filePath) {
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

func formatter(param gin.LogFormatterParams, bodyStr string) string {
	var statusColor, methodColor, resetColor string
	if param.IsOutputColor() {
		statusColor = param.StatusCodeColor()
		methodColor = param.MethodColor()
		resetColor = param.ResetColor()
	}

	if param.Latency > time.Minute {
		param.Latency = param.Latency - param.Latency%time.Second
	}
	return fmt.Sprintf("[EDance] %v |%s %3d %s| %13v | %15s |%s %-7s %s %#v%s  %s \n",
		param.TimeStamp.Format("2006/01/02 - 15:04:05"),
		statusColor,
		param.StatusCode,
		resetColor,
		param.Latency,
		param.ClientIP,
		methodColor,
		param.Method,
		resetColor,
		param.Path,
		bodyStr,
		param.ErrorMessage,
	)
}

func XLog() gin.HandlerFunc {
	out := io.MultiWriter(*LogFile(), os.Stdout)
	var notLogged []string
	var skip map[string]struct{}
	if length := len(notLogged); length > 0 {
		skip = make(map[string]struct{}, length)
		for _, path := range notLogged {
			skip[path] = struct{}{}
		}
	}
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		//get request body
		buff, _ := ioutil.ReadAll(c.Request.Body)
		bodyStr := ""
		if buff != nil && len(buff) > 0 {
			t := io.NopCloser(bytes.NewBuffer(buff))
			c.Request.Body = t
			bodyStr = string(buff)
			bodyStr = strings.ReplaceAll(bodyStr, "\n", "")
			bodyStr = strings.ReplaceAll(bodyStr, "\"", "")
			bodyStr = strings.ReplaceAll(bodyStr, "\t", "")
			bodyStr = strings.ReplaceAll(bodyStr, " ", "")
		}
		c.Next()
		// Log only when path is not being skipped
		if _, ok := skip[path]; !ok {
			param := gin.LogFormatterParams{Request: c.Request,
				Keys: c.Keys}
			param.TimeStamp = time.Now()
			param.Latency = param.TimeStamp.Sub(start)
			param.ClientIP = c.ClientIP()
			param.Method = c.Request.Method
			param.StatusCode = c.Writer.Status()
			param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()
			param.BodySize = c.Writer.Size()
			if raw != "" {
				path = path + "?" + raw
			}
			param.Path = path
			fmt.Fprint(out, formatter(param, bodyStr))
		}
	}
}
