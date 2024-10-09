package logging

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
	"sync"
)

type Logger = *zap.SugaredLogger

type Field = zap.Field

var defaultLogger Logger
var logLevel = zapcore.InfoLevel
var loggerOnce sync.Once

func SetLogLevel(level string) error {
	switch strings.ToLower(level) {
	case "debug":
		logLevel = zapcore.DebugLevel
	case "info":
		logLevel = zapcore.InfoLevel
	case "warn":
		logLevel = zapcore.WarnLevel
	case "error":
		logLevel = zapcore.ErrorLevel
	case "panic":
		logLevel = zapcore.PanicLevel
	case "fatal":
		logLevel = zapcore.FatalLevel
	default:
		return fmt.Errorf("invalid log level %s", level)
	}
	return nil
}

func New(name string, fields ...Field) Logger {
	logger := newLogger(name)

	if len(fields) > 0 {
		var args = make([]any, len(fields))
		for i, field := range fields {
			args[i] = field
		}

		logger = logger.With(args...)
	}
	return logger
}

func DefaultLogger() Logger {
	loggerOnce.Do(func() {
		defaultLogger = newLogger("default")
	})
	return defaultLogger
}

func NewField(key string, value string) Field {
	return zap.String(key, value)
}

func newLogger(name string) Logger {
	return zap.New(zapcore.NewTee(
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(humanEncoderConfig()),
			zapcore.AddSync(os.Stdout),
			logLevel,
		),
	), zap.AddStacktrace(zap.ErrorLevel)).Named(name).Sugar()
}

func encoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func humanEncoderConfig() zapcore.EncoderConfig {
	cfg := encoderConfig()
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	cfg.EncodeDuration = zapcore.StringDurationEncoder
	return cfg
}
