package log

import (
	"fmt"
	"strings"
	"time"
)

// Level data type
type Level int

// Log level constants
const (
	SILENT Level = iota
	ERROR
	WARNING
	INFO
	DEBUG
	TRACE
)

var timeFormat = "2006/01/02 - 15:04:05"
var logLevel = INFO
var indent = 0

var logLevels = map[string]Level{
	"silent":  SILENT,
	"error":   ERROR,
	"warning": WARNING,
	"info":    INFO,
	"debug":   DEBUG,
	"trace":   TRACE,
}

var logIcons = map[string]string{
	"[ERR] ": "🔥",
	"[WRN] ": "⚠️",
	"[INF] ": "ℹ️",
	"[DBG] ": "🎯",
	"[TRC] ": "👀",
}

// SetLogLevelByString sets logging level using its string name
func SetLogLevelByString(newLevel string) (Level, error) {
	oldLevel := logLevel
	newLogLevelValue, ok := logLevels[strings.ToLower(newLevel)]
	if ok {
		logLevel = newLogLevelValue
	} else {
		err := fmt.Errorf("'%s' is invalid value for log level", newLevel)
		Error(err.Error())
		return oldLevel, err
	}

	return oldLevel, nil
}

// Error logs error message
func Error(v ...interface{}) {
	if logLevel >= ERROR {
		out("[ERR] ", fmt.Sprint(v...))
	}
}

// Errorf logs formatted error message
func Errorf(format string, v ...interface{}) {
	if logLevel >= ERROR {
		out("[ERR] ", format, v...)
	}
}

// Warn logs warning message
func Warn(v ...interface{}) {
	if logLevel >= WARNING {
		out("[WRN] ", fmt.Sprint(v...))
	}
}

// Warnf logs formatted warning message
func Warnf(format string, v ...interface{}) {
	if logLevel >= WARNING {
		out("[WRN] ", format, v...)
	}
}

// Info logs informational message
func Info(v ...interface{}) {
	if logLevel >= INFO {
		out("[INF] ", fmt.Sprint(v...))
	}
}

// Infof logs formatted informational message
func Infof(format string, v ...interface{}) {
	if logLevel >= INFO {
		out("[INF] ", format, v...)
	}
}

// Debug logs debug info
func Debug(v ...interface{}) {
	if logLevel >= DEBUG {
		out("[DBG] ", fmt.Sprint(v...))
	}
}

// Debugf logs formatted debug info
func Debugf(format string, v ...interface{}) {
	if logLevel >= DEBUG {
		out("[DBG] ", format, v...)
	}
}

// Trace logs trace info
func Trace(v ...interface{}) {
	if logLevel >= TRACE {
		out("[TRC] ", fmt.Sprint(v...))
	}
}

// Tracef logs formatted trace info
func Tracef(format string, v ...interface{}) {
	if logLevel >= TRACE {
		out("[TRC] ", format, v...)
	}
}

// Fatal logs error and panics
func Fatal(v interface{}) {
	Error(v)
	panic(v)
}

// Fatalf logs formatted error and panics
func Fatalf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	Error(msg)
	panic(msg)
}

// SetIndent func for indents in log
func SetIndent(ind int) {
	indent = ind
}

func out(level string, format string, v ...interface{}) {
	if len(v) > 0 {
		template := strings.Join(
			[]string{
				level,
				time.Now().Format(timeFormat),
				" | ",
				logIcons[level],
				" ",
				strings.Repeat(" ", indent),
				" ",
				format,
				"\n",
			},
			"")
		fmt.Printf(template, v...)
	} else {
		fmt.Printf(
			"%s%v | %s %s %s\n",
			level,
			time.Now().Format(timeFormat),
			logIcons[level],
			strings.Repeat(" ", indent),
			format,
		)
	}
}

// TimeTrack function for printing function time in log
func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	Infof("%s took %s -", name, elapsed)
}
