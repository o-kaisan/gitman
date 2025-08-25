package common

import (
	"log/slog"
	"os"
	"strings"
)

var Log *slog.Logger

func SetupGlobalLogger(debug bool) {
	level := slog.LevelInfo
	if strings.ToLower(GetEnvWithString("GITMAN_DEBUG", "false")) == "true" || debug {
		level = slog.LevelDebug
	}

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})

	logger := slog.New(handler)

	slog.SetDefault(logger)
}
