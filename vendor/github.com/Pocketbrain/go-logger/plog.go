package plog

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

// Logger is the interface for loggers
type Logger interface {
	Panic(...interface{})
	Fatal(...interface{})
	Warn(...interface{})
	Error(...interface{})
	Info(...interface{})
	Debug(...interface{})

	SetLevel(logrus.Level)

	SetJSONFormatter()
	SetTextFormatter()

	SetStderr()
	SetStdout()

	SetDebugOptions()

	Entry() *logrus.Entry
	With(field logrus.Fields) Logger
}

var oriLogger = logrus.New()

// New will instantiate new logger
func New(opts ...Option) Logger {
	options := defaultOption
	for _, o := range opts {
		o(&options)
	}

	if inDebugMode() {
		options.level = DebugLevel
		options.format = flags.formatter
		options.output = flags.stdout
	}

	l := logrus.New()
	l.Formatter = options.format
	l.Out = options.output
	l.Level = options.level

	return logger{entry: logrus.NewEntry(l)}
}

// SetDebugOptions will provide the defaults for debugging purposes
func (l logger) SetDebugOptions() {
	if inDebugMode() {
		return
	}
	l.SetLevel(DebugLevel)
	l.SetTextFormatter()
}

// Panic logs a message at level Panic on the standard logger.
func (l logger) Panic(args ...interface{}) {
	l.source().Panic(args...)
}

// Fatal logs a message at level Fatal on the standard logger.
func (l logger) Fatal(args ...interface{}) {
	l.source().Fatal(args...)
}

// Warn logs a message at level Warn on the standard logger.
func (l logger) Warn(args ...interface{}) {
	l.source().Warn(args...)
}

// Error logs a message at level Error on the standard logger.
func (l logger) Error(args ...interface{}) {
	l.source().Error(args...)
}

// Info logs a message at level Info on the standard logger.
func (l logger) Info(args ...interface{}) {
	l.source().Info(args...)
}

// Debug logs a message at level Debug on the standard logger.
func (l logger) Debug(args ...interface{}) {
	l.source().Debug(args...)
}

// Entry will return the logrus.Entry
func (l logger) Entry() *logrus.Entry {
	return l.entry
}

// With will add extra fields
func (l logger) With(fields logrus.Fields) Logger {
	return logger{l.entry.WithFields(fields)}
}

// source add the source fields to the logger
// contain filename, function and line where the logging happens
func (l logger) source() *logrus.Entry {
	pc, file, line, ok := runtime.Caller(2)

	if !ok {
		return nil
	}

	slash := strings.LastIndex(file, "/")
	file = file[slash+1:]
	funcName := runtime.FuncForPC(pc).Name()
	fName := strings.Split(path.Base(funcName), ".")

	return l.entry.WithFields(logrus.Fields{
		"file":     fmt.Sprintf("%s", file),
		"function": fmt.Sprintf("%s", fName[len(fName)-1]),
		"line":     line,
	})

}

// SetLevel will set the standar logger level
func (l logger) SetLevel(level logrus.Level) {
	if inDebugMode() {
		return
	}
	l.entry.Logger.Level = level
}

// SetJSONFormatter will set the standard formatter with JSON
func (l logger) SetJSONFormatter() {
	if inDebugMode() {
		return
	}
	l.entry.Logger.Formatter = &logrus.JSONFormatter{}
}

// SetTextFormatter will set the standard formatter with Text
func (l logger) SetTextFormatter() {
	if inDebugMode() {
		return
	}
	l.entry.Logger.Formatter = &logrus.TextFormatter{FullTimestamp: true}
}

// SetStderr will set the standard output with Stderr
func (l logger) SetStderr() {
	if inDebugMode() {
		return
	}
	l.entry.Logger.Out = os.Stderr
}

// SetStdout will set the standard output with Stdout
func (l logger) SetStdout() {
	if inDebugMode() {
		return
	}
	l.entry.Logger.Out = os.Stdout
}
