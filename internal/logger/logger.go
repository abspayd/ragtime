package logger

import (
	"io"
	"log/slog"
)

var (
	Log *slog.Logger
)

func New(w io.Writer, verbose bool) {
	level := slog.LevelInfo
	if verbose {
		level = slog.LevelDebug
	}
	handler := slog.NewTextHandler(w, &slog.HandlerOptions{
		Level: level,
	})
	Log = slog.New(handler)
}
