// Package logger provides structured logging setup using Go's log/slog.
// It configures a global slog logger with either a text (dev) or JSON (prod) handler.
package logger

import (
	"log/slog"
	"os"
)

// Init initialises the global slog logger. The mode parameter selects the
// output format and minimum log level:
//   - "dev"  → text handler, debug level (human-friendly console output)
//   - "prod" → JSON handler, info level (machine-parseable structured logs)
//   - "" or any other value → treated as "dev"
func Init(mode string) {
	var handler slog.Handler

	switch mode {
	case "prod", "production", "release":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	default:
		// dev / any unrecognised value → pretty text with debug level
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	}

	slog.SetDefault(slog.New(handler))
}
