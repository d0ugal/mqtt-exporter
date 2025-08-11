package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/d0ugal/mqtt-exporter/internal/collectors"
	"github.com/d0ugal/mqtt-exporter/internal/config"
	"github.com/d0ugal/mqtt-exporter/internal/logging"
	"github.com/d0ugal/mqtt-exporter/internal/metrics"
	"github.com/d0ugal/mqtt-exporter/internal/server"
)

// hasEnvironmentVariables checks if any MQTT_EXPORTER_* environment variables are set
func hasEnvironmentVariables() bool {
	envVars := []string{
		"MQTT_EXPORTER_SERVER_HOST",
		"MQTT_EXPORTER_SERVER_PORT",
		"MQTT_EXPORTER_LOG_LEVEL",
		"MQTT_EXPORTER_LOG_FORMAT",
		"MQTT_EXPORTER_METRICS_DEFAULT_INTERVAL",
		"MQTT_EXPORTER_MQTT_BROKER",
		"MQTT_EXPORTER_MQTT_CLIENT_ID",
		"MQTT_EXPORTER_MQTT_USERNAME",
		"MQTT_EXPORTER_MQTT_PASSWORD",
		"MQTT_EXPORTER_MQTT_TOPICS",
		"MQTT_EXPORTER_MQTT_QOS",
		"MQTT_EXPORTER_MQTT_CLEAN_SESSION",
		"MQTT_EXPORTER_MQTT_KEEP_ALIVE",
		"MQTT_EXPORTER_MQTT_CONNECT_TIMEOUT",
	}

	for _, envVar := range envVars {
		if os.Getenv(envVar) != "" {
			return true
		}
	}

	return false
}

func main() {
	var (
		configPath    string
		configFromEnv bool
	)

	flag.StringVar(&configPath, "config", "config.yaml", "Path to configuration file")
	flag.BoolVar(&configFromEnv, "config-from-env", false, "Load configuration from environment variables only")
	flag.Parse()

	// Use environment variable if config flag is not provided
	if configPath == "config.yaml" && !configFromEnv {
		if envConfig := os.Getenv("CONFIG_PATH"); envConfig != "" {
			configPath = envConfig
		}
	}

	// Check if we should use environment-only configuration
	if !configFromEnv {
		// Check explicit flag first
		if os.Getenv("MQTT_EXPORTER_CONFIG_FROM_ENV") == "true" {
			configFromEnv = true
		} else if hasEnvironmentVariables() {
			// Auto-detect environment variables and use them
			configFromEnv = true
		}
	}

	// Load configuration
	cfg, err := config.LoadConfig(configPath, configFromEnv)
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Configure logging
	logging.Configure(&cfg.Logging)

	// Initialize metrics
	metricsRegistry := metrics.NewRegistry()

	// Create collectors
	mqttCollector := collectors.NewMQTTCollector(cfg, metricsRegistry)

	// Start collectors
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mqttCollector.Start(ctx)

	// Create and start server
	srv := server.New(cfg, metricsRegistry)

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		slog.Info("Shutting down gracefully...")
		cancel()

		if err := srv.Shutdown(); err != nil {
			slog.Error("Failed to shutdown server gracefully", "error", err)
		}
	}()

	// Start server
	if err := srv.Start(); err != nil {
		slog.Error("Server failed", "error", err)
		os.Exit(1)
	}
}
