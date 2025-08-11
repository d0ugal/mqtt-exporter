# MQTT Exporter

A Prometheus exporter for MQTT message monitoring. This exporter connects to an MQTT broker, subscribes to configured topics, and exposes metrics about message counts and connection status.

## Features

- Monitor MQTT message counts per topic
- Track message bytes received
- Connection status monitoring
- Error tracking
- Last message timestamps
- Configurable topic subscriptions
- Prometheus metrics endpoint
- Health check endpoint

## Quick Start

### Using Docker

```bash
# Pull the latest image
docker pull ghcr.io/d0ugal/mqtt-exporter:latest

# Run with default configuration
docker run -p 8080:8080 ghcr.io/d0ugal/mqtt-exporter:latest

# Run with custom config (mount config file)
docker run -p 8080:8080 \
  -v $(pwd)/config.yaml:/root/config.yaml \
  ghcr.io/d0ugal/mqtt-exporter:latest --config /root/config.yaml
```

### Building from Source

```bash
# Clone the repository
git clone https://github.com/d0ugal/mqtt-exporter.git
cd mqtt-exporter

# Build the application
make build

# Run the exporter
./mqtt-exporter
```

## Configuration

Copy the example configuration and customize it for your environment:

```bash
cp config.example.yaml config.yaml
```

### Configuration Options

```yaml
server:
  host: "0.0.0.0"  # Server host
  port: 8080       # Server port

logging:
  level: "info"    # Log level: debug, info, warn, error
  format: "json"   # Log format: json, text

metrics:
  collection:
    default_interval: "30s"  # Default collection interval

mqtt:
  broker: "localhost:1883"   # MQTT broker address
  client_id: "mqtt-exporter" # MQTT client ID
  username: ""               # MQTT username (optional)
  password: ""               # MQTT password (optional)
  topics:
    - "#"                    # Topics to monitor (wildcards supported)
  qos: 1                     # Quality of Service level
  clean_session: true        # Clean session on connect
  keep_alive: 60            # Keep alive interval in seconds
  connect_timeout: 30        # Connection timeout in seconds
```

### Environment Variable Configuration

For containerized deployments, you can configure the application entirely through environment variables. This is especially useful for Kubernetes, Docker Compose, and other container orchestration systems.

#### Automatic Detection

The application automatically detects when environment variables are present and uses them for configuration. You don't need to explicitly set `MQTT_EXPORTER_CONFIG_FROM_ENV=true` unless you want to force environment-only mode.

**Behavior:**
- If any `MQTT_EXPORTER_*` environment variables are detected, the application will use environment variable configuration
- If no environment variables are found, it falls back to the YAML configuration file
- The `MQTT_EXPORTER_CONFIG_FROM_ENV=true` flag can be used to explicitly force environment-only mode

#### Environment Variable Format

All environment variables are prefixed with `MQTT_EXPORTER_`:

- `MQTT_EXPORTER_CONFIG_FROM_ENV=true` - Force environment-only configuration mode (optional)
- `MQTT_EXPORTER_SERVER_HOST` - Server host (default: "0.0.0.0")
- `MQTT_EXPORTER_SERVER_PORT` - Server port (default: 8080)
- `MQTT_EXPORTER_LOG_LEVEL` - Log level: debug, info, warn, error (default: "info")
- `MQTT_EXPORTER_LOG_FORMAT` - Log format: json, text (default: "json")
- `MQTT_EXPORTER_METRICS_DEFAULT_INTERVAL` - Default collection interval (default: "30s")

#### MQTT Configuration

- `MQTT_EXPORTER_MQTT_BROKER` - MQTT broker address (required)
- `MQTT_EXPORTER_MQTT_CLIENT_ID` - MQTT client ID (default: "mqtt-exporter")
- `MQTT_EXPORTER_MQTT_USERNAME` - MQTT username (optional)
- `MQTT_EXPORTER_MQTT_PASSWORD` - MQTT password (optional)
- `MQTT_EXPORTER_MQTT_TOPICS` - Comma-separated list of topics (default: "#")
- `MQTT_EXPORTER_MQTT_QOS` - Quality of Service level (default: 1)
- `MQTT_EXPORTER_MQTT_CLEAN_SESSION` - Clean session on connect (default: true)
- `MQTT_EXPORTER_MQTT_KEEP_ALIVE` - Keep alive interval in seconds (default: 60)
- `MQTT_EXPORTER_MQTT_CONNECT_TIMEOUT` - Connection timeout in seconds (default: 30)

