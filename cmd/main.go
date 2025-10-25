package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/d0ugal/mqtt-exporter/internal/collectors"
	"github.com/d0ugal/mqtt-exporter/internal/config"
	"github.com/d0ugal/mqtt-exporter/internal/metrics"
	"github.com/d0ugal/mqtt-exporter/internal/version"
	"github.com/d0ugal/promexporter/app"
	"github.com/d0ugal/promexporter/logging"
	promexporter_metrics "github.com/d0ugal/promexporter/metrics"
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
	// Parse command line flags
	var showVersion bool
	flag.BoolVar(&showVersion, "version", false, "Show version information")
	flag.BoolVar(&showVersion, "v", false, "Show version information")

	var (
		configPath    string
		configFromEnv bool
	)

	flag.StringVar(&configPath, "config", "config.yaml", "Path to configuration file")
	flag.BoolVar(&configFromEnv, "config-from-env", false, "Load configuration from environment variables only")
	flag.Parse()

	// Show version if requested
	if showVersion {
		fmt.Printf("mqtt-exporter %s\n", version.Version)
		fmt.Printf("Commit: %s\n", version.Commit)
		fmt.Printf("Build Date: %s\n", version.BuildDate)
		os.Exit(0)
	}

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

	// Configure logging using promexporter
	logging.Configure(&logging.Config{
		Level:  cfg.Logging.Level,
		Format: cfg.Logging.Format,
	})

	// Initialize metrics registry using promexporter
	metricsRegistry := promexporter_metrics.NewRegistry("mqtt_exporter_info")

	// Set version info metric with mqtt-exporter version information
	metricsRegistry.VersionInfo.WithLabelValues(version.Version, version.Commit, version.BuildDate).Set(1)

	// Add custom metrics to the registry
	mqttRegistry := metrics.NewMQTTRegistry(metricsRegistry)

	// Create collector
	mqttCollector := collectors.NewMQTTCollector(cfg, mqttRegistry)

	// Create and run application using promexporter
	application := app.New("mqtt-exporter").
		WithConfig(&cfg.BaseConfig).
		WithMetrics(metricsRegistry).
		WithCollector(mqttCollector).
		Build()

	if err := application.Run(); err != nil {
		slog.Error("Application failed", "error", err)
		os.Exit(1)
	}
}
