package main

import (
	"os"
	"testing"
)

func TestHasEnvironmentVariables(t *testing.T) {
	// Test with no environment variables set
	if hasEnvironmentVariables() {
		t.Error("hasEnvironmentVariables() should return false when no env vars are set")
	}

	// Test with a single environment variable set
	if err := os.Setenv("MQTT_EXPORTER_MQTT_BROKER", "localhost:1883"); err != nil {
		t.Fatalf("Failed to set environment variable: %v", err)
	}

	defer func() {
		if err := os.Unsetenv("MQTT_EXPORTER_MQTT_BROKER"); err != nil {
			t.Logf("Failed to unset environment variable: %v", err)
		}
	}()

	if !hasEnvironmentVariables() {
		t.Error("hasEnvironmentVariables() should return true when MQTT_EXPORTER_MQTT_BROKER is set")
	}

	// Test with multiple environment variables set
	if err := os.Setenv("MQTT_EXPORTER_SERVER_HOST", "127.0.0.1"); err != nil {
		t.Fatalf("Failed to set environment variable: %v", err)
	}

	defer func() {
		if err := os.Unsetenv("MQTT_EXPORTER_SERVER_HOST"); err != nil {
			t.Logf("Failed to unset environment variable: %v", err)
		}
	}()

	if !hasEnvironmentVariables() {
		t.Error("hasEnvironmentVariables() should return true when multiple env vars are set")
	}

	// Test with non-MQTT_EXPORTER environment variable (should not affect result)
	if err := os.Setenv("OTHER_VAR", "value"); err != nil {
		t.Fatalf("Failed to set environment variable: %v", err)
	}

	defer func() {
		if err := os.Unsetenv("OTHER_VAR"); err != nil {
			t.Logf("Failed to unset environment variable: %v", err)
		}
	}()

	if !hasEnvironmentVariables() {
		t.Error("hasEnvironmentVariables() should return true when MQTT_EXPORTER_ vars are set, regardless of other vars")
	}
}
