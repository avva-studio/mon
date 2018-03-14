package plog

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestCreateWithDebugOptions(t *testing.T) {
	l := New(DebugOptions())

	if l == (logger{}) {
		t.Error("Logger should be instantiated")
	}

	lrText := new(logrus.TextFormatter)
	lrText.FullTimestamp = true
	lrDebug := logrus.DebugLevel

	assert.Equal(t, l.Entry().Logger.Formatter, lrText, "Should have Text formatter")
	assert.Equal(t, l.Entry().Logger.Level, lrDebug, "Should have debug level")
}

func TestCreateWithFormatOptions(t *testing.T) {
	l := New(TextFormat())

	if l == (logger{}) {
		t.Error("Logger should be instantiated")
	}

	lrText := new(logrus.TextFormatter)
	lrText.FullTimestamp = true

	assert.Equal(t, l.Entry().Logger.Formatter, lrText, "Should have Text formatter")

	lJSON := New(JSONFormat())

	lrJSON := new(logrus.JSONFormatter)
	assert.Equal(t, lJSON.Entry().Logger.Formatter, lrJSON, "Should have JSON formatter")
}

func TestCreateLoggerWithSetter(t *testing.T) {
	l := New()

	lrJSON := new(logrus.JSONFormatter)
	assert.Equal(t, l.Entry().Logger.Formatter, lrJSON, "Should have JSON formatter")

	l.SetLevel(WarnLevel)
	lrWarn := logrus.WarnLevel
	assert.Equal(t, l.Entry().Logger.Level, lrWarn, "Should have warn level")

	l.SetDebugOptions()
	lrDebug := logrus.DebugLevel
	lrText := new(logrus.TextFormatter)
	lrText.FullTimestamp = true
	assert.Equal(t, l.Entry().Logger.Formatter, lrText, "Should have Text formatter")
	assert.Equal(t, l.Entry().Logger.Level, lrDebug, "Should have debug level")

}

func TestLevelDebug(t *testing.T) {
	l := New(LevelDebug())
	assert.Equal(t, l.Entry().Logger.Level, logrus.DebugLevel, "Level should be equal to debug")
}

func TestLevelInfo(t *testing.T) {
	l := New(LevelInfo())

	assert.Equal(t, l.Entry().Logger.Level, logrus.InfoLevel, "Level should be equal to info")
}

func TestLevelError(t *testing.T) {
	l := New(LevelError())

	assert.Equal(t, l.Entry().Logger.Level, logrus.ErrorLevel, "Level should be equal to error")
}

func TestLevelWarn(t *testing.T) {
	l := New(LevelWarn())

	assert.Equal(t, l.Entry().Logger.Level, logrus.WarnLevel, "Level should be equal to error")
}

func TestOutputStdErr(t *testing.T) {
	l := New(OutputStdErr())

	lr := logrus.New()
	out := lr.Out

	assert.Equal(t, l.Entry().Logger.Out, out, "Should be the same error output, default setting logrus")

	lr.Out = os.Stdout
	out = lr.Out

	assert.NotEqual(t, l.Entry().Logger.Out, out, "Should not log to error output")
}

func TestOutputStdOut(t *testing.T) {
	l := New(OutputStdOut())

	lr := logrus.New()
	lr.Out = os.Stdout
	out := lr.Out

	assert.Equal(t, l.Entry().Logger.Out, out, "Should be the same standard output")

	lr.Out = os.Stderr
	out = lr.Out

	assert.NotEqual(t, l.Entry().Logger.Out, out, "Should not log to standard output")
}
