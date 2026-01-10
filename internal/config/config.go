package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	promexporter_config "github.com/d0ugal/promexporter/config"
	"gopkg.in/yaml.v3"
)

// Duration uses promexporter Duration type
type Duration = promexporter_config.Duration

type Config struct {
	promexporter_config.BaseConfig

	MQTT MQTTConfig `yaml:"mqtt"`
}

type MQTTConfig struct {
	Broker           string   `yaml:"broker"`
	ClientID         string   `yaml:"client_id"`
	Username         string   `yaml:"username"`
	Password         string   `yaml:"password"`
	Topics           []string `yaml:"topics"`
	QoS              int      `yaml:"qos"`
	CleanSession     bool     `yaml:"clean_session"`
	KeepAlive        Duration `yaml:"keep_alive"`
	ConnectTimeout   Duration `yaml:"connect_timeout"`
	DisableSysTopics bool     `yaml:"disable_sys_topics"`
}

// LoadConfig loads configuration from either a YAML file or environment variables
func LoadConfig(path string, configFromEnv bool) (*Config, error) {
	if configFromEnv {
		return loadFromEnv()
	}

	return Load(path)
}

// Load loads configuration from a YAML file
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
	setDefaults(&config)

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return &config, nil
}

// loadFromEnv loads configuration from environment variables
func loadFromEnv() (*Config, error) {
	config := &Config{}

	// Load base configuration from environment
	baseConfig := &promexporter_config.BaseConfig{}

	// Server configuration
	if host := os.Getenv("MQTT_EXPORTER_SERVER_HOST"); host != "" {
		baseConfig.Server.Host = host
	} else {
		baseConfig.Server.Host = "0.0.0.0"
	}

	if portStr := os.Getenv("MQTT_EXPORTER_SERVER_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err != nil {
			return nil, fmt.Errorf("invalid server port: %w", err)
		} else {
			baseConfig.Server.Port = port
		}
	} else {
		baseConfig.Server.Port = 8080
	}

	// Logging configuration
	if level := os.Getenv("MQTT_EXPORTER_LOG_LEVEL"); level != "" {
		baseConfig.Logging.Level = level
	} else {
		baseConfig.Logging.Level = "info"
	}

	if format := os.Getenv("MQTT_EXPORTER_LOG_FORMAT"); format != "" {
		baseConfig.Logging.Format = format
	} else {
		baseConfig.Logging.Format = "json"
	}

	// Metrics configuration
	if intervalStr := os.Getenv("MQTT_EXPORTER_METRICS_DEFAULT_INTERVAL"); intervalStr != "" {
		if interval, err := time.ParseDuration(intervalStr); err != nil {
			return nil, fmt.Errorf("invalid metrics default interval: %w", err)
		} else {
			baseConfig.Metrics.Collection.DefaultInterval = promexporter_config.Duration{Duration: interval}
			baseConfig.Metrics.Collection.DefaultIntervalSet = true
		}
	} else {
		baseConfig.Metrics.Collection.DefaultInterval = promexporter_config.Duration{Duration: time.Second * 30}
	}

	config.BaseConfig = *baseConfig

	// Apply generic environment variables (TRACING_ENABLED, PROFILING_ENABLED, etc.)
	// These are handled by promexporter and are shared across all exporters
	if err := promexporter_config.ApplyGenericEnvVars(&config.BaseConfig); err != nil {
		return nil, fmt.Errorf("failed to apply generic environment variables: %w", err)
	}

	// MQTT configuration
	if broker := os.Getenv("MQTT_EXPORTER_MQTT_BROKER"); broker != "" {
		config.MQTT.Broker = broker
	} else {
		config.MQTT.Broker = "tcp://localhost:1883"
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
		config.MQTT.Topics = strings.Split(topicsStr, ",")
	} else {
		config.MQTT.Topics = []string{"#"}
	}

	if qosStr := os.Getenv("MQTT_EXPORTER_MQTT_QOS"); qosStr != "" {
		if qos, err := strconv.Atoi(qosStr); err != nil {
			return nil, fmt.Errorf("invalid MQTT QoS: %w", err)
		} else {
			config.MQTT.QoS = qos
		}
	} else {
		config.MQTT.QoS = 0
	}

	if cleanSessionStr := os.Getenv("MQTT_EXPORTER_MQTT_CLEAN_SESSION"); cleanSessionStr != "" {
		if cleanSession, err := strconv.ParseBool(cleanSessionStr); err != nil {
			return nil, fmt.Errorf("invalid MQTT clean session: %w", err)
		} else {
			config.MQTT.CleanSession = cleanSession
		}
	} else {
		config.MQTT.CleanSession = true
	}

	if keepAliveStr := os.Getenv("MQTT_EXPORTER_MQTT_KEEP_ALIVE"); keepAliveStr != "" {
		if keepAlive, err := time.ParseDuration(keepAliveStr); err != nil {
			return nil, fmt.Errorf("invalid MQTT keep alive: %w", err)
		} else {
			config.MQTT.KeepAlive = Duration{Duration: keepAlive}
		}
	} else {
		config.MQTT.KeepAlive = Duration{Duration: time.Second * 60}
	}

	if connectTimeoutStr := os.Getenv("MQTT_EXPORTER_MQTT_CONNECT_TIMEOUT"); connectTimeoutStr != "" {
		if connectTimeout, err := time.ParseDuration(connectTimeoutStr); err != nil {
			return nil, fmt.Errorf("invalid MQTT connect timeout: %w", err)
		} else {
			config.MQTT.ConnectTimeout = Duration{Duration: connectTimeout}
		}
	} else {
		config.MQTT.ConnectTimeout = Duration{Duration: time.Second * 30}
	}

	if disableSysTopicsStr := os.Getenv("MQTT_EXPORTER_MQTT_DISABLE_SYS_TOPICS"); disableSysTopicsStr != "" {
		if disableSysTopics, err := strconv.ParseBool(disableSysTopicsStr); err != nil {
			return nil, fmt.Errorf("invalid MQTT disable sys topics: %w", err)
		} else {
			config.MQTT.DisableSysTopics = disableSysTopics
		}
	}

	// Set defaults for any missing values
	setDefaults(config)

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return config, nil
}

