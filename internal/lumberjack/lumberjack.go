package lumberjack

import (
	"io"
	"log/slog"
)

var Logger *slog.Logger

func Init(w io.Writer) {
	handler := slog.NewTextHandler(w, nil)
	logger := slog.New(handler)

	Logger = logger
}
