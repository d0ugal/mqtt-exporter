package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Duration represents a time duration that can be parsed from strings
type Duration struct {
	time.Duration
}

// UnmarshalYAML implements custom unmarshaling for duration strings
func (d *Duration) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value interface{}
	if err := unmarshal(&value); err != nil {
		return err
	}

	switch v := value.(type) {
	case string:
		duration, err := time.ParseDuration(v)
		if err != nil {
			return fmt.Errorf("invalid duration format '%s': %w", v, err)
		}

		d.Duration = duration
	case int:
		// Backward compatibility: treat as seconds
		d.Duration = time.Duration(v) * time.Second
	case int64:
		// Backward compatibility: treat as seconds
		d.Duration = time.Duration(v) * time.Second
	default:
		return fmt.Errorf("duration must be a string (e.g., '60s', '1h') or integer (seconds)")
	}

	return nil
}

// Seconds returns the duration in seconds
func (d *Duration) Seconds() int {
	return int(d.Duration.Seconds())
}

type Config struct {
	Server  ServerConfig  `yaml:"server"`
	Logging LoggingConfig `yaml:"logging"`
	Metrics MetricsConfig `yaml:"metrics"`
	MQTT    MQTTConfig    `yaml:"mqtt"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type LoggingConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"` // "json" or "text"
}

type MetricsConfig struct {
	Collection CollectionConfig `yaml:"collection"`
}

type CollectionConfig struct {
	DefaultInterval Duration `yaml:"default_interval"`
	// Track if the value was explicitly set
	DefaultIntervalSet bool `yaml:"-"`
}

// UnmarshalYAML implements custom unmarshaling to track if the value was set
func (c *CollectionConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// Create a temporary struct to unmarshal into
	type tempCollectionConfig struct {
		DefaultInterval Duration `yaml:"default_interval"`
	}

	var temp tempCollectionConfig
	if err := unmarshal(&temp); err != nil {
		return err
	}

	c.DefaultInterval = temp.DefaultInterval
	c.DefaultIntervalSet = true

	return nil
}