// setDefaults sets default values for configuration
func setDefaults(config *Config) {
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
		config.Metrics.Collection.DefaultInterval = promexporter_config.Duration{Duration: time.Second * 30}
	}

	if config.MQTT.Broker == "" {
		config.MQTT.Broker = "tcp://localhost:1883"
	}

	if config.MQTT.ClientID == "" {
		config.MQTT.ClientID = "mqtt-exporter"
	}

	if len(config.MQTT.Topics) == 0 {
		config.MQTT.Topics = []string{"#"}
	}

	if config.MQTT.KeepAlive.Duration == 0 {
		config.MQTT.KeepAlive = Duration{Duration: time.Second * 60}
	}

	if config.MQTT.ConnectTimeout.Duration == 0 {
		config.MQTT.ConnectTimeout = Duration{Duration: time.Second * 30}
	}
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
		return fmt.Errorf("mqtt client id is required")
	}

	if c.MQTT.QoS < 0 || c.MQTT.QoS > 2 {
		return fmt.Errorf("mqtt qos must be between 0 and 2, got %d", c.MQTT.QoS)
	}

	if c.MQTT.KeepAlive.Seconds() < 0 {
		return fmt.Errorf("mqtt keep alive must be non-negative, got %d", c.MQTT.KeepAlive.Seconds())
	}

	if c.MQTT.ConnectTimeout.Seconds() < 1 {
		return fmt.Errorf("mqtt connect timeout must be at least 1 second, got %d", c.MQTT.ConnectTimeout.Seconds())
	}

	return nil
}

// GetDefaultInterval returns the default collection interval
func (c *Config) GetDefaultInterval() int {
	return c.Metrics.Collection.DefaultInterval.Seconds()
}

// ParseStringList parses a comma-separated string into a slice of strings
func ParseStringList(input string) []string {
	if input == "" {
		return []string{}
	}

	parts := strings.Split(input, ",")
	result := make([]string, 0, len(parts))

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}

// ParseBool parses a string to boolean
func ParseBool(input string) (bool, error) {
	switch strings.ToLower(strings.TrimSpace(input)) {
	case "true", "1", "yes", "on":
		return true, nil
	case "false", "0", "no", "off":
		return false, nil
	default:
		return false, fmt.Errorf("invalid boolean value: %s", input)
	}
}

// ParseInt parses a string to int
func ParseInt(input string) (int, error) {
	return strconv.Atoi(strings.TrimSpace(input))
}
