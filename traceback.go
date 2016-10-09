package traceback

import (
	"github.com/Sirupsen/logrus"
	"gopkg.in/sensorbee/sensorbee.v0/core"
	"gopkg.in/sensorbee/sensorbee.v0/data"
	"runtime"
)

type TracebackFormatter struct {
	parent logrus.Formatter
	level  logrus.Level
}

func (f *TracebackFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// Only print stacktrace if the level of the entry is higher than the specified log level.
	if entry.Level <= f.level {
		traceback := []map[string]interface{}{}
		for level := 2; ; level += 1 {
			pc, file, line, ok := runtime.Caller(level)
			if !ok {
				break
			}
			traceback = append(traceback, map[string]interface{}{
				"func": runtime.FuncForPC(pc).Name(),
				"file": file,
				"line": line,
			})
		}
		entry.Data["traceback"] = traceback
	}
	return f.parent.Format(entry)
}

type tracebackUDS struct{}

func NewTracebackUDS(ctx *core.Context, params data.Map) (core.SharedState, error) {
	level := logrus.WarnLevel
	if v, ok := params["level"]; ok {
		levelStr, err := data.ToString(v)
		if err != nil {
			return nil, err
		}

		level, err = logrus.ParseLevel(levelStr)
		if err != nil {
			return nil, err
		}
	}

	entry := ctx.Log()
	logger := entry.Logger
	logger.Formatter = &TracebackFormatter{
		parent: logger.Formatter,
		level:  level,
	}

	entry.Infof("enabled traceback (min log level: %s)", level.String())

	return &tracebackUDS{}, nil
}

func (t *tracebackUDS) Terminate(ctx *core.Context) error {
	return nil
}
