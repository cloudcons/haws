package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
)

// LogLevel defines the level of logging
type LogLevel int

const (
	// LevelDebug is for detailed debugging information
	LevelDebug LogLevel = iota
	// LevelInfo is for general operational entries
	LevelInfo
	// LevelWarn is for warning events
	LevelWarn
	// LevelError is for error events
	LevelError
)

var (
	// CurrentLevel is the current log level
	CurrentLevel = LevelInfo
	
	// Color formatters
	debugColor = color.New(color.FgCyan)
	infoColor  = color.New(color.FgGreen)
	warnColor  = color.New(color.FgYellow)
	errorColor = color.New(color.FgRed)
	fatalColor = color.New(color.FgRed, color.Bold)
)

// SetLevel sets the current log level
func SetLevel(level LogLevel) {
	CurrentLevel = level
}

// Debug logs debug level messages if the current log level permits
func Debug(format string, args ...interface{}) {
	if CurrentLevel <= LevelDebug {
		log("DEBUG", format, args...)
	}
}

// Info logs info level messages if the current log level permits
func Info(format string, args ...interface{}) {
	if CurrentLevel <= LevelInfo {
		log("INFO ", format, args...)
	}
}

// Warn logs warning level messages if the current log level permits
func Warn(format string, args ...interface{}) {
	if CurrentLevel <= LevelWarn {
		log("WARN ", format, args...)
	}
}

// Error logs error level messages
func Error(format string, args ...interface{}) {
	if CurrentLevel <= LevelError {
		log("ERROR", format, args...)
	}
}

// Fatal logs an error message and then exits
func Fatal(format string, args ...interface{}) {
	log("FATAL", format, args...)
	os.Exit(1)
}

// log formats and prints a log message
func log(level string, format string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, args...)
	
	// Format: [timestamp] LEVEL: message
	fmt.Printf("[%s] ", timestamp)
	
	// Select the appropriate color based on the level
	var levelColor *color.Color
	switch level {
	case "DEBUG":
		levelColor = debugColor
	case "INFO ":
		levelColor = infoColor
	case "WARN ":
		levelColor = warnColor
	case "ERROR":
		levelColor = errorColor
	case "FATAL":
		levelColor = fatalColor
	default:
		levelColor = infoColor
	}
	
	levelColor.Printf("%s: ", level)
	fmt.Println(message)
}
