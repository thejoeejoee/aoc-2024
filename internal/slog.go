package internal

import (
	"log/slog"
	"os"
)

func init() {
	l := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelDebug,
	})
	slog.SetDefault(slog.New(l))
}
