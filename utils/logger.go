package utils

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var currLogLevel zapcore.Level

const LogLevelEnvKey = "LOG_LEVEL"

// Basic interface for a Logger
type Logger interface {
	Named(string) Logger
	WithCtx(context.Context) Logger
	Sync() error
	ToggleDebug()

	Debugw(string, ...interface{})
	Infow(string, ...interface{})
	Warnw(string, ...interface{})
	Errorw(string, ...interface{})
	Error(string, error, ...interface{})

	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Warnf(string, ...interface{})
	Errorf(string, ...interface{})
	Fatalf(string, ...interface{})
}

// Create a new logger at the speficied level
func NewLogger(logLevel string) *ctxLogger {
	return newLogger(logLevel)
}

// Create a new logger pulling the level from the environment
func NewLoggerFromEnv() *ctxLogger {
	return newLogger(os.Getenv(LogLevelEnvKey))
}

// Helper function that creates and returns an instance of a logger
func newLogger(logLevel string) *ctxLogger {
	lvl, err := zapcore.ParseLevel(logLevel)
	if err != nil {
		panic("unsupported log level: " + logLevel)
	}

	currLogLevel = lvl

	clog := &ctxLogger{
		origLvl: lvl,
		ctxMap:  make(map[string]interface{}),
	}

	lvlFunc := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l >= currLogLevel
	})

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		os.Stdout,
		lvlFunc,
	)

	logger := zap.New(core)
	logger = logger.WithOptions(zap.AddCaller(), zap.AddCallerSkip(1))

	clog.logger = logger.Sugar()

	return clog
}

// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=
// Custom Logger that includes contextual values
// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=
type ctxLogger struct {
	origLvl zapcore.Level
	logger  *zap.SugaredLogger
	ctxMap  map[string]interface{}
}

func (l *ctxLogger) ToggleDebug() {
	var state string

	if currLogLevel == zapcore.DebugLevel {
		currLogLevel = l.origLvl
		state = "DISABLED"
	} else {
		currLogLevel = zapcore.DebugLevel
		state = "ENABLED"
	}

	l.logger.Infof("DEBUG logging has been %s", state)
}

// Create a new Named logger
func (l *ctxLogger) Named(name string) Logger {
	logger := l.logger.Named(name)

	return &ctxLogger{
		logger: logger,
		ctxMap: l.ctxMap,
	}
}

func (l *ctxLogger) Sync() error {
	return l.logger.Sync()
}

// Enrich our logger with a context
func (l *ctxLogger) WithCtx(ctx context.Context) Logger {
	if ctx == nil {
		ctx = context.Background()
	}

	newLogger := ctxLogger{
		logger: l.logger,
		ctxMap: make(map[string]interface{}),
	}

	newLogger.updateTrackedValues(ctx)

	return &newLogger
}

func (l *ctxLogger) Debugw(msg string, kvPairs ...interface{}) {
	l.logger.With(l.fields()...).Debugw(msg, kvPairs...)
}

func (l *ctxLogger) Infow(msg string, kvPairs ...interface{}) {
	l.logger.With(l.fields()...).Infow(msg, kvPairs...)
}

func (l *ctxLogger) Warnw(msg string, kvPairs ...interface{}) {
	l.logger.With(l.fields()...).Warnw(msg, kvPairs...)
}

func (l *ctxLogger) Errorw(msg string, kvPairs ...interface{}) {
	l.logger.With(l.fields()...).Errorw(msg, kvPairs...)
}

func (l *ctxLogger) Error(msg string, err error, kvPairs ...interface{}) {
	if err != nil {
		kvPairs = append(kvPairs, "error.message", err.Error())
	}

	l.logger.With(l.fields()...).Errorw(msg, kvPairs...)
}

func (l *ctxLogger) Debugf(msg string, kvPairs ...interface{}) {
	l.logger.With(l.fields()...).Debugf(msg, kvPairs...)
}

func (l *ctxLogger) Infof(format string, kvPairs ...interface{}) {
	l.logger.With(l.fields()...).Infof(format, kvPairs...)
}

func (l *ctxLogger) Warnf(format string, kvPairs ...interface{}) {
	l.logger.With(l.fields()...).Warnf(format, kvPairs...)
}

func (l *ctxLogger) Errorf(format string, kvPairs ...interface{}) {
	l.logger.With(l.fields()...).Errorf(format, kvPairs...)
}

func (l *ctxLogger) Fatalf(format string, kvPairs ...interface{}) {
	l.logger.With(l.fields()...).Fatalf(format, kvPairs...)
}

// Extract the values from our context map so they can be logged
func (l *ctxLogger) fields() (kvPairs []interface{}) {
	for k, v := range l.ctxMap {
		kvPairs = append(kvPairs, k, v)
	}

	return
}

func (l *ctxLogger) updateTrackedValues(ctx context.Context) {
	l.updateContextMap(ctx)
}

func (l *ctxLogger) updateContextMap(ctx context.Context) {
	fieldMap := GetFieldMapFromContext(ctx)

	for k, v := range fieldMap {
		l.ctxMap[k] = v
	}
}
