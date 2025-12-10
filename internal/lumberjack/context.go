package lumberjack

import (
	"context"
	"log/slog"
)

type ctxKey string

const loggerKey ctxKey = "logger"

func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func FromContext(ctx context.Context) *slog.Logger {
	if ctx == nil {
		return Logger
	}
	if l, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
		return l
	}
	return Logger
}

func NewRequestLogger(requestID string) *slog.Logger {
	return Logger.With("request_id", requestID)
}
