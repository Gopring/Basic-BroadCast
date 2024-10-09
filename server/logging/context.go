package logging

import "context"

type loggerKey struct{}

func With(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

func From(ctx context.Context) Logger {
	if ctx == nil {
		return defaultLogger
	}

	logger, ok := ctx.Value(loggerKey{}).(Logger)
	if !ok {
		return defaultLogger
	}
	return logger
}
