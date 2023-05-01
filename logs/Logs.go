package logs

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
)

type LogLevel int

const CriticalLevel = 1
const ErrorLevel = 2
const WarningLevel = 3
const InfoLevel = 4
const DebugLevel = 5

//goland:noinspection GoUnusedExportedFunction
func LogLevelName(level LogLevel) string {
	list := map[LogLevel]string{
		DebugLevel:    "debug",
		InfoLevel:     "information",
		WarningLevel:  "warning",
		ErrorLevel:    "error",
		CriticalLevel: "critical error",
	}
	return list[level]
}

//goland:noinspection SpellCheckingInspection
func LogLevelTag(level LogLevel) string {
	list := map[LogLevel]string{
		DebugLevel:    "DEBUG",
		InfoLevel:     "INFO",
		WarningLevel:  "WARN",
		ErrorLevel:    "ERROR",
		CriticalLevel: "CRIT",
	}
	return list[level]
}

//goland:noinspection GoNameStartsWithPackageName
const LogsRuntimeCallerSkip = 4

type Logs struct {
	level    LogLevel
	fileName string
	stdOut   bool
}

var logger *Logs = nil

func Logger() *Logs {
	if logger == nil {
		logger = new(Logs)
		logger.stdOut = false
		logger.fileName = ""
		logger.level = InfoLevel
	}
	return logger
}

func setLevel(level LogLevel) {
	Debug("Logs", fmt.Sprintf("SetLevel: %s", LogLevelTag(level)), nil)
	Logger().level = level
}

//goland:noinspection GoUnusedExportedFunction
func SetLevelDebug() { setLevel(DebugLevel) }

//goland:noinspection GoUnusedExportedFunction
func SetLevelInfo() { setLevel(InfoLevel) }

//goland:noinspection GoUnusedExportedFunction
func SetLevelWarn() { setLevel(WarningLevel) }

//goland:noinspection GoUnusedExportedFunction
func SetLevelError() { setLevel(ErrorLevel) }

//goland:noinspection GoUnusedExportedFunction
func SetStdOut(enable bool) {
	Debug("Logs", fmt.Sprintf("SetStdOut: '%t'", enable), nil)
	if enable {
		Info("Logs", fmt.Sprintf("Using StdOut"), nil)
		log.SetOutput(os.Stdout)
		Logger().stdOut = enable
	}
}

//goland:noinspection GoUnusedExportedFunction
func SetFileName(fileName string) {
	Debug("Logs", fmt.Sprintf("SetFileName: '%s'", fileName), nil)
	if len(fileName) > 0 {
		f, e := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
		if e == nil {
			Info("Logs", fmt.Sprintf("Using '%s'", fileName), nil)
			log.SetOutput(f)
			Logger().fileName = fileName
		} else {
			Warn("Logs", fmt.Sprintf("SetFileName: cannot write in file '%s'", fileName), e)
		}
	}
}

func contextMethod(callerSkip int) string {
	if callerSkip > 0 {
		callerSkip += LogsRuntimeCallerSkip
	} else {
		callerSkip = LogsRuntimeCallerSkip
	}
	pc, _, _, _ := runtime.Caller(callerSkip)
	runtimeContext := runtime.FuncForPC(pc).Name()
	return strings.Split(runtimeContext, "/")[1]
}

func formatError(prefix string, e error) string {
	if e != nil {
		return fmt.Sprintf("%s: %s", prefix, e.Error())
	}
	return ""
}

func formatLog(level LogLevel, title string, message string) string {
	logMessage := fmt.Sprintf("[%s]", LogLevelTag(level))
	logMessage = fmt.Sprintf("%-8s", logMessage)
	if len(title) > 0 && len(message) > 0 {
		logMessage += fmt.Sprintf("%-24s %s", title, message)
	} else if len(message) > 0 {
		logMessage += fmt.Sprintf("%s", message)
	} else {
		logMessage = ""
	}
	return logMessage
}

func logMessage(level LogLevel, title string, message string, e error) {
	switch title {
	case "", "current":
		title = contextMethod(0)
	case "parent":
		title = contextMethod(1)
	}
	if level <= Logger().level && len(message) > 0 {
		if logMessage := formatLog(level, title, message); len(logMessage) > 0 {
			log.Printf("%s", logMessage)

		}
	}
	if e != nil {
		if level > ErrorLevel {
			level = ErrorLevel
		}
		if logMessage := formatLog(level, title, e.Error()); len(logMessage) > 0 {
			log.Printf("%s", logMessage)
		}
		_, _ = os.Stderr.WriteString(fmt.Sprintf("%s\n", formatError(title, e)))
	}
}

//goland:noinspection GoUnusedExportedFunction
func Debug(title string, message string, e error) {
	logMessage(DebugLevel, title, message, e)
}

//goland:noinspection GoUnusedExportedFunction
func Info(title string, message string, e error) {
	logMessage(InfoLevel, title, message, e)
}

//goland:noinspection GoUnusedExportedFunction
func Warn(title string, message string, e error) {
	logMessage(WarningLevel, title, message, e)
}

//goland:noinspection GoUnusedExportedFunction
func Error(title string, message string, e error) {
	logMessage(ErrorLevel, title, message, e)
}

//goland:noinspection GoUnusedExportedFunction
func Critical(title string, message string, e error) {
	logMessage(ErrorLevel, title, message, e)
}

//goland:noinspection GoUnusedExportedFunction
func CriticalExit(title string, message string, e error) {
	Critical(title, message, e)
	os.Exit(1)
}
