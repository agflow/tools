package log

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"time"
)

const (
	reset      = "\033[0m"
	red        = "\033[31m"
	green      = "\033[32m"
	yellow     = "\033[33m"
	timeFormat = "2006/01/02 15:04:05"
)

//nolint: gochecknoinits
func init() {
	log.SetFlags(0)
}

// Info logs in info level
func Info(v ...interface{}) {
	nowStr := time.Now().Format(timeFormat)
	log.Printf("%s%s %s %s%s", green, nowStr, "[INFO]", fmt.Sprint(v...), reset)
}

// Infof logs in info level with a format
func Infof(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	Info(msg)
}

// Warn logs in warn level
func Warn(v ...interface{}) {
	nowStr := time.Now().Format(timeFormat)
	log.Printf("%s%s %s %s%s", yellow, nowStr, "[WARN]", fmt.Sprint(v...), reset)
}

// Warnf logs in warn level with a format
func Warnf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	Warn(msg)
}

// Debug logs in debug level
func Debug(v ...interface{}) {
	nowStr := time.Now().Format(timeFormat)

	execLine := ""
	_, file, line, ok := runtime.Caller(1)
	if ok {
		execLine = fmt.Sprintf("%s:%d", filepath.Base(file), line)
	}

	log.Printf("%s%s %s %s %s%s", yellow, nowStr, execLine, "[DEBUG]", fmt.Sprint(v...), reset)
}

// Debugf logs in debug level with a format
func Debugf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	Debug(msg)
}

// IfErrorDiffNil logs Error if argument is not nil
func IfErrorDiffNil(v interface{}) {
	if v != nil {
		Error(v)
	}
}

// Error logs in error level
func Error(v ...interface{}) {
	nowStr := time.Now().Format(timeFormat)
	log.Printf("%s%s %s %s%s", red, nowStr, "[ERROR]", fmt.Sprint(v...), reset)
}

// Errorf logs in error level with a format
func Errorf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	Error(msg)
}

// Fatalf logs in error level with a format
func Fatalf(format string, v ...interface{}) {
	log.Fatalf(format, v...)
}

// Fatal logs in error level
func Fatal(v ...interface{}) {
	log.Fatal(v...)
}
