package logger

import (
	"fmt"
	"strings"

	"github.com/artwebs/aogo/utils"
)

// Log levels to control the logging output.
const (
	LoggerFatal = iota
	LoggerError
	LoggerWarn
	LoggerNotice
	LoggerInfo
	LoggerDebug
)

var LoggerLevel = LoggerDebug
var LoggerText = map[int]string{LoggerFatal: "Fatal ", LoggerError: "Error ", LoggerWarn: "Warn  ", LoggerNotice: "Notice", LoggerInfo: "Info  ", LoggerDebug: "Debug "}
var LoggerColor = map[int]string{LoggerFatal: "\033[0;36;48m", LoggerError: "\033[0;35;48m", LoggerWarn: "\033[0;31;48m", LoggerNotice: "\033[0;33;48m", LoggerInfo: "\033[0;32;48m", LoggerDebug: "\033[0;34;48m"}

// SetLogLevel sets the global log level used by the simple
// logger.
func SetLevel(l int) {
	LoggerLevel = l
}

func Fatal(v ...interface{}) {
	print(LoggerFatal, v...)
}

// Error logs a message at error level.
func Error(v ...interface{}) {
	print(LoggerError, v...)
}

// Error logs a message at error level.
func ErrorTag(v ...interface{}) {
	if len(v) > 1 {
		v[0] = utils.Tag(v[0]) + "=>"
	}
	Error(v...)
}

// Warning logs a message at warning level.
func Warn(v ...interface{}) {
	print(LoggerWarn, v...)
}

func WarnTag(v ...interface{}) {
	if len(v) > 1 {
		v[0] = utils.Tag(v[0]) + "=>"
	}
	Warn(v...)
}

func Notice(v ...interface{}) {
	print(LoggerNotice, v...)
}

func NoticeTag(v ...interface{}) {
	if len(v) > 1 {
		v[0] = utils.Tag(v[0]) + "=>"
	}
	Notice(v...)
}

// Info logs a message at info level.
func Info(v ...interface{}) {
	print(LoggerInfo, v...)
}

func InfoTag(v ...interface{}) {
	if len(v) > 1 {
		v[0] = utils.Tag(v[0]) + "=>"
	}
	Info(v...)
}

// Debug logs a message at debug level.
func Debug(v ...interface{}) {
	print(LoggerDebug, v...)
}

func DebugTag(v ...interface{}) {
	if len(v) > 1 {
		v[0] = utils.Tag(v[0]) + "=>"
	}
	Debug(v...)
}

func generateFmtStr(n int) string {
	return strings.Repeat("%v ", n)
}

func print(l int, v ...interface{}) {
	if l <= LoggerLevel {
		text := ""
		color := ""
		if val, ok := LoggerText[l]; ok {
			text = val
		}
		if val, ok := LoggerColor[l]; ok {
			color = val
		}
		fmt.Println(utils.NowDateTime() + " [" + color + strings.ToUpper(text) + "\033[0m] " + fmt.Sprintf(generateFmtStr(len(v)), v...))
	}
}
