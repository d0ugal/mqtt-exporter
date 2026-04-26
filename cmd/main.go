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
	flag.BoolVar(&configFromEnv, "config-from-env", false, "Deprecated: env vars are always applied; this flag is a no-op")
	flag.Parse()

	// Show version if requested
	if showVersion {
		fmt.Printf("mqtt-exporter %s\n", version.Version)
		fmt.Printf("Commit: %s\n", version.Commit)
		fmt.Printf("Build Date: %s\n", version.BuildDate)
		os.Exit(0)
	}

	if configPath == "config.yaml" {
		if envConfig := os.Getenv("CONFIG_PATH"); envConfig != "" {
			configPath = envConfig
		}
	}

	if configFromEnv || os.Getenv("MQTT_EXPORTER_CONFIG_FROM_ENV") == "true" {
		fmt.Fprintln(os.Stderr, "Warning: --config-from-env / MQTT_EXPORTER_CONFIG_FROM_ENV is deprecated and has no effect. Env vars are always applied on top of yaml config.")
	}

	cfg, err := config.LoadConfig(configPath)
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

	// Add custom metrics to the registry
	mqttRegistry := metrics.NewMQTTRegistry(metricsRegistry)

	// Create and run application using promexporter
	application := app.New("MQTT Exporter").
		WithConfig(&cfg.BaseConfig).
		WithMetrics(metricsRegistry).
		WithVersionInfo(version.Version, version.Commit, version.BuildDate).
		Build()

	// Create collector with app reference for tracing
	mqttCollector := collectors.NewMQTTCollector(cfg, mqttRegistry, application)
	application.WithCollector(mqttCollector)

	if err := application.Run(); err != nil {
		slog.Error("Application failed", "error", err)
		os.Exit(1)
	}
}
