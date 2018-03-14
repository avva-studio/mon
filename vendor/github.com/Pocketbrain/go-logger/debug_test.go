package plog

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"

	"github.com/stretchr/testify/assert"
)

func TestInDebugMode(t *testing.T) {
	debug := inDebugMode()
	assert.Equal(t, false, debug)
	flags.debug = true
	debug = inDebugMode()
	assert.Equal(t, true, debug)
	//for other tests we need to make sure to keep debug mode off
	flags.debug = false
}

func TestStartDebugModeNoFlags(t *testing.T) {
	flagSet := new(flag.FlagSet)
	err := StartDebugMode(flagSet)
	assert.Equal(t, defaultDebugFlags(), flags)
	assert.NotNil(t, err)
}

func TestStartDebugModeDebugFlag(t *testing.T) {
	flags = defaultDebugFlags()
	flagSet := new(flag.FlagSet)
	flagSet.Bool("debug", true, "debug test flag")
	err := StartDebugMode(flagSet)
	expected := defaultDebugFlags()
	expected.debug = true
	assert.Equal(t, expected, flags)
	assert.Nil(t, err)
	//for other tests we need to make sure to keep debug mode off
	flags.debug = false
}

func TestStartDebugModeAllFlags(t *testing.T) {
	flags = defaultDebugFlags()
	flagSet := new(flag.FlagSet)
	flagSet.Bool("debug", true, "debug test flag")
	flagSet.String("formatter", "text", "formatter test flag")
	flagSet.String("stdoutput", "out", "stdoutput test flag")
	err := StartDebugMode(flagSet)
	expected := defaultDebugFlags()
	expected.debug = true
	expected.formatter = &logrus.TextFormatter{FullTimestamp: true}
	expected.stdout = os.Stdout
	assert.Equal(t, expected, flags)
	assert.Nil(t, err)
	//for other tests we need to make sure to keep debug mode off
	flags.debug = false
}

func TestStartDebugModeNoDebugFlag(t *testing.T) {
	flags = defaultDebugFlags()
	flagSet := new(flag.FlagSet)
	flagSet.String("formatter", "text", "formatter test flag")
	flagSet.String("stdoutput", "out", "stdoutput test flag")
	err := StartDebugMode(flagSet)
	expected := defaultDebugFlags()
	assert.Equal(t, expected, flags)
	assert.NotNil(t, err)
}

func TestStartDebugModeErrorFlags(t *testing.T) {
	flags = defaultDebugFlags()
	flagSet := new(flag.FlagSet)
	flagSet.Bool("debug", true, "debug test flag")
	flagSet.String("formatter", "wrongvalue", "formatter test flag")
	flagSet.String("stdoutput", "wrongvalue", "stdoutput test flag")
	err := StartDebugMode(flagSet)
	expected := defaultDebugFlags()
	assert.Equal(t, expected, flags)
	assert.NotNil(t, err)
	//for other tests we need to make sure to keep debug mode off
	flags.debug = false
}
