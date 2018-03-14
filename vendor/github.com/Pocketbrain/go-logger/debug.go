package plog

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

type debugFlags struct {
	debug     bool
	formatter logrus.Formatter
	stdout    *os.File
}

var flags = defaultDebugFlags()

// StartDebugMode will start both base logger and plog in a debug mode if
// the flag "debug" is set to true. In debug mode logger is set to the debug level
// and setting level, formatter or stdoutput won't have any affect.
// Optional flags:
// - formatter - allowed values "json"(default) and "text"
// - stdoutput - allowed values "err"(default) and "out"
func StartDebugMode(fs *flag.FlagSet) error {
	flags = defaultDebugFlags()
	dMode, err := fs.GetBool("debug")
	if err != nil {
		return err
	}
	if !dMode {
		return nil
	}
	//set base logger default options
	SetStderr()
	SetJSONFormatter()
	SetLevel(DebugLevel)

	err = setFormatterFlag(fs)
	if err != nil {
		return err
	}

	err = setStdoutputFlag(fs)
	if err != nil {
		return err
	}

	//we lock the settings, when debug is true plog neither base logger options can be changed
	flags.debug = true
	return nil
}

func setFormatterFlag(fs *flag.FlagSet) error {
	f, err := fs.GetString("formatter")
	//as formatter flag is optional we do not report the error
	if err != nil {
		return nil
	}
	switch f {
	case "text":
		flags.formatter = &logrus.TextFormatter{FullTimestamp: true}
		SetTextFormatter()
	case "json":
		//already set as default value
	default:
		return fmt.Errorf("invalid formatter flag value %s", f)
	}
	return nil
}

func setStdoutputFlag(fs *flag.FlagSet) error {
	s, err := fs.GetString("stdoutput")
	//as stdoutput flag is optional we do not report the error
	if err != nil {
		return nil
	}
	switch s {
	case "out":
		flags.stdout = os.Stdout
		SetStdout()
	case "err":
		//already set as default value
	default:
		return fmt.Errorf("invalid stdoutput flag value %s", s)
	}
	return nil
}

func defaultDebugFlags() *debugFlags {
	return &debugFlags{
		formatter: &logrus.JSONFormatter{},
		stdout:    os.Stderr,
	}
}

func inDebugMode() bool {
	return flags.debug
}
