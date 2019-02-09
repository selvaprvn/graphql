package logger

import "log"

var golbalLogger *log.Logger

// SetLogger set global logger
func SetLogger(logger *log.Logger) {
	golbalLogger = logger
}

// Info log
func Info(v ...interface{}) {
	if golbalLogger != nil {
		golbalLogger.Println(v...)
	}
}
