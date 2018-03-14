package plog

import (
	"os"

	"github.com/sirupsen/logrus"
)

var defaultOption = option{
	format: &logrus.JSONFormatter{},
	output: os.Stderr,
	level:  logrus.InfoLevel,
}

// DebugOptions will provide the defaults for debugging purposes
func DebugOptions() Option {
	return func(o *option) {
		o.output = os.Stdout
		o.format = &logrus.TextFormatter{FullTimestamp: true}
		o.level = logrus.DebugLevel
	}
}

// Option function to process given options
type Option func(*option)

type option struct {
	format logrus.Formatter
	output *os.File
	level  logrus.Level
}

// TextFormat will set the formatter to text
func TextFormat() Option {
	return func(o *option) {
		o.format = &logrus.TextFormatter{FullTimestamp: true}
	}
}

// JSONFormat will set the formatter to json
func JSONFormat() Option {
	return func(o *option) {
		o.format = &logrus.JSONFormatter{}
	}
}

// OutputStdErr will set the output to stderr
func OutputStdErr() Option {
	return func(o *option) {
		o.output = os.Stderr
	}
}

// OutputStdOut will set the output to stdout
func OutputStdOut() Option {
	return func(o *option) {
		o.output = os.Stdout
	}
}

// LevelWarn will set the level to Warn
func LevelWarn() Option {
	return func(o *option) {
		o.level = WarnLevel
	}
}

// LevelDebug will set the level to Debug
func LevelDebug() Option {
	return func(o *option) {
		o.level = DebugLevel
	}
}

// LevelInfo will set the level to Info
func LevelInfo() Option {
	return func(o *option) {
		o.level = InfoLevel
	}
}

// LevelError will set the level to Error
func LevelError() Option {
	return func(o *option) {
		o.level = ErrorLevel
	}
}
