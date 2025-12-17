package woodpecker

import (
	"io"
	"log/slog"
)

var Logger *slog.Logger

func Init(w io.Writer) {
	handler := slog.NewJSONHandler(w, nil)
	logger := slog.New(handler)

	Logger = logger
}
