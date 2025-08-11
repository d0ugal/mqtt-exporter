package config

import (
	"os"
	"testing"
	"time"
)

func TestDuration_UnmarshalYAML(t *testing.T) {
	// Test with a simple string value
	var d Duration

	err := d.UnmarshalYAML(func(v interface{}) error {
		// Simulate YAML unmarshaling by setting the value
		*(v.(*interface{})) = "30s"
		return nil
	})
	if err != nil {
		t.Errorf("Duration.UnmarshalYAML() unexpected error: %v", err)
		return
	}

	expected := 30 * time.Second
	if d.Duration != expected {
		t.Errorf("Duration.UnmarshalYAML() = %v, want %v", d.Duration, expected)
	}
}

func TestDuration_UnmarshalYAML_Integer(t *testing.T) {
	// Test with integer value (backward compatibility)
	var d Duration

	err := d.UnmarshalYAML(func(v interface{}) error {
		// Simulate YAML unmarshaling by setting the value
		*(v.(*interface{})) = 60
		return nil
	})
	if err != nil {
		t.Errorf("Duration.UnmarshalYAML() unexpected error: %v", err)
		return
	}

	expected := 60 * time.Second
	if d.Duration != expected {
		t.Errorf("Duration.UnmarshalYAML() = %v, want %v", d.Duration, expected)
	}
}

func TestDuration_UnmarshalYAML_Invalid(t *testing.T) {
	// Test with invalid value
	var d Duration

	err := d.UnmarshalYAML(func(v interface{}) error {
		// Simulate YAML unmarshaling by setting the value
		*(v.(*interface{})) = "invalid"
		return nil
	})
	if err == nil {
		t.Error("Duration.UnmarshalYAML() expected error, got nil")
	}
}

func TestDuration_Seconds(t *testing.T) {
	d := Duration{30 * time.Second}
	if got := d.Seconds(); got != 30 {
		t.Errorf("Duration.Seconds() = %v, want %v", got, 30)
	}
}

func TestLoad(t *testing.T) {
	// Create a temporary config file
	configContent := `
server:
  host: "127.0.0.1"
  port: 9090

logging:
  level: "debug"
  format: "text"

metrics:
  collection:
    default_interval: "45s"

mqtt:
  broker: "localhost:1883"
  client_id: "test-client"
  username: "testuser"
  password: "testpass"
  topics:
    - "test/topic"
    - "another/topic"
  qos: 2
  clean_session: false
  keep_alive: 120
  connect_timeout: 60
`

	tmpfile, err := os.CreateTemp("", "config_test_*.yaml")
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := os.Remove(tmpfile.Name()); err != nil {
			t.Logf("Failed to remove temp file: %v", err)
		}
	}()

	if _, err := tmpfile.Write([]byte(configContent)); err != nil {
		t.Fatal(err)
	}

	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Test loading the config
	cfg, err := Load(tmpfile.Name())
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}

	// Verify the loaded configuration
	if cfg.Server.Host != "127.0.0.1" {
		t.Errorf("Server.Host = %v, want %v", cfg.Server.Host, "127.0.0.1")
	}

	if cfg.Server.Port != 9090 {
		t.Errorf("Server.Port = %v, want %v", cfg.Server.Port, 9090)
	}

	if cfg.Logging.Level != "debug" {
		t.Errorf("Logging.Level = %v, want %v", cfg.Logging.Level, "debug")
	}

	if cfg.Logging.Format != "text" {
		t.Errorf("Logging.Format = %v, want %v", cfg.Logging.Format, "text")
	}

	if cfg.Metrics.Collection.DefaultInterval.Seconds() != 45 {
		t.Errorf("Metrics.Collection.DefaultInterval = %v, want %v", cfg.Metrics.Collection.DefaultInterval.Seconds(), 45)
	}

	if cfg.MQTT.Broker != "localhost:1883" {
		t.Errorf("MQTT.Broker = %v, want %v", cfg.MQTT.Broker, "localhost:1883")
	}

	if cfg.MQTT.ClientID != "test-client" {
		t.Errorf("MQTT.ClientID = %v, want %v", cfg.MQTT.ClientID, "test-client")
	}

	if cfg.MQTT.Username != "testuser" {
		t.Errorf("MQTT.Username = %v, want %v", cfg.MQTT.Username, "testuser")
	}

	if cfg.MQTT.Password != "testpass" {
		t.Errorf("MQTT.Password = %v, want %v", cfg.MQTT.Password, "testpass")
	}

	if len(cfg.MQTT.Topics) != 2 {
		t.Errorf("MQTT.Topics length = %v, want %v", len(cfg.MQTT.Topics), 2)
	}

	if cfg.MQTT.QoS != 2 {
		t.Errorf("MQTT.QoS = %v, want %v", cfg.MQTT.QoS, 2)
	}

	if cfg.MQTT.CleanSession != false {
		t.Errorf("MQTT.CleanSession = %v, want %v", cfg.MQTT.CleanSession, false)
	}

	if cfg.MQTT.KeepAlive != 120 {
		t.Errorf("MQTT.KeepAlive = %v, want %v", cfg.MQTT.KeepAlive, 120)
	}

	if cfg.MQTT.ConnectTimeout != 60 {
		t.Errorf("MQTT.ConnectTimeout = %v, want %v", cfg.MQTT.ConnectTimeout, 60)
	}
}

