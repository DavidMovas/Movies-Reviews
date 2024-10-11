package log

import (
	"context"
	"log/slog"
	"os"
)

const (
	contextLoggerKey = "logger"
)

type leveler struct {
	level slog.Level
}

func (l leveler) Level() slog.Level {
	return l.level
}

func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, contextLoggerKey, logger)
}

func FromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(contextLoggerKey).(*slog.Logger); ok {
		return logger
	}

	return slog.Default()
}

func SetupLogger(isLocal bool, level string) (*slog.Logger, error) {
	l, err := newLevelFromString(level)
	if err != nil {
		return nil, err
	}

	opts := &slog.HandlerOptions{Level: l}

	var handler slog.Handler
	if isLocal {
		handler = slog.NewTextHandler(os.Stdout, opts)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	return slog.New(handler), nil
}

func newLevelFromString(level string) (slog.Level, error) {
	switch level {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return slog.LevelInfo, nil
	}
}
