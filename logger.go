package gtmlp

import (
	"log/slog"
	"os"
	"sync"
)

var (
	globalLogger *slog.Logger
	globalLevel  slog.Level
	loggerMutex  sync.RWMutex
)

func init() {
	// Default: Warn level, text handler to stderr (production-safe)
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelWarn,
	})
	globalLogger = slog.New(handler)
	globalLevel = slog.LevelWarn
}

// SetLogger configures the global logger
// Use this to customize the logger handler (JSON vs Text, output destination, etc.)
//
// Example:
//
//	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
//	    Level: slog.LevelInfo,
//	})
//	gtmlp.SetLogger(slog.New(handler))
func SetLogger(logger *slog.Logger) {
	loggerMutex.Lock()
	defer loggerMutex.Unlock()
	globalLogger = logger
}

// SetLogLevel changes the global log level by creating a new default handler
// Available levels: slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError
//
// Default: slog.LevelWarn (production-safe)
//
// Note: This recreates the handler with default settings (TextHandler to stderr).
// If you're using a custom handler (custom writer, JSON format, etc.), use SetLogger instead.
//
// Example:
//
//	// Development: enable Info logs
//	gtmlp.SetLogLevel(slog.LevelInfo)
//
//	// Troubleshooting: enable Debug logs
//	gtmlp.SetLogLevel(slog.LevelDebug)
//
//	// Production: use default Warn level (no call needed)
//
//	// For custom handlers, use SetLogger:
//	handler := slog.NewJSONHandler(myWriter, &slog.HandlerOptions{Level: slog.LevelDebug})
//	gtmlp.SetLogger(slog.New(handler))
func SetLogLevel(level slog.Level) {
	loggerMutex.Lock()
	defer loggerMutex.Unlock()

	globalLevel = level

	// Create new default handler with updated level
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: level,
	})

	globalLogger = slog.New(handler)
}

// GetLogger returns the current global logger
// Useful for testing and debugging
func GetLogger() *slog.Logger {
	loggerMutex.RLock()
	defer loggerMutex.RUnlock()
	return globalLogger
}

// getLogger returns the global logger for internal use (thread-safe)
func getLogger() *slog.Logger {
	loggerMutex.RLock()
	defer loggerMutex.RUnlock()
	return globalLogger
}