func TestLoad_Defaults(t *testing.T) {
	// Create a minimal config file
	configContent := `
mqtt:
  broker: "localhost:1883"
`

	tmpfile, err := os.CreateTemp("", "config_test_*.yaml")
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := os.Remove(tmpfile.Name()); err != nil {
			t.Logf("Failed to remove temp file: %v", err)
		}
	}()

	if _, err := tmpfile.Write([]byte(configContent)); err != nil {
		t.Fatal(err)
	}

	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Test loading the config
	cfg, err := Load(tmpfile.Name())
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}

	// Verify defaults are set
	if cfg.Server.Host != "0.0.0.0" {
		t.Errorf("Server.Host default = %v, want %v", cfg.Server.Host, "0.0.0.0")
	}

	if cfg.Server.Port != 8080 {
		t.Errorf("Server.Port default = %v, want %v", cfg.Server.Port, 8080)
	}

	if cfg.Logging.Level != "info" {
		t.Errorf("Logging.Level default = %v, want %v", cfg.Logging.Level, "info")
	}

	if cfg.Logging.Format != "json" {
		t.Errorf("Logging.Format default = %v, want %v", cfg.Logging.Format, "json")
	}

	if cfg.Metrics.Collection.DefaultInterval.Seconds() != 30 {
		t.Errorf("Metrics.Collection.DefaultInterval default = %v, want %v", cfg.Metrics.Collection.DefaultInterval.Seconds(), 30)
	}

	if cfg.MQTT.ClientID != "mqtt-exporter" {
		t.Errorf("MQTT.ClientID default = %v, want %v", cfg.MQTT.ClientID, "mqtt-exporter")
	}

	if cfg.MQTT.QoS != 1 {
		t.Errorf("MQTT.QoS default = %v, want %v", cfg.MQTT.QoS, 1)
	}

	if cfg.MQTT.KeepAlive != 60 {
		t.Errorf("MQTT.KeepAlive default = %v, want %v", cfg.MQTT.KeepAlive, 60)
	}

	if cfg.MQTT.ConnectTimeout != 30 {
		t.Errorf("MQTT.ConnectTimeout default = %v, want %v", cfg.MQTT.ConnectTimeout, 30)
	}

	if len(cfg.MQTT.Topics) != 1 || cfg.MQTT.Topics[0] != "#" {
		t.Errorf("MQTT.Topics default = %v, want %v", cfg.MQTT.Topics, []string{"#"})
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := Load("nonexistent.yaml")
	if err == nil {
		t.Error("Load() expected error for nonexistent file, got nil")
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	// Create a file with invalid YAML
	configContent := `
server:
  host: "127.0.0.1"
  port: invalid_port
`

	tmpfile, err := os.CreateTemp("", "config_test_*.yaml")
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := os.Remove(tmpfile.Name()); err != nil {
			t.Logf("Failed to remove temp file: %v", err)
		}
	}()

	if _, err := tmpfile.Write([]byte(configContent)); err != nil {
		t.Fatal(err)
	}

	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	_, err = Load(tmpfile.Name())
	if err == nil {
		t.Error("Load() expected error for invalid YAML, got nil")
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: Config{
				Server:  ServerConfig{Host: "0.0.0.0", Port: 8080},
				Logging: LoggingConfig{Level: "info", Format: "json"},
				Metrics: MetricsConfig{
					Collection: CollectionConfig{
						DefaultInterval: Duration{30 * time.Second},
					},
				},
				MQTT: MQTTConfig{
					Broker:         "localhost:1883",
					ClientID:       "test",
					QoS:            1,
					KeepAlive:      60,
					ConnectTimeout: 30,
					Topics:         []string{"test/topic"},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid server port",
			config: Config{
				Server:  ServerConfig{Host: "0.0.0.0", Port: 0},
				Logging: LoggingConfig{Level: "info", Format: "json"},
				Metrics: MetricsConfig{
					Collection: CollectionConfig{
						DefaultInterval: Duration{30 * time.Second},
					},
				},
				MQTT: MQTTConfig{
					Broker:         "localhost:1883",
					ClientID:       "test",
					QoS:            1,
					KeepAlive:      60,
					ConnectTimeout: 30,
					Topics:         []string{"test/topic"},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid logging level",
			config: Config{
				Server:  ServerConfig{Host: "0.0.0.0", Port: 8080},
				Logging: LoggingConfig{Level: "invalid", Format: "json"},
				Metrics: MetricsConfig{
					Collection: CollectionConfig{
						DefaultInterval: Duration{30 * time.Second},
					},
				},
				MQTT: MQTTConfig{
					Broker:         "localhost:1883",
					ClientID:       "test",
					QoS:            1,
					KeepAlive:      60,
					ConnectTimeout: 30,
					Topics:         []string{"test/topic"},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid logging format",
			config: Config{
				Server:  ServerConfig{Host: "0.0.0.0", Port: 8080},
				Logging: LoggingConfig{Level: "info", Format: "invalid"},
				Metrics: MetricsConfig{
					Collection: CollectionConfig{
						DefaultInterval: Duration{30 * time.Second},
					},
				},
				MQTT: MQTTConfig{
					Broker:         "localhost:1883",
					ClientID:       "test",
					QoS:            1,
					KeepAlive:      60,
					ConnectTimeout: 30,
					Topics:         []string{"test/topic"},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid metrics interval",
			config: Config{
				Server:  ServerConfig{Host: "0.0.0.0", Port: 8080},
				Logging: LoggingConfig{Level: "info", Format: "json"},
				Metrics: MetricsConfig{
					Collection: CollectionConfig{
						DefaultInterval: Duration{0},
					},
				},
				MQTT: MQTTConfig{
					Broker:         "localhost:1883",
					ClientID:       "test",
					QoS:            1,
					KeepAlive:      60,
					ConnectTimeout: 30,
					Topics:         []string{"test/topic"},
				},
			},
			wantErr: true,
		},
		{
			name: "missing MQTT broker",
			config: Config{
				Server:  ServerConfig{Host: "0.0.0.0", Port: 8080},
				Logging: LoggingConfig{Level: "info", Format: "json"},
				Metrics: MetricsConfig{
					Collection: CollectionConfig{
						DefaultInterval: Duration{30 * time.Second},
					},
				},
				MQTT: MQTTConfig{
					Broker:         "",
					ClientID:       "test",
					QoS:            1,
					KeepAlive:      60,
					ConnectTimeout: 30,
					Topics:         []string{"test/topic"},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid MQTT QoS",
			config: Config{
				Server:  ServerConfig{Host: "0.0.0.0", Port: 8080},
				Logging: LoggingConfig{Level: "info", Format: "json"},
				Metrics: MetricsConfig{
					Collection: CollectionConfig{
						DefaultInterval: Duration{30 * time.Second},
					},
				},
				MQTT: MQTTConfig{
					Broker:         "localhost:1883",
					ClientID:       "test",
					QoS:            3,
					KeepAlive:      60,
					ConnectTimeout: 30,
					Topics:         []string{"test/topic"},
				},
			},
			wantErr: true,
		},
		{
			name: "empty MQTT topics",
			config: Config{
				Server:  ServerConfig{Host: "0.0.0.0", Port: 8080},
				Logging: LoggingConfig{Level: "info", Format: "json"},
				Metrics: MetricsConfig{
					Collection: CollectionConfig{
						DefaultInterval: Duration{30 * time.Second},
					},
				},
				MQTT: MQTTConfig{
					Broker:         "localhost:1883",
					ClientID:       "test",
					QoS:            1,
					KeepAlive:      60,
					ConnectTimeout: 30,
					Topics:         []string{},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfig_GetDefaultInterval(t *testing.T) {
	cfg := Config{
		Metrics: MetricsConfig{
			Collection: CollectionConfig{
				DefaultInterval: Duration{45 * time.Second},
			},
		},
	}

	if got := cfg.GetDefaultInterval(); got != 45 {
		t.Errorf("Config.GetDefaultInterval() = %v, want %v", got, 45)
	}
}

// setTestEnvVars sets multiple environment variables for testing
func setTestEnvVars(t *testing.T, envVars map[string]string) func() {
	for key, value := range envVars {
		if err := os.Setenv(key, value); err != nil {
			t.Fatalf("Failed to set environment variable %s: %v", key, err)
		}
	}

	// Return cleanup function
	return func() {
		for key := range envVars {
			if err := os.Unsetenv(key); err != nil {
				t.Logf("Failed to unset environment variable %s: %v", key, err)
			}
		}
	}
}

func TestLoadConfig_EnvironmentVariables(t *testing.T) {
	// Set environment variables for testing
	envVars := map[string]string{
		"MQTT_EXPORTER_SERVER_HOST":              "127.0.0.1",
		"MQTT_EXPORTER_SERVER_PORT":              "9090",
		"MQTT_EXPORTER_LOG_LEVEL":                "debug",
		"MQTT_EXPORTER_LOG_FORMAT":               "text",
		"MQTT_EXPORTER_METRICS_DEFAULT_INTERVAL": "45s",
		"MQTT_EXPORTER_MQTT_BROKER":              "localhost:1883",
		"MQTT_EXPORTER_MQTT_CLIENT_ID":           "test-client",
		"MQTT_EXPORTER_MQTT_USERNAME":            "testuser",
		"MQTT_EXPORTER_MQTT_PASSWORD":            "testpass",
		"MQTT_EXPORTER_MQTT_TOPICS":              "test/topic,another/topic",
		"MQTT_EXPORTER_MQTT_QOS":                 "2",
		"MQTT_EXPORTER_MQTT_CLEAN_SESSION":       "false",
		"MQTT_EXPORTER_MQTT_KEEP_ALIVE":          "120",
		"MQTT_EXPORTER_MQTT_CONNECT_TIMEOUT":     "60",
	}

	cleanup := setTestEnvVars(t, envVars)
	defer cleanup()

	// Load configuration from environment variables
	cfg, err := LoadConfig("", true)
	if err != nil {
		t.Fatalf("LoadConfig() unexpected error: %v", err)
	}

	// Verify server configuration
	if cfg.Server.Host != "127.0.0.1" {
		t.Errorf("Server.Host = %v, want %v", cfg.Server.Host, "127.0.0.1")
	}

	if cfg.Server.Port != 9090 {
		t.Errorf("Server.Port = %v, want %v", cfg.Server.Port, 9090)
	}

	// Verify logging configuration
	if cfg.Logging.Level != "debug" {
		t.Errorf("Logging.Level = %v, want %v", cfg.Logging.Level, "debug")
	}

	if cfg.Logging.Format != "text" {
		t.Errorf("Logging.Format = %v, want %v", cfg.Logging.Format, "text")
	}

	// Verify metrics configuration
	if cfg.Metrics.Collection.DefaultInterval.Seconds() != 45 {
		t.Errorf("Metrics.Collection.DefaultInterval = %v, want %v", cfg.Metrics.Collection.DefaultInterval.Seconds(), 45)
	}

	// Verify MQTT configuration
	if cfg.MQTT.Broker != "localhost:1883" {
		t.Errorf("MQTT.Broker = %v, want %v", cfg.MQTT.Broker, "localhost:1883")
	}

	if cfg.MQTT.ClientID != "test-client" {
		t.Errorf("MQTT.ClientID = %v, want %v", cfg.MQTT.ClientID, "test-client")
	}

	if cfg.MQTT.Username != "testuser" {
		t.Errorf("MQTT.Username = %v, want %v", cfg.MQTT.Username, "testuser")
	}

	if cfg.MQTT.Password != "testpass" {
		t.Errorf("MQTT.Password = %v, want %v", cfg.MQTT.Password, "testpass")
	}

	if len(cfg.MQTT.Topics) != 2 {
		t.Errorf("MQTT.Topics length = %v, want %v", len(cfg.MQTT.Topics), 2)
	}

	if cfg.MQTT.Topics[0] != "test/topic" || cfg.MQTT.Topics[1] != "another/topic" {
		t.Errorf("MQTT.Topics = %v, want %v", cfg.MQTT.Topics, []string{"test/topic", "another/topic"})
	}

	if cfg.MQTT.QoS != 2 {
		t.Errorf("MQTT.QoS = %v, want %v", cfg.MQTT.QoS, 2)
	}

	if cfg.MQTT.CleanSession {
		t.Errorf("MQTT.CleanSession = %v, want %v", cfg.MQTT.CleanSession, false)
	}

	if cfg.MQTT.KeepAlive != 120 {
		t.Errorf("MQTT.KeepAlive = %v, want %v", cfg.MQTT.KeepAlive, 120)
	}

	if cfg.MQTT.ConnectTimeout != 60 {
		t.Errorf("MQTT.ConnectTimeout = %v, want %v", cfg.MQTT.ConnectTimeout, 60)
	}
}

func TestLoadConfig_EnvironmentVariables_Defaults(t *testing.T) {
	// Set only required environment variable
	if err := os.Setenv("MQTT_EXPORTER_MQTT_BROKER", "localhost:1883"); err != nil {
		t.Fatalf("Failed to set environment variable: %v", err)
	}

	// Clean up environment variable after test
	defer func() {
		if err := os.Unsetenv("MQTT_EXPORTER_MQTT_BROKER"); err != nil {
			t.Logf("Failed to unset environment variable: %v", err)
		}
	}()

	// Load configuration from environment variables
	cfg, err := LoadConfig("", true)
	if err != nil {
		t.Fatalf("LoadConfig() unexpected error: %v", err)
	}

	// Verify defaults are set correctly
	if cfg.Server.Host != "0.0.0.0" {
		t.Errorf("Server.Host = %v, want %v", cfg.Server.Host, "0.0.0.0")
	}

	if cfg.Server.Port != 8080 {
		t.Errorf("Server.Port = %v, want %v", cfg.Server.Port, 8080)
	}

	if cfg.Logging.Level != "info" {
		t.Errorf("Logging.Level = %v, want %v", cfg.Logging.Level, "info")
	}

	if cfg.Logging.Format != "json" {
		t.Errorf("Logging.Format = %v, want %v", cfg.Logging.Format, "json")
	}

	if cfg.MQTT.ClientID != "mqtt-exporter" {
		t.Errorf("MQTT.ClientID = %v, want %v", cfg.MQTT.ClientID, "mqtt-exporter")
	}

	if cfg.MQTT.QoS != 1 {
		t.Errorf("MQTT.QoS = %v, want %v", cfg.MQTT.QoS, 1)
	}

	if !cfg.MQTT.CleanSession {
		t.Errorf("MQTT.CleanSession = %v, want %v", cfg.MQTT.CleanSession, true)
	}

	if cfg.MQTT.KeepAlive != 60 {
		t.Errorf("MQTT.KeepAlive = %v, want %v", cfg.MQTT.KeepAlive, 60)
	}

	if cfg.MQTT.ConnectTimeout != 30 {
		t.Errorf("MQTT.ConnectTimeout = %v, want %v", cfg.MQTT.ConnectTimeout, 30)
	}

	if len(cfg.MQTT.Topics) != 1 || cfg.MQTT.Topics[0] != "#" {
		t.Errorf("MQTT.Topics = %v, want %v", cfg.MQTT.Topics, []string{"#"})
	}
}

func TestLoadConfig_EnvironmentVariables_MissingBroker(t *testing.T) {
	// Don't set MQTT broker environment variable
	cfg, err := LoadConfig("", true)
	if err == nil {
		t.Error("LoadConfig() expected error for missing MQTT broker, got nil")
	}

	if cfg != nil {
		t.Error("LoadConfig() expected nil config for missing MQTT broker")
	}
}

func TestLoadConfig_EnvironmentVariables_InvalidValues(t *testing.T) {
	// Set invalid environment variables
	if err := os.Setenv("MQTT_EXPORTER_MQTT_BROKER", "localhost:1883"); err != nil {
		t.Fatalf("Failed to set environment variable: %v", err)
	}

	if err := os.Setenv("MQTT_EXPORTER_SERVER_PORT", "invalid"); err != nil {
		t.Fatalf("Failed to set environment variable: %v", err)
	}

	if err := os.Setenv("MQTT_EXPORTER_MQTT_QOS", "invalid"); err != nil {
		t.Fatalf("Failed to set environment variable: %v", err)
	}

	if err := os.Setenv("MQTT_EXPORTER_MQTT_CLEAN_SESSION", "invalid"); err != nil {
		t.Fatalf("Failed to set environment variable: %v", err)
	}

	if err := os.Setenv("MQTT_EXPORTER_METRICS_DEFAULT_INTERVAL", "invalid"); err != nil {
		t.Fatalf("Failed to set environment variable: %v", err)
	}

	// Clean up environment variables after test
	defer func() {
		if err := os.Unsetenv("MQTT_EXPORTER_MQTT_BROKER"); err != nil {
			t.Logf("Failed to unset environment variable: %v", err)
		}

		if err := os.Unsetenv("MQTT_EXPORTER_SERVER_PORT"); err != nil {
			t.Logf("Failed to unset environment variable: %v", err)
		}

		if err := os.Unsetenv("MQTT_EXPORTER_MQTT_QOS"); err != nil {
			t.Logf("Failed to unset environment variable: %v", err)
		}

		if err := os.Unsetenv("MQTT_EXPORTER_MQTT_CLEAN_SESSION"); err != nil {
			t.Logf("Failed to unset environment variable: %v", err)
		}

		if err := os.Unsetenv("MQTT_EXPORTER_METRICS_DEFAULT_INTERVAL"); err != nil {
			t.Logf("Failed to unset environment variable: %v", err)
		}
	}()

	// Test invalid server port
	_, err := LoadConfig("", true)
	if err == nil {
		t.Error("LoadConfig() expected error for invalid server port, got nil")
	}

	// Test invalid QoS
	if err := os.Unsetenv("MQTT_EXPORTER_SERVER_PORT"); err != nil {
		t.Logf("Failed to unset environment variable: %v", err)
	}

	_, err = LoadConfig("", true)
	if err == nil {
		t.Error("LoadConfig() expected error for invalid QoS, got nil")
	}

	// Test invalid clean session
	if err := os.Unsetenv("MQTT_EXPORTER_MQTT_QOS"); err != nil {
		t.Logf("Failed to unset environment variable: %v", err)
	}

	_, err = LoadConfig("", true)
	if err == nil {
		t.Error("LoadConfig() expected error for invalid clean session, got nil")
	}

	// Test invalid metrics interval
	if err := os.Unsetenv("MQTT_EXPORTER_MQTT_CLEAN_SESSION"); err != nil {
		t.Logf("Failed to unset environment variable: %v", err)
	}

	_, err = LoadConfig("", true)
	if err == nil {
		t.Error("LoadConfig() expected error for invalid metrics interval, got nil")
	}
}

func TestParseStringList(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: nil,
		},
		{
			name:     "single item",
			input:    "test",
			expected: []string{"test"},
		},
		{
			name:     "multiple items",
			input:    "test1,test2,test3",
			expected: []string{"test1", "test2", "test3"},
		},
		{
			name:     "items with spaces",
			input:    " test1 , test2 , test3 ",
			expected: []string{"test1", "test2", "test3"},
		},
		{
			name:     "empty items",
			input:    "test1,,test2, ,test3",
			expected: []string{"test1", "test2", "test3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseStringList(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("parseStringList() length = %v, want %v", len(result), len(tt.expected))
				return
			}

			for i, item := range result {
				if item != tt.expected[i] {
					t.Errorf("parseStringList()[%d] = %v, want %v", i, item, tt.expected[i])
				}
			}
		})
	}
}

func TestParseBool(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
		wantErr  bool
	}{
		{"true", "true", true, false},
		{"1", "1", true, false},
		{"yes", "yes", true, false},
		{"on", "on", true, false},
		{"false", "false", false, false},
		{"0", "0", false, false},
		{"no", "no", false, false},
		{"off", "off", false, false},
		{"invalid", "invalid", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseBool(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseBool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if result != tt.expected {
				t.Errorf("parseBool() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestParseInt(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
		wantErr  bool
	}{
		{"positive", "123", 123, false},
		{"zero", "0", 0, false},
		{"negative", "-456", -456, false},
		{"invalid", "invalid", 0, true},
		{"float", "123.45", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseInt(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if result != tt.expected {
				t.Errorf("parseInt() = %v, want %v", result, tt.expected)
			}
		})
	}
}
