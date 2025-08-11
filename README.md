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

# Run with custom config
docker run -p 8080:8080 \
  -v $(pwd)/config.yaml:/app/config.yaml \
  ghcr.io/d0ugal/mqtt-exporter:latest
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

# Run with custom config
docker run -p 8080:8080 \
  -v $(pwd)/config.yaml:/app/config.yaml \
  mqtt-exporter
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
