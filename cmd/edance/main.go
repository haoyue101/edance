package main

import (
	"bytes"
	"edance/server"
	"edance/util/xlog"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

const PortEnv = "EDANCE_PORT"
const DefaultPort = "10088"
const LogFile = "D:/project-go/edance/log/edance.log"

func main() {
	eng := gin.Default()
	xlog.Init(LogFile)
	eng.Use(log())
	server.InitRouters(eng)
	_ = eng.Run(":" + resolvePort())
}

func resolvePort() string {
	if port := os.Getenv(PortEnv); port != "" {
		portInt, err := strconv.Atoi(port)
		if err != nil || 0 < portInt || portInt > 65535 {
			xlog.Warn("Environment variable " + PortEnv + " is invalid, using default port " + DefaultPort)
			return DefaultPort
		}
		return port
	} else {
		xlog.Info("Environment variable " + PortEnv + " is undefined, using default port " + DefaultPort)
		return DefaultPort
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

func log() gin.HandlerFunc {
	out := io.MultiWriter(*xlog.LogFile(), os.Stdout)
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
