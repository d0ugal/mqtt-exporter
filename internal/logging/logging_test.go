package logging

import (
	"bytes"
	"log/slog"
	"os"
	"strings"
	"testing"

	"github.com/d0ugal/mqtt-exporter/internal/config"
)

func TestConfigure_JSONFormat(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	defer func() { os.Stdout = oldStdout }()

	// Configure logging
	cfg := &config.LoggingConfig{
		Level:  "debug",
		Format: "json",
	}
	Configure(cfg)

	// Test logging
	slog.Info("test message", "key", "value")

	// Close pipe and read output
	if err := w.Close(); err != nil {
		t.Fatalf("Failed to close pipe: %v", err)
	}

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("Failed to read from pipe: %v", err)
	}

	output := buf.String()

	// Verify JSON format
	if !strings.Contains(output, `"msg":"test message"`) {
		t.Errorf("Expected JSON log output, got: %s", output)
	}

	if !strings.Contains(output, `"key":"value"`) {
		t.Errorf("Expected JSON log output with key-value pair, got: %s", output)
	}
}

func TestConfigure_TextFormat(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	defer func() { os.Stdout = oldStdout }()

	// Configure logging
	cfg := &config.LoggingConfig{
		Level:  "info",
		Format: "text",
	}
	Configure(cfg)

	// Test logging
	slog.Info("test message", "key", "value")

	// Close pipe and read output
	if err := w.Close(); err != nil {
		t.Fatalf("Failed to close pipe: %v", err)
	}

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("Failed to read from pipe: %v", err)
	}

	output := buf.String()

	// Verify text format
	if !strings.Contains(output, "test message") {
		t.Errorf("Expected text log output, got: %s", output)
	}

	if !strings.Contains(output, "key=value") {
		t.Errorf("Expected text log output with key-value pair, got: %s", output)
	}
}

func TestConfigure_DefaultLevel(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	defer func() { os.Stdout = oldStdout }()

	// Configure logging with empty level (should default to info)
	cfg := &config.LoggingConfig{
		Level:  "",
		Format: "text",
	}
	Configure(cfg)

	// Test debug logging (should not appear with info level)
	slog.Debug("debug message")
	slog.Info("info message")

	// Close pipe and read output
	if err := w.Close(); err != nil {
		t.Fatalf("Failed to close pipe: %v", err)
	}

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("Failed to read from pipe: %v", err)
	}

	output := buf.String()

	// Verify debug message is not logged, but info is
	if strings.Contains(output, "debug message") {
		t.Errorf("Debug message should not be logged with info level, got: %s", output)
	}

	if !strings.Contains(output, "info message") {
		t.Errorf("Info message should be logged, got: %s", output)
	}
}

func TestConfigure_DefaultFormat(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	defer func() { os.Stdout = oldStdout }()

	// Configure logging with empty format (should default to text)
	cfg := &config.LoggingConfig{
		Level:  "info",
		Format: "",
	}
	Configure(cfg)

	// Test logging
	slog.Info("test message")

	// Close pipe and read output
	if err := w.Close(); err != nil {
		t.Fatalf("Failed to close pipe: %v", err)
	}

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("Failed to read from pipe: %v", err)
	}

	output := buf.String()

	// Verify text format (default) - check for text structure
	if !strings.Contains(output, "msg=") || !strings.Contains(output, "level=") {
		t.Errorf("Expected text log output (default), got: %s", output)
	}
}

func TestConfigure_AllLevels(t *testing.T) {
	levels := []string{"debug", "info", "warn", "error"}

	for _, level := range levels {
		t.Run(level, func(t *testing.T) {
			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			defer func() { os.Stdout = oldStdout }()

			// Configure logging
			cfg := &config.LoggingConfig{
				Level:  level,
				Format: "text",
			}
			Configure(cfg)

			// Test logging at the configured level
			switch level {
			case "debug":
				slog.Debug("debug message")
			case "info":
				slog.Info("info message")
			case "warn":
				slog.Warn("warn message")
			case "error":
				slog.Error("error message")
			}

			// Close pipe and read output
			if err := w.Close(); err != nil {
				t.Fatalf("Failed to close pipe: %v", err)
			}

			var buf bytes.Buffer
			if _, err := buf.ReadFrom(r); err != nil {
				t.Fatalf("Failed to read from pipe: %v", err)
			}

			output := buf.String()

			// Verify the message is logged
			expectedMessage := level + " message"
			if !strings.Contains(output, expectedMessage) {
				t.Errorf("Expected %s message to be logged, got: %s", level, output)
			}
		})
	}
}
