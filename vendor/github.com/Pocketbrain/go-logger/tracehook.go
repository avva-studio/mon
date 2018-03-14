package plog

import (
	"path"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
)

// TraceHook struct for satisfying interface
type TraceHook struct {
}

// Levels where the trace should be applicable for
func (hook TraceHook) Levels() []log.Level {
	return log.AllLevels
}

// Fire will run the function when an entry will be made
func (hook TraceHook) Fire(entry *log.Entry) error {
	if pc, file, line, ok := runtime.Caller(6); ok {
		funcName := runtime.FuncForPC(pc).Name()

		entry.Data["file"] = path.Base(file)
		fName := strings.Split(path.Base(funcName), ".")
		entry.Data["function"] = fName[0]
		entry.Data["line"] = line
	}

	return nil
}
