package plog

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestBaseLoggerDebugOptions(t *testing.T) {
	SetDebugOptions()

	lr := new(logrus.TextFormatter)
	lr.FullTimestamp = true
	assert.Equal(t, oriLogger.Formatter, lr, "The debug options should be set equal")
	assert.Equal(t, oriLogger.Level, logrus.DebugLevel, "the debug options should be set equal")
}

func TestBaseLoggerDebug(t *testing.T) {
	var buffer bytes.Buffer

	SetDebugOptions()
	SetJSONFormatter()

	oriLogger.Out = &buffer

	Debug("Debug test logging")

	m := make(map[string]interface{})
	err := json.Unmarshal(buffer.Bytes(), &m)

	if err != nil {
		t.Error("json marshall failed")
	}

	assert.Equal(t, "Debug test logging", m["msg"], "Message is not equal")
	assert.Equal(t, "base_test.go", m["file"], "Filename is not as expected")
	assert.Equal(t, "TestBaseLoggerDebug", m["function"], "Function name is not as expected")
	assert.Equal(t, float64(30), m["line"], "Line is not as expected")
}

func TestBaseLoggerError(t *testing.T) {
	var buffer bytes.Buffer

	SetDebugOptions()
	SetJSONFormatter()

	oriLogger.Out = &buffer

	Error("Error test logging")

	m := make(map[string]interface{})
	err := json.Unmarshal(buffer.Bytes(), &m)

	if err != nil {
		t.Error("json marshall failed")
	}

	assert.Equal(t, "Error test logging", m["msg"], "Message is not equal")
	assert.Equal(t, "base_test.go", m["file"], "Filename is not as expected")
	assert.Equal(t, "TestBaseLoggerError", m["function"], "Function name is not as expected")
	assert.Equal(t, float64(53), m["line"], "Line is not as expected")
}

func TestBaseLoggerInfo(t *testing.T) {
	var buffer bytes.Buffer

	SetDebugOptions()
	SetJSONFormatter()

	oriLogger.Out = &buffer

	Info("Info test logging")

	m := make(map[string]interface{})
	err := json.Unmarshal(buffer.Bytes(), &m)

	if err != nil {
		t.Error("json marshall failed")
	}

	assert.Equal(t, "Info test logging", m["msg"], "Message is not equal")
	assert.Equal(t, "base_test.go", m["file"], "Filename is not as expected")
	assert.Equal(t, "TestBaseLoggerInfo", m["function"], "Function name is not as expected")
	assert.Equal(t, float64(76), m["line"], "Line is not as expected")
}

func TestBaseLoggerWarn(t *testing.T) {
	var buffer bytes.Buffer

	SetDebugOptions()
	SetJSONFormatter()

	oriLogger.Out = &buffer

	Warn("Warn test logging")

	m := make(map[string]interface{})
	err := json.Unmarshal(buffer.Bytes(), &m)

	if err != nil {
		t.Error("json marshall failed")
	}

	assert.Equal(t, "Warn test logging", m["msg"], "Message is not equal")
	assert.Equal(t, "base_test.go", m["file"], "Filename is not as expected")
	assert.Equal(t, "TestBaseLoggerWarn", m["function"], "Function name is not as expected")
	assert.Equal(t, float64(99), m["line"], "Line is not as expected")
}

func TestBaseLoggerWithFields(t *testing.T) {
	var buffer bytes.Buffer

	SetDebugOptions()
	SetJSONFormatter()

	oriLogger.Out = &buffer

	With(logrus.Fields{
		"extraFields": "Extra",
	}).Warn("Warn test logging")

	m := make(map[string]interface{})
	err := json.Unmarshal(buffer.Bytes(), &m)

	if err != nil {
		t.Error("json marshall failed")
	}

	assert.Equal(t, "Warn test logging", m["msg"], "Message is not equal")
	assert.Equal(t, "base_test.go", m["file"], "Filename is not as expected")
	assert.Equal(t, "Extra", m["extraFields"], "Extra Fields is not available as expected")
	assert.Equal(t, "TestBaseLoggerWithFields", m["function"], "Function name is not as expected")
	assert.Equal(t, float64(124), m["line"], "Line is not as expected")
}

func TestBaseLoggerSetLevel(t *testing.T) {
	SetLevel(DebugLevel)
	assert.Equal(t, oriLogger.Level, logrus.DebugLevel, "Level should be equal to debug")

	SetLevel(WarnLevel)
	assert.Equal(t, oriLogger.Level, logrus.WarnLevel, "Level should be equal to warn")

	SetLevel(InfoLevel)
	assert.Equal(t, oriLogger.Level, logrus.InfoLevel, "Level should be equal to info")

	SetLevel(ErrorLevel)
	assert.Equal(t, oriLogger.Level, logrus.ErrorLevel, "Level should be equal to error")

}

func TestBaseLoggerFormatter(t *testing.T) {
	SetJSONFormatter()
	lrJSON := new(logrus.JSONFormatter)
	assert.Equal(t, oriLogger.Formatter, lrJSON, "Should have JSON formatter")

	SetTextFormatter()
	lrText := new(logrus.TextFormatter)
	lrText.FullTimestamp = true
	assert.Equal(t, oriLogger.Formatter, lrText, "Should have Text formatter")

}

func TestBaseLoggerSetOutput(t *testing.T) {
	SetStderr()
	assert.Equal(t, oriLogger.Out, os.Stderr, "Output should be equal to stderr")

	SetStdout()
	assert.Equal(t, oriLogger.Out, os.Stdout, "Output should be equal to stdout")
}