type MQTTConfig struct {
	Broker         string   `yaml:"broker"`
	ClientID       string   `yaml:"client_id"`
	Username       string   `yaml:"username"`
	Password       string   `yaml:"password"`
	Topics         []string `yaml:"topics"`
	QoS            int      `yaml:"qos"`
	CleanSession   bool     `yaml:"clean_session"`
	KeepAlive      int      `yaml:"keep_alive"`
	ConnectTimeout int      `yaml:"connect_timeout"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Set defaults
	if config.Server.Host == "" {
		config.Server.Host = "0.0.0.0"
	}

	if config.Server.Port == 0 {
		config.Server.Port = 8080
	}

	if config.Logging.Level == "" {
		config.Logging.Level = "info"
	}

	if config.Logging.Format == "" {
		config.Logging.Format = "json"
	}

	if !config.Metrics.Collection.DefaultIntervalSet {
		config.Metrics.Collection.DefaultInterval = Duration{time.Second * 30}
	}

	if config.MQTT.ClientID == "" {
		config.MQTT.ClientID = "mqtt-exporter"
	}

	if config.MQTT.QoS == 0 {
		config.MQTT.QoS = 1
	}

	if config.MQTT.KeepAlive == 0 {
		config.MQTT.KeepAlive = 60
	}

	if config.MQTT.ConnectTimeout == 0 {
		config.MQTT.ConnectTimeout = 30
	}

	if len(config.MQTT.Topics) == 0 {
		config.MQTT.Topics = []string{"#"}
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return &config, nil
}

// LoadConfig loads configuration from either a YAML file or environment variables
// If configFromEnv is true, it will load from environment variables only
func LoadConfig(path string, configFromEnv bool) (*Config, error) {
	if configFromEnv {
		return loadFromEnv()
	}

	return Load(path)
}

// loadFromEnv loads configuration from environment variables
func loadFromEnv() (*Config, error) {
	config := &Config{}

	// Server configuration
	if host := os.Getenv("MQTT_EXPORTER_SERVER_HOST"); host != "" {
		config.Server.Host = host
	} else {
		config.Server.Host = "0.0.0.0"
	}

	if portStr := os.Getenv("MQTT_EXPORTER_SERVER_PORT"); portStr != "" {
		if port, err := parseInt(portStr); err != nil {
			return nil, fmt.Errorf("invalid server port: %w", err)
		} else {
			config.Server.Port = port
		}
	} else {
		config.Server.Port = 8080
	}

	// Logging configuration
	if level := os.Getenv("MQTT_EXPORTER_LOG_LEVEL"); level != "" {
		config.Logging.Level = level
	} else {
		config.Logging.Level = "info"
	}

	if format := os.Getenv("MQTT_EXPORTER_LOG_FORMAT"); format != "" {
		config.Logging.Format = format
	} else {
		config.Logging.Format = "json"
	}

	// Metrics configuration
	if intervalStr := os.Getenv("MQTT_EXPORTER_METRICS_DEFAULT_INTERVAL"); intervalStr != "" {
		if interval, err := time.ParseDuration(intervalStr); err != nil {
			return nil, fmt.Errorf("invalid metrics default interval: %w", err)
		} else {
			config.Metrics.Collection.DefaultInterval = Duration{interval}
			config.Metrics.Collection.DefaultIntervalSet = true
		}
	} else {
		config.Metrics.Collection.DefaultInterval = Duration{time.Second * 30}
	}

	// MQTT configuration
	if broker := os.Getenv("MQTT_EXPORTER_MQTT_BROKER"); broker != "" {
		config.MQTT.Broker = broker
	} else {
		return nil, fmt.Errorf("MQTT broker is required (MQTT_EXPORTER_MQTT_BROKER)")
	}

	if clientID := os.Getenv("MQTT_EXPORTER_MQTT_CLIENT_ID"); clientID != "" {
		config.MQTT.ClientID = clientID
	} else {
		config.MQTT.ClientID = "mqtt-exporter"
	}

	if username := os.Getenv("MQTT_EXPORTER_MQTT_USERNAME"); username != "" {
		config.MQTT.Username = username
	}

	if password := os.Getenv("MQTT_EXPORTER_MQTT_PASSWORD"); password != "" {
		config.MQTT.Password = password
	}

	if topicsStr := os.Getenv("MQTT_EXPORTER_MQTT_TOPICS"); topicsStr != "" {
		config.MQTT.Topics = parseStringList(topicsStr)
	} else {
		config.MQTT.Topics = []string{"#"}
	}

	if qosStr := os.Getenv("MQTT_EXPORTER_MQTT_QOS"); qosStr != "" {
		if qos, err := parseInt(qosStr); err != nil {
			return nil, fmt.Errorf("invalid MQTT QoS: %w", err)
		} else {
			config.MQTT.QoS = qos
		}
	} else {
		config.MQTT.QoS = 1
	}

	if cleanSessionStr := os.Getenv("MQTT_EXPORTER_MQTT_CLEAN_SESSION"); cleanSessionStr != "" {
		if cleanSession, err := parseBool(cleanSessionStr); err != nil {
			return nil, fmt.Errorf("invalid MQTT clean session: %w", err)
		} else {
			config.MQTT.CleanSession = cleanSession
		}
	} else {
		config.MQTT.CleanSession = true
	}

	if keepAliveStr := os.Getenv("MQTT_EXPORTER_MQTT_KEEP_ALIVE"); keepAliveStr != "" {
		if keepAlive, err := parseInt(keepAliveStr); err != nil {
			return nil, fmt.Errorf("invalid MQTT keep alive: %w", err)
		} else {
			config.MQTT.KeepAlive = keepAlive
		}
	} else {
		config.MQTT.KeepAlive = 60
	}

	if connectTimeoutStr := os.Getenv("MQTT_EXPORTER_MQTT_CONNECT_TIMEOUT"); connectTimeoutStr != "" {
		if connectTimeout, err := parseInt(connectTimeoutStr); err != nil {
			return nil, fmt.Errorf("invalid MQTT connect timeout: %w", err)
		} else {
			config.MQTT.ConnectTimeout = connectTimeout
		}
	} else {
		config.MQTT.ConnectTimeout = 30
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return config, nil
}

// parseInt parses a string to int
func parseInt(s string) (int, error) {
	var i int

	_, err := fmt.Sscanf(s, "%d", &i)
	if err != nil {
		return 0, err
	}
	// Check if there are any remaining characters (like decimal points)
	if len(fmt.Sprintf("%d", i)) != len(s) {
		return 0, fmt.Errorf("invalid integer format: %s", s)
	}

	return i, nil
}

// parseBool parses a string to bool
func parseBool(s string) (bool, error) {
	switch s {
	case "true", "1", "yes", "on":
		return true, nil
	case "false", "0", "no", "off":
		return false, nil
	default:
		return false, fmt.Errorf("invalid boolean value: %s", s)
	}
}

// parseStringList parses a comma-separated string into a slice of strings
func parseStringList(s string) []string {
	if s == "" {
		return nil
	}

	// Split by comma and trim whitespace
	parts := strings.Split(s, ",")

	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}

// Validate performs comprehensive validation of the configuration
func (c *Config) Validate() error {
	// Validate server configuration
	if err := c.validateServerConfig(); err != nil {
		return fmt.Errorf("server config: %w", err)
	}

	// Validate logging configuration
	if err := c.validateLoggingConfig(); err != nil {
		return fmt.Errorf("logging config: %w", err)
	}

	// Validate metrics configuration
	if err := c.validateMetricsConfig(); err != nil {
		return fmt.Errorf("metrics config: %w", err)
	}

	// Validate MQTT configuration
	if err := c.validateMQTTConfig(); err != nil {
		return fmt.Errorf("mqtt config: %w", err)
	}

	return nil
}

func (c *Config) validateServerConfig() error {
	if c.Server.Port < 1 || c.Server.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535, got %d", c.Server.Port)
	}

	return nil
}

func (c *Config) validateLoggingConfig() error {
	validLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	if !validLevels[c.Logging.Level] {
		return fmt.Errorf("invalid logging level: %s", c.Logging.Level)
	}

	validFormats := map[string]bool{
		"json": true,
		"text": true,
	}
	if !validFormats[c.Logging.Format] {
		return fmt.Errorf("invalid logging format: %s", c.Logging.Format)
	}

	return nil
}

func (c *Config) validateMetricsConfig() error {
	if c.Metrics.Collection.DefaultInterval.Seconds() < 1 {
		return fmt.Errorf("default interval must be at least 1 second, got %d", c.Metrics.Collection.DefaultInterval.Seconds())
	}

	if c.Metrics.Collection.DefaultInterval.Seconds() > 86400 {
		return fmt.Errorf("default interval must be at most 86400 seconds (24 hours), got %d", c.Metrics.Collection.DefaultInterval.Seconds())
	}

	return nil
}

func (c *Config) validateMQTTConfig() error {
	if c.MQTT.Broker == "" {
		return fmt.Errorf("mqtt broker is required")
	}

	if c.MQTT.ClientID == "" {
		return fmt.Errorf("mqtt client_id is required")
	}

	if c.MQTT.QoS < 0 || c.MQTT.QoS > 2 {
		return fmt.Errorf("mqtt qos must be between 0 and 2, got %d", c.MQTT.QoS)
	}

	if c.MQTT.KeepAlive < 1 {
		return fmt.Errorf("mqtt keep_alive must be at least 1, got %d", c.MQTT.KeepAlive)
	}

	if c.MQTT.ConnectTimeout < 1 {
		return fmt.Errorf("mqtt connect_timeout must be at least 1, got %d", c.MQTT.ConnectTimeout)
	}

	if len(c.MQTT.Topics) == 0 {
		return fmt.Errorf("at least one mqtt topic must be configured")
	}

	return nil
}

// GetDefaultInterval returns the default collection interval
func (c *Config) GetDefaultInterval() int {
	return c.Metrics.Collection.DefaultInterval.Seconds()
}
