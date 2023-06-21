package cockroach

import (
	"github.com/djmarrerajr/common-lib/utils"
)

type LogWriter struct {
	utils.Logger
}

func (l LogWriter) Printf(msg string, data ...interface{}) {
	var isError bool

	for _, item := range data {
		switch item.(type) {
		case error:
			isError = true
		}
	}

	if isError {
		l.Errorf(msg, data...)
	} else {
		l.Infof(msg, data...)
	}
}

func NewGormLogger(logger utils.Logger) *LogWriter {
	return &LogWriter{logger}
}