#### Docker Compose Example

```yaml
version: '3.8'
services:
  mqtt-exporter:
    image: ghcr.io/d0ugal/mqtt-exporter:latest
    ports:
      - "8080:8080"
    environment:
      - MQTT_EXPORTER_MQTT_BROKER=mqtt://localhost:1883
      - MQTT_EXPORTER_MQTT_USERNAME=user
      - MQTT_EXPORTER_MQTT_PASSWORD=pass
      - MQTT_EXPORTER_MQTT_TOPICS=sensor/+/temperature,device/+/status
      - MQTT_EXPORTER_LOG_LEVEL=debug
    restart: unless-stopped
```

#### Kubernetes Example

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: mqtt-exporter
spec:
  replicas: 1
  selector:
    matchLabels:
      app: mqtt-exporter
  template:
    metadata:
      labels:
        app: mqtt-exporter
    spec:
      containers:
      - name: mqtt-exporter
        image: ghcr.io/d0ugal/mqtt-exporter:latest
        ports:
        - containerPort: 8080
        env:
        - name: MQTT_EXPORTER_MQTT_BROKER
          value: "mqtt://mqtt-broker:1883"
        - name: MQTT_EXPORTER_MQTT_USERNAME
          valueFrom:
            secretKeyRef:
              name: mqtt-credentials
              key: username
        - name: MQTT_EXPORTER_MQTT_PASSWORD
          valueFrom:
            secretKeyRef:
              name: mqtt-credentials
              key: password
        - name: MQTT_EXPORTER_MQTT_TOPICS
          value: "sensor/+/temperature,device/+/status"
```

#### Command Line Usage

```bash
# Use environment variables (automatically detected)
MQTT_EXPORTER_MQTT_BROKER=localhost:1883 ./mqtt-exporter

# Use environment variables with explicit flag
./mqtt-exporter --config-from-env

# Use environment variable for config path
CONFIG_PATH=/path/to/config.yaml ./mqtt-exporter

# Force environment-only mode via environment variable
MQTT_EXPORTER_CONFIG_FROM_ENV=true ./mqtt-exporter
```

## Metrics

The following Prometheus metrics are exposed at `/metrics`:

- `mqtt_messages_total`: Total number of MQTT messages received (by topic)
- `mqtt_message_bytes_total`: Total bytes received in MQTT messages (by topic)
- `mqtt_connection_status`: MQTT connection status (1 = connected, 0 = disconnected)
- `mqtt_connection_errors_total`: Total number of MQTT connection errors
- `mqtt_topic_last_message_timestamp`: Timestamp of the last message received per topic

## Endpoints

- `/`: Service information
- `/health`: Health check endpoint
- `/metrics`: Prometheus metrics endpoint

## Development

### Prerequisites

- Go 1.24 or later
- Docker (for linting and testing)

### Setup

```bash
# Install dependencies
go mod download

# Run tests
make test

# Format and lint code
make lint

# Build the application
make build
```

### Available Make Targets

- `make help`: Show available targets
- `make build`: Build the application
- `make test`: Run tests with coverage
- `make lint`: Format code and run golangci-lint
- `make fmt`: Format code only
- `make lint-only`: Run linting only
- `make clean`: Clean build artifacts

### Testing

```bash
# Run all tests
make test

# Run tests with verbose output
go test -v ./...

# Run tests with race detection
go test -race ./...
```

## Docker Development

```bash
# Build Docker image
docker build -t mqtt-exporter .

# Run container
docker run -p 8080:8080 mqtt-exporter

# Run with custom config (mount config file)
docker run -p 8080:8080 \
  -v $(pwd)/config.yaml:/root/config.yaml \
  mqtt-exporter --config /root/config.yaml
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Code Style

- Follow Go best practices and conventions
- Run `make lint` before submitting PRs
- Ensure all tests pass
- Add tests for new features

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Security

If you discover any security-related issues, please email security@example.com instead of using the issue tracker.

## Support

- Create an issue for bug reports or feature requests
- Check the [documentation](docs/) for detailed guides
- Join our [Discussions](https://github.com/d0ugal/mqtt-exporter/discussions) for questions and community support
