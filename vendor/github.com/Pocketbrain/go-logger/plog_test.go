package plog

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestCreateLogger(t *testing.T) {
	l := New()

	if l == (logger{}) {
		t.Error("Logger should be instantiated")
	}

	lrJSON := new(logrus.JSONFormatter)
	lrLevel := logrus.InfoLevel

	assert.Equal(t, l.Entry().Logger.Formatter, lrJSON, "Should have JSON formatter")
	assert.Equal(t, l.Entry().Logger.Level, lrLevel, "Should have Info level")
}

func TestLoggerDebug(t *testing.T) {
	var buffer bytes.Buffer

	l := New(DebugOptions())
	l.SetJSONFormatter()

	l.Entry().Logger.Out = &buffer

	l.Debug("Debug test logging")

	m := make(map[string]interface{})
	err := json.Unmarshal(buffer.Bytes(), &m)

	if err != nil {
		t.Error("json marshall failed")
	}

	assert.Equal(t, "Debug test logging", m["msg"], "Message is not equal")
	assert.Equal(t, "plog_test.go", m["file"], "Filename is not as expected")
	assert.Equal(t, "TestLoggerDebug", m["function"], "Function name is not as expected")
	assert.Equal(t, float64(35), m["line"], "Line is not as expected")
}

func TestLoggerError(t *testing.T) {
	var buffer bytes.Buffer

	l := New(DebugOptions())
	l.SetJSONFormatter()

	l.Entry().Logger.Out = &buffer

	l.Error("Error test logging")

	m := make(map[string]interface{})
	err := json.Unmarshal(buffer.Bytes(), &m)

	if err != nil {
		t.Error("json marshall failed")
	}

	assert.Equal(t, "Error test logging", m["msg"], "Message is not equal")
	assert.Equal(t, "plog_test.go", m["file"], "Filename is not as expected")
	assert.Equal(t, "TestLoggerError", m["function"], "Function name is not as expected")
	assert.Equal(t, float64(58), m["line"], "Line is not as expected")
}

func TestLoggerInfo(t *testing.T) {
	var buffer bytes.Buffer

	l := New(DebugOptions())
	l.SetJSONFormatter()

	l.Entry().Logger.Out = &buffer

	l.Info("Info test logging")

	m := make(map[string]interface{})
	err := json.Unmarshal(buffer.Bytes(), &m)

	if err != nil {
		t.Error("json marshall failed")
	}

	assert.Equal(t, "Info test logging", m["msg"], "Message is not equal")
	assert.Equal(t, "plog_test.go", m["file"], "Filename is not as expected")
	assert.Equal(t, "TestLoggerInfo", m["function"], "Function name is not as expected")
	assert.Equal(t, float64(81), m["line"], "Line is not as expected")
}

func TestLoggerWarn(t *testing.T) {
	var buffer bytes.Buffer

	l := New(DebugOptions())
	l.SetJSONFormatter()

	l.Entry().Logger.Out = &buffer

	l.Warn("Warn test logging")

	m := make(map[string]interface{})
	err := json.Unmarshal(buffer.Bytes(), &m)

	if err != nil {
		t.Error("json marshall failed")
	}

	assert.Equal(t, "Warn test logging", m["msg"], "Message is not equal")
	assert.Equal(t, "plog_test.go", m["file"], "Filename is not as expected")
	assert.Equal(t, "TestLoggerWarn", m["function"], "Function name is not as expected")
	assert.Equal(t, float64(104), m["line"], "Line is not as expected")
}

func TestLoggerWithFields(t *testing.T) {
	var buffer bytes.Buffer

	l := New(DebugOptions())
	l.SetJSONFormatter()

	l.Entry().Logger.Out = &buffer

	l.With(logrus.Fields{
		"extraFields": "Extra",
	}).Warn("Warn test logging")

	m := make(map[string]interface{})
	err := json.Unmarshal(buffer.Bytes(), &m)

	if err != nil {
		t.Error("json marshall failed")
	}

	assert.Equal(t, "Warn test logging", m["msg"], "Message is not equal")
	assert.Equal(t, "plog_test.go", m["file"], "Filename is not as expected")
	assert.Equal(t, "Extra", m["extraFields"], "Extra Fields is not available as expected")
	assert.Equal(t, "TestLoggerWithFields", m["function"], "Function name is not as expected")
	assert.Equal(t, float64(129), m["line"], "Line is not as expected")
}

func TestLoggerFormatter(t *testing.T) {
	l := New()

	l.SetJSONFormatter()
	lrJSON := new(logrus.JSONFormatter)
	assert.Equal(t, l.Entry().Logger.Formatter, lrJSON, "Should have JSON formatter")

	l.SetTextFormatter()
	lrText := new(logrus.TextFormatter)
	lrText.FullTimestamp = true
	assert.Equal(t, l.Entry().Logger.Formatter, lrText, "Should have Text formatter")

}

func TestLoggerOutput(t *testing.T) {
	l := New()
	l.SetStderr()
	assert.Equal(t, l.Entry().Logger.Out, os.Stderr, "Output should be equal to stderr")

	l.SetStdout()
	assert.Equal(t, l.Entry().Logger.Out, os.Stdout, "Output should be equal to stdout")
}
