package lumberjack

import (
	"io"
	"log/slog"
)

var Logger *slog.Logger

type ContextHandler struct {
	slog.Handler
}

func Init(w io.Writer) {
	handler := slog.NewTextHandler(w, nil)
	logger := slog.New(handler)

	Logger = logger
}
