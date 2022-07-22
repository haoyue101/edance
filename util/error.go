package util

import (
	"fmt"
	"github.com/pkg/errors"
	"runtime"
	"strings"
)

func Wrap(err error, cause ...string) error {
	if cause == nil {
		return errors.Wrap(err, StackTrace(0))
	} else {
		format := "["
		for i := 0; i < len(cause)-1; i++ {
			format += "%s, "
		}
		format += "]"
		return errors.Wrapf(err, StackTrace(0), format, cause)
	}
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
	return fmt.Sprintf("\n\t%s.%s:%d", file, funcName, line)
}
