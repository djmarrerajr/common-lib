package utils_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/djmarrerajr/common-lib/utils"
)

const (
	loggerName   = "unit-test-logger"
	loggerMsg    = "test-message"
	loggerKey    = "test_key"
	loggerCtxKey = "ctx_key"
	loggerValue  = "test-value"
	loggerFormat = "%s - %s = %s"
)

type LoggerTestSuite struct {
	suite.Suite

	ctx        context.Context
	stdout     []byte
	readStream *os.File
	origStdout *os.File
}

func (l *LoggerTestSuite) SetupTest() {
	l.ctx = context.Background()

	// we need to redirect stdout because we are wrapping zap
	l.stdout = make([]byte, 2048)

	readStream, writeStream, err := os.Pipe()
	l.readStream = readStream
	l.NoError(err, "unable to create the pipe needed for stdout redirection")

	origStdout := os.Stdout
	l.origStdout = origStdout
	os.Stdout = writeStream
}

func (l *LoggerTestSuite) TeardownTest() {
	os.Stdout = l.origStdout
}

type loggerFunc func(string, ...interface{})
type logMsg struct {
	Level           string  `json:"level"`
	Logger          string  `json:"logger"`
	Message         string  `json:"msg"`
	Value           string  `json:"test_key,omitempty"`
	CtxVal          string  `json:"ctx_key,omitempty"`
	ErrorMessage    *string `json:"error.message,omitempty"`
	ErrorStackTrace *string `json:"error.stack,omitempty"`
}

func (l *LoggerTestSuite) TestLogger_LogMessagesWithKeyValuePairs() {
	logger := utils.NewLogger("DEBUG").Named(loggerName)

	ctx := utils.AddMapToContext(context.Background(), map[string]interface{}{
		loggerCtxKey: loggerValue,
	})

	loggerFuncs := map[string]loggerFunc{
		"debug": logger.WithCtx(ctx).Debugw,
		"info":  logger.WithCtx(ctx).Infow,
		"warn":  logger.WithCtx(ctx).Warnw,
		"error": logger.WithCtx(ctx).Errorw,
	}

	// log a message at each level and verify the result...
	for level, levelFunc := range loggerFuncs {
		levelFunc(loggerMsg, loggerKey, loggerValue)

		entry := l.getLogEntry()

		l.Equalf(loggerName, entry.Logger, "invalid logger name")
		l.Equalf(level, entry.Level, "invalid level")
		l.Equalf(loggerMsg, entry.Message, "invalid message")
		l.Equalf(loggerValue, entry.Value, "invalid key/value")
		l.Equalf(loggerValue, entry.CtxVal, "invalid context value")
	}
}

func (l *LoggerTestSuite) TestLogger_LogFormattedAtDifferentLevels() {
	logger := utils.NewLogger("DEBUG").Named(loggerName)

	ctx := utils.AddMapToContext(context.Background(), map[string]interface{}{
		loggerCtxKey: loggerValue,
	})

	loggerFuncs := map[string]loggerFunc{
		"debug": logger.WithCtx(ctx).Debugf,
		"info":  logger.WithCtx(ctx).Infof,
		"warn":  logger.WithCtx(ctx).Warnf,
		"error": logger.WithCtx(ctx).Errorf,
	}

	// log a message at each level and verify the result...
	for level, levelFunc := range loggerFuncs {
		levelFunc(loggerFormat, loggerMsg, loggerKey, loggerValue)

		entry := l.getLogEntry()

		l.Equalf(loggerName, entry.Logger, "invalid logger name")
		l.Equalf(level, entry.Level, "invalid level")
		l.Equalf(fmt.Sprintf(loggerFormat, loggerMsg, loggerKey, loggerValue), entry.Message, "invalid message")
		l.Equalf(loggerValue, entry.CtxVal, "invalid context value")
	}
}

func (l *LoggerTestSuite) TestLogger_Error_AddsErrorMessage() {
	err := fmt.Errorf("a regular error with no stack trace")

	logger := utils.NewLogger("DEBUG").Named(loggerName)

	logger.Error("foo", err)

	entry := l.getLogEntry()

	l.Equal(*entry.ErrorMessage, err.Error())
	l.Nil(entry.ErrorStackTrace)
}

func (l *LoggerTestSuite) getLogEntry() logMsg {
	entry := logMsg{}

	bufSize, err := l.readStream.Read(l.stdout)
	l.NoError(err, "unable to read from stdout redirection pipe")

	err = json.Unmarshal(l.stdout[:bufSize], &entry)
	l.NoError(err, "unable to unmarshal logger output")

	return entry
}

func TestLogger(t *testing.T) {
	suite.Run(t, new(LoggerTestSuite))
}
