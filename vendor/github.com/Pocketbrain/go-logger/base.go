package plog

import (
	"os"

	"github.com/sirupsen/logrus"
)

type logger struct {
	entry *logrus.Entry
}

var baseLogger = logger{
	entry: logrus.NewEntry(oriLogger),
}

func init() {
	SetJSONFormatter()
}

const (
	// PanicLevel is the Panic Level on Logrus
	PanicLevel = logrus.PanicLevel
	// FatalLevel is the Fatal Level on Logrus
	FatalLevel = logrus.FatalLevel
	// DebugLevel is the Debug Level on Logrus
	DebugLevel = logrus.DebugLevel
	// ErrorLevel is the Error Level on Logrus
	ErrorLevel = logrus.ErrorLevel
	// WarnLevel is the Warn Level on Logrus
	WarnLevel = logrus.WarnLevel
	// InfoLevel is the Info Level on Logrus
	InfoLevel = logrus.InfoLevel
)

// SetDebugOptions will provide the defaults for debugging purposes
func SetDebugOptions() {
	SetStdout()
	SetLevel(DebugLevel)
	SetTextFormatter()
}

// Panic logs a message at Panic level
func Panic(args ...interface{}) {
	baseLogger.source().Panic(args...)
}

// Fatal logs a message at Fatal level
func Fatal(args ...interface{}) {
	baseLogger.source().Fatal(args...)
}

// Debug logs a message at Debug level
func Debug(args ...interface{}) {
	baseLogger.source().Debug(args...)
}

// Info logs a message at Info level
func Info(args ...interface{}) {
	baseLogger.source().Info(args...)
}

// Warn logs a message at Warn level
func Warn(args ...interface{}) {
	baseLogger.source().Warn(args...)
}

// Error logs a message at Error level
func Error(args ...interface{}) {
	baseLogger.source().Error(args...)
}

// With adds a field to the logger.
func With(fields logrus.Fields) Logger {
	return baseLogger.With(fields)
}

// SetJSONFormatter will set the baseLogger formatter with JSON
func SetJSONFormatter() {
	if inDebugMode() {
		return
	}
	oriLogger.Formatter = &logrus.JSONFormatter{}
}

// SetTextFormatter will set the baseLogger formatter with Text
func SetTextFormatter() {
	if inDebugMode() {
		return
	}
	oriLogger.Formatter = &logrus.TextFormatter{FullTimestamp: true}
}

// SetStdout will set the baseLogger output
func SetStdout() {
	if inDebugMode() {
		return
	}
	oriLogger.Out = os.Stdout
}

// SetStderr will set the baseLogger output
func SetStderr() {
	if inDebugMode() {
		return
	}
	oriLogger.Out = os.Stderr
}

// SetLevel will set the baseLogger level
func SetLevel(l logrus.Level) {
	if inDebugMode() {
		return
	}
	oriLogger.Level = l
}
