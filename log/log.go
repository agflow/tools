package log

import (
	"fmt"
	"log"
	"os"
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

// Level type
type Level uint32

const (
	// PanicLvl level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLvl Level = iota
	// FatalLvl level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLvl
	// ErrorLvl level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLvl
	// WarnLvl level. Non-critical entries that deserve eyes.
	WarnLvl
	// InfoLvl level. General operational entries about what's going on inside the
	// application.
	InfoLvl
	// DebugLvl level. Usually only enabled when debugging. Very verbose logging.
	DebugLvl
)

// Hook is an alias for the hook function
type Hook = func(MetaInfo) error

// Logger defines the logger to be used
type Logger struct {
	Hooks []Hook
}

// nolint: gochecknoglobals
var logger Logger

// nolint: gochecknoinits
func init() {
	log.SetFlags(0)
	logger = Logger{}
}

// MetaInfo is the metadata of the log
type MetaInfo struct {
	Msg string
	Lvl Level
}

// Info logs in info level
func Info(v ...interface{}) {
	exec(MetaInfo{Msg: fmt.Sprint(v...), Lvl: InfoLvl})
}

// Infof logs in info level with a format
func Infof(format string, v ...interface{}) {
	Info(fmt.Sprintf(format, v...))
}

// Warn logs in warn level
func Warn(v ...interface{}) {
	exec(MetaInfo{Msg: fmt.Sprint(v...), Lvl: WarnLvl})
}

// Warnf logs in warn level with a format
func Warnf(format string, v ...interface{}) {
	Warn(fmt.Sprintf(format, v...))
}

// Debug logs in debug level
func Debug(v ...interface{}) {
	exec(MetaInfo{Msg: fmt.Sprint(v...), Lvl: DebugLvl})
}

// Debugf logs in debug level with a format
func Debugf(format string, v ...interface{}) {
	Debug(fmt.Sprintf(format, v...))
}

// IfErrorDiffNil logs Error if argument is not nil. DEPRECATED
func IfErrorDiffNil(v interface{}) {
	if v != nil {
		Error(v)
	}
}

// ErrorType logs the error if the argument is not nil
func ErrorType(err error) {
	if err != nil {
		Error(err)
	}
}

// Error logs in error level
func Error(v ...interface{}) {
	exec(MetaInfo{Msg: fmt.Sprint(v...), Lvl: ErrorLvl})
}

// Errorf logs in error level with a format
func Errorf(format string, v ...interface{}) {
	Error(fmt.Sprintf(format, v...))
}

// Fatalf logs in error level with a format
func Fatalf(format string, v ...interface{}) {
	Fatal(fmt.Sprintf(format, v...))
}

// Fatal logs in error level
func Fatal(v ...interface{}) {
	exec(MetaInfo{Msg: fmt.Sprint(v...), Lvl: FatalLvl})
}

// Panic logs in error level with a posterior Panic()
func Panic(v ...interface{}) {
	exec(MetaInfo{Msg: fmt.Sprint(v...), Lvl: PanicLvl})
}

// Panicf logs in error level with a posterior Panic() with format
func Panicf(format string, v ...interface{}) {
	Panic(fmt.Sprintf(format, v...))
}

// AddHook adds hook to the logger
func AddHook(hook Hook) {
	logger.Hooks = append(logger.Hooks, hook)
}

// EmptyHooks empties the hooks of the logger
func EmptyHooks(hook Hook) {
	logger.Hooks = []Hook{}
}

func (li *MetaInfo) log() {
	nowStr := time.Now().Format(timeFormat)
	switch li.Lvl {
	case PanicLvl:
		log.Printf("%s%s %s %s%s", red, nowStr, "[PANIC]", li.Msg, reset)
	case FatalLvl:
		log.Printf("%s%s %s %s%s", red, nowStr, "[FATAL]", li.Msg, reset)
	case ErrorLvl:
		log.Printf("%s%s %s %s%s", red, nowStr, "[ERROR]", li.Msg, reset)
	case WarnLvl:
		log.Printf("%s%s %s %s%s", yellow, nowStr, "[WARN]", li.Msg, reset)
	case InfoLvl:
		log.Printf("%s%s %s %s%s", green, nowStr, "[INFO]", li.Msg, reset)
	case DebugLvl:
		execLine := ""
		_, file, line, ok := runtime.Caller(1)
		if ok {
			execLine = fmt.Sprintf("%s:%d", filepath.Base(file), line)
		}
		log.Printf("%s%s %s %s %s%s", yellow, nowStr, execLine, "[DEBUG]", li.Msg, reset)
	}
}

func exec(li MetaInfo) {
	li.log()
	for _, hook := range logger.Hooks {
		nowStr := time.Now().Format(timeFormat)
		if err := hook(li); err != nil {
			log.Printf("%s%s %s %v%s", red, nowStr, "[ERROR]", err, reset)
		}
	}
	switch li.Lvl {
	case PanicLvl:
		panic(li.Msg)
	case FatalLvl:
		os.Exit(1)
	}
}
