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
