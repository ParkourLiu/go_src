package log

import (
	"os"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

const (
	LevelDebug = 1 << iota
	LevelInfo
	LevelWarn
	LevelError
)

type Logger struct {
	logger log.Logger
}

var defaultLogger *Logger

func init() {
	var myTimestamp log.Valuer = func() interface{} { return time.Now().String() }
	logger := log.NewLogfmtLogger(os.Stdout)
	logger = log.With(logger, "ts", myTimestamp)
	logger = log.With(logger, "caller", log.Caller(7))
	defaultLogger = &Logger{logger: logger}
}

func SetDefaultLogLevel(level int) {
	defaultLogger.SetLogLevel(level)
}

func With(keyvals ...interface{}) {
	defaultLogger.With(keyvals...)
}

func GetDefaultLogger() *Logger {
	return defaultLogger
}

func Debug(keyvals ...interface{}) {
	defaultLogger.Debug(keyvals...)
}

func Info(keyvals ...interface{}) {
	defaultLogger.Info(keyvals...)
}

func Warn(keyvals ...interface{}) {
	defaultLogger.Warn(keyvals...)
}

func Error(keyvals ...interface{}) {
	defaultLogger.Error(keyvals...)
}

func NewLogger() *Logger {
	var myTimestamp log.Valuer = func() interface{} { return time.Now().String() }
	logger := log.NewLogfmtLogger(os.Stdout)
	logger = log.With(logger, "ts", myTimestamp)
	logger = log.With(logger, "caller", log.Caller(6))
	return &Logger{logger: logger}
}

func (l *Logger) GetDefaultKitLogger() *log.Logger {
	return &l.logger
}

func (l *Logger) With(keyvals ...interface{}) {
	l.logger = log.With(l.logger, keyvals...)
}

func (l *Logger) Debug(keyvals ...interface{}) {
	level.Debug(l.logger).Log(keyvals...)
}

func (l *Logger) Info(keyvals ...interface{}) {
	level.Info(l.logger).Log(keyvals...)
	//log.With(l.logger, level.Key(), level.InfoValue()).Log(keyvals...)
}

func (l *Logger) Warn(keyvals ...interface{}) {
	level.Warn(l.logger).Log(keyvals...)
}

func (l *Logger) Error(keyvals ...interface{}) {
	level.Error(l.logger).Log(keyvals...)
}

func (l *Logger) SetLogLevel(levelval int) {
	switch levelval {
	case LevelDebug:
		l.logger = level.NewFilter(l.logger, level.AllowDebug())
	case LevelInfo:
		l.logger = level.NewFilter(l.logger, level.AllowInfo())
	case LevelWarn:
		l.logger = level.NewFilter(l.logger, level.AllowWarn())
	case LevelError:
		l.logger = level.NewFilter(l.logger, level.AllowError())
	default:
		l.logger = level.NewFilter(l.logger, level.AllowAll())
	}
}
