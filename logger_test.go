package gtmlp

import (
	"bytes"
	"log/slog"
	"strings"
	"sync"
	"testing"
)

// TestDefaultLogger tests that the default logger is configured correctly
func TestDefaultLogger(t *testing.T) {
	logger := GetLogger()
	if logger == nil {
		t.Fatal("Default logger should not be nil")
	}

	// Default level should be Warn
	if globalLevel != slog.LevelWarn {
		t.Errorf("Expected default level to be Warn, got %v", globalLevel)
	}
}

// TestSetLogger tests setting a custom logger
func TestSetLogger(t *testing.T) {
	// Save original logger
	originalLogger := GetLogger()
	defer SetLogger(originalLogger)

	// Create custom logger with buffer
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	customLogger := slog.New(handler)

	// Set custom logger
	SetLogger(customLogger)

	// Verify it was set
	if GetLogger() != customLogger {
		t.Error("Custom logger was not set correctly")
	}

	// Test logging
	GetLogger().Info("test message", "key", "value")

	output := buf.String()
	if !strings.Contains(output, "test message") {
		t.Errorf("Expected log output to contain 'test message', got: %s", output)
	}
	if !strings.Contains(output, "key=value") {
		t.Errorf("Expected log output to contain 'key=value', got: %s", output)
	}
}

// TestSetLogLevel tests changing the log level
func TestSetLogLevel(t *testing.T) {
	// Save original logger
	originalLogger := GetLogger()
	originalLevel := globalLevel
	defer func() {
		SetLogger(originalLogger)
		globalLevel = originalLevel
	}()

	tests := []struct {
		name  string
		level slog.Level
	}{
		{"Debug", slog.LevelDebug},
		{"Info", slog.LevelInfo},
		{"Warn", slog.LevelWarn},
		{"Error", slog.LevelError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetLogLevel(tt.level)

			if globalLevel != tt.level {
				t.Errorf("Expected level %v, got %v", tt.level, globalLevel)
			}
		})
	}
}

// TestLogLevels tests that different log levels work correctly
func TestLogLevels(t *testing.T) {
	// Save original state
	originalLogger := GetLogger()
	defer SetLogger(originalLogger)

	// Create logger with buffer
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	SetLogger(slog.New(handler))

	// Debug should not appear (level is Info)
	GetLogger().Debug("debug message")
	if strings.Contains(buf.String(), "debug message") {
		t.Error("Debug message should not appear when level is Info")
	}

	// Info should appear
	buf.Reset()
	GetLogger().Info("info message")
	if !strings.Contains(buf.String(), "info message") {
		t.Error("Info message should appear when level is Info")
	}

	// Warn should appear
	buf.Reset()
	GetLogger().Warn("warn message")
	if !strings.Contains(buf.String(), "warn message") {
		t.Error("Warn message should appear when level is Info")
	}

	// Error should appear
	buf.Reset()
	GetLogger().Error("error message")
	if !strings.Contains(buf.String(), "error message") {
		t.Error("Error message should appear when level is Info")
	}
}

// TestLogOutput tests structured logging output
func TestLogOutput(t *testing.T) {
	// Save original state
	originalLogger := GetLogger()
	defer SetLogger(originalLogger)

	// Create logger with buffer
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	SetLogger(slog.New(handler))

	// Log with structured fields
	GetLogger().Info("test message",
		"url", "https://example.com",
		"status", 200,
		"duration_ms", 150)

	output := buf.String()

	// Check all fields are present
	expectedFields := []string{"test message", "url=https://example.com", "status=200", "duration_ms=150"}
	for _, field := range expectedFields {
		if !strings.Contains(output, field) {
			t.Errorf("Expected output to contain '%s', got: %s", field, output)
		}
	}
}

// TestJSONHandler tests JSON output format
func TestJSONHandler(t *testing.T) {
	// Save original state
	originalLogger := GetLogger()
	defer SetLogger(originalLogger)

	// Create JSON logger with buffer
	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	SetLogger(slog.New(handler))

	// Log message
	GetLogger().Info("json test", "key", "value")

	output := buf.String()

	// Should contain JSON structure
	if !strings.Contains(output, `"msg":"json test"`) {
		t.Errorf("Expected JSON output to contain message, got: %s", output)
	}
	if !strings.Contains(output, `"key":"value"`) {
		t.Errorf("Expected JSON output to contain key-value, got: %s", output)
	}
}

// TestConcurrentLogging tests thread-safety of logger
func TestConcurrentLogging(t *testing.T) {
	// Save original state
	originalLogger := GetLogger()
	defer SetLogger(originalLogger)

	// Create logger with buffer
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	SetLogger(slog.New(handler))

	// Run concurrent logging
	var wg sync.WaitGroup
	numGoroutines := 100
	numLogsPerGoroutine := 10

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numLogsPerGoroutine; j++ {
				GetLogger().Info("concurrent log", "goroutine", id, "iteration", j)
			}
		}(i)
	}

	wg.Wait()

	// Just verify it didn't crash and produced output
	output := buf.String()
	if len(output) == 0 {
		t.Error("Expected log output from concurrent logging")
	}
}

// TestGetLogger tests the GetLogger function
func TestGetLogger(t *testing.T) {
	logger1 := GetLogger()
	logger2 := GetLogger()

	// Should return the same instance
	if logger1 != logger2 {
		t.Error("GetLogger should return the same instance")
	}
}

// TestSetLogLevel_ResetsToDefault tests that SetLogLevel creates new default handler
func TestSetLogLevel_ResetsToDefault(t *testing.T) {
	// Save original state
	originalLogger := GetLogger()
	defer SetLogger(originalLogger)

	// Set custom JSON handler with buffer
	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	SetLogger(slog.New(handler))

	// Verify custom handler works
	GetLogger().Info("before level change")
	if !strings.Contains(buf.String(), "before level change") {
		t.Error("Custom handler should work before SetLogLevel")
	}

	// Change level - this should reset to default handler (stderr, not buffer)
	SetLogLevel(slog.LevelDebug)

	// Clear buffer
	buf.Reset()

	// Log something - should NOT go to buffer anymore (goes to stderr)
	GetLogger().Debug("after level change")

	output := buf.String()

	// Buffer should be empty because SetLogLevel created a new default handler
	if output != "" {
		t.Errorf("Expected empty buffer after SetLogLevel (resets to default handler), got: %s", output)
	}
}
