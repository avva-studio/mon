package plog

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestHookFires(t *testing.T) {
	var buffer bytes.Buffer

	l := logrus.New()
	l.Hooks = logrus.LevelHooks{}
	l.Hooks.Add(TraceHook{})
	l.Out = &buffer
	l.Formatter = &logrus.JSONFormatter{}

	l.Error("Error test logging")

	m := make(map[string]interface{})
	err := json.Unmarshal(buffer.Bytes(), &m)

	assert.Nil(t, err, "Json Marshall failed")

	assert.Equal(t, "Error test logging", m["msg"], "Message is not equal")
	assert.Equal(t, "tracehook_test.go", m["file"], "Filename is not as expected")
	assert.Equal(t, "go-logger", m["function"], "Function name is not as expected")
	assert.Equal(t, float64(21), m["line"], "Line is not as expected")
}
