package log

import (
	"strings"
	"github.com/artwebs/aogo/utils"
	"github.com/astaxie/beego/logs"
)

// Log levels to control the logging output.
const (
	LevelEmergency = iota
	LevelAlert
	LevelCritical
	LevelError
	LevelWarning
	LevelNotice
	LevelInformational
	LevelDebug
)

// SetLogLevel sets the global log level used by the simple
// logger.
func SetLevel(l int) {
	BeeLogger.SetLevel(l)
}

func SetLogFuncCall(b bool) {
	BeeLogger.EnableFuncCallDepth(b)
	BeeLogger.SetLogFuncCallDepth(3)
}

// logger references the used application logger.
var BeeLogger *logs.BeeLogger

func NewLogger(channellen int64) {
	BeeLogger = logs.NewLogger(channellen)
}

// SetLogger sets a new logger.
func SetLogger(adaptername string, config string) error {
	err := BeeLogger.SetLogger(adaptername, config)
	if err != nil {
		return err
	}
	return nil
}

func Emergency(v ...interface{}) {
	BeeLogger.Emergency(generateFmtStr(len(v)), v...)
}

func Alert(v ...interface{}) {
	BeeLogger.Alert(generateFmtStr(len(v)), v...)
}

// Critical logs a message at critical level.
func Critical(v ...interface{}) {
	BeeLogger.Critical(generateFmtStr(len(v)), v...)
}

// Error logs a message at error level.
func Error(v ...interface{}) {
	BeeLogger.Error(generateFmtStr(len(v)), v...)
}

// Warning logs a message at warning level.
func Warning(v ...interface{}) {
	BeeLogger.Warning(generateFmtStr(len(v)), v...)
}

// compatibility alias for Warning()
func Warn(v ...interface{}) {
	BeeLogger.Warn(generateFmtStr(len(v)), v...)
}

func WarnTag(v ...interface{}) {
	if len(v)>1 {
		v[0] = utils.Tag(v[0])+"=>"
	}
	Warn(v...)
}

func Notice(v ...interface{}) {
	BeeLogger.Notice(generateFmtStr(len(v)), v...)
}

func NoticeTag(v ...interface{}) {
	if len(v)>1 {
		v[0] = utils.Tag(v[0])+"=>"
	}
	Notice(v...)
}


// Info logs a message at info level.
func Informational(v ...interface{}) {
	BeeLogger.Informational(generateFmtStr(len(v)), v...)
}

// compatibility alias for Warning()
func Info(v ...interface{}) {
	BeeLogger.Info(generateFmtStr(len(v)), v...)
}

func InfoTag(v ...interface{}) {
	if len(v)>1 {
		v[0] = utils.Tag(v[0])+"=>"
	}
	Info(v...)
}

// Debug logs a message at debug level.
func Debug(v ...interface{}) {
	BeeLogger.Debug(generateFmtStr(len(v)), v...)
}

func DebugTag(v ...interface{}) {
	if len(v)>1 {
		v[0] = utils.Tag(v[0])+"=>"
	}
	Debug(v...)
}

// Trace logs a message at trace level.
// compatibility alias for Warning()
func Trace(v ...interface{}) {
	BeeLogger.Trace(generateFmtStr(len(v)), v...)
}

func generateFmtStr(n int) string {
	return strings.Repeat("%v ", n)
}

