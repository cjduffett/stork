package logger

import "fmt"

// LogLevel indicates the current level of logging, from 0-2:
// 0 - Info only
// 1 - Info and Warning
// 2 - Info, Warning, and Debug
var LogLevel int

const (
	DefaultLevel = 0
	InfoLevel    = 0
	WarnLevel    = 1
	DebugLevel   = 2
)

// Info logs non-critical information to the console at any LogLevel.
func Info(v ...interface{}) {
	logStork("INFO", v...)
}

// Error logs critical information to the console at any LogLevel.
func Error(v ...interface{}) {
	logStork("ERROR", v...)
}

// Warning logs warning information to the console at LogLevel 1 or higher.
func Warning(v ...interface{}) {
	if LogLevel >= 1 {
		logStork("WARN", v...)
	}
}

// Debug logs debug information to the console at LogLevel 2 or higher.
func Debug(v ...interface{}) {
	if LogLevel == 2 {
		logStork("DEBUG", v...)
	}
}

func logStork(level string, v ...interface{}) {
	fmt.Println(append([]interface{}{fmt.Sprintf("[Stork] [%s] ", level)}, v...)...)
}
