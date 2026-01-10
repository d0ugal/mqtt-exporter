# MQTT Exporter

A Prometheus exporter for MQTT message monitoring that connects to an MQTT broker and exposes metrics about message counts and connection status.

**Image**: `ghcr.io/d0ugal/mqtt-exporter:v1.25.27`

## Metrics

### MQTT Metrics

#### Standard Metrics
- `mqtt_messages_total` - Total number of MQTT messages received (by topic)
- `mqtt_message_bytes_total` - Total bytes received in MQTT messages (by topic)
- `mqtt_connection_status` - MQTT connection status (1 = connected, 0 = disconnected)
- `mqtt_connection_errors_total` - Total number of MQTT connection errors
- `mqtt_reconnects_total` - Total number of MQTT reconnection attempts
- `mqtt_topic_last_message_timestamp` - Timestamp of the last message received per topic

#### Broker $SYS Topic Metrics (enabled by default)
These metrics are automatically collected from the broker's $SYS topics:

**Client Metrics:**
- `mqtt_sys_broker_clients_connected` - Number of currently connected clients
- `mqtt_sys_broker_clients_disconnected` - Total number of persistent clients currently disconnected
- `mqtt_sys_broker_clients_expired_total` - Number of disconnected persistent clients that have been expired and removed
- `mqtt_sys_broker_clients_total` - Total number of active and inactive clients
- `mqtt_sys_broker_clients_maximum` - Maximum number of clients connected simultaneously

**Message Metrics:**
- `mqtt_sys_broker_messages_received_total` - Total messages of any type received since startup
- `mqtt_sys_broker_messages_sent_total` - Total messages of any type sent since startup
- `mqtt_sys_broker_messages_inflight` - Number of messages with QoS>0 awaiting acknowledgments
- `mqtt_sys_broker_store_messages_count` - Number of messages currently held in the message store
- `mqtt_sys_broker_store_messages_bytes` - Bytes currently held by message payloads in the message store

**Byte Metrics:**
- `mqtt_sys_broker_bytes_received_total` - Total bytes received since the broker started
- `mqtt_sys_broker_bytes_sent_total` - Total bytes sent since the broker started

**Publish Metrics:**
- `mqtt_sys_broker_publish_received_total` - Total PUBLISH messages received since startup
- `mqtt_sys_broker_publish_sent_total` - Total PUBLISH messages sent since startup
- `mqtt_sys_broker_publish_dropped_total` - Total publish messages dropped due to inflight/queuing limits

**Other Metrics:**
- `mqtt_sys_broker_subscriptions_count` - Total number of subscriptions active on the broker
- `mqtt_sys_broker_retained_messages_count` - Total number of retained messages active on the broker
- `mqtt_sys_broker_heap_current_bytes` - Current heap memory size in bytes
- `mqtt_sys_broker_heap_maximum_bytes` - Largest amount of heap memory used by the broker
- `mqtt_sys_broker_version_info` - Broker version information (gauge set to 1 with version label)

**Load Average Metrics (1min, 5min, 15min intervals):**
- `mqtt_sys_broker_load_connections` - Moving average of CONNECT packets received
- `mqtt_sys_broker_load_bytes_received` - Moving average of bytes received
- `mqtt_sys_broker_load_bytes_sent` - Moving average of bytes sent
- `mqtt_sys_broker_load_messages_received` - Moving average of all MQTT messages received
- `mqtt_sys_broker_load_messages_sent` - Moving average of all MQTT messages sent
- `mqtt_sys_broker_load_publish_received` - Moving average of PUBLISH messages received
- `mqtt_sys_broker_load_publish_sent` - Moving average of PUBLISH messages sent
- `mqtt_sys_broker_load_publish_dropped` - Moving average of dropped PUBLISH messages
- `mqtt_sys_broker_load_sockets` - Moving average of socket connections opened to the broker

**Note:** $SYS topic monitoring is enabled by default. Some metrics may not be available depending on your MQTT broker implementation (e.g., heap metrics, load averages). Set `disable_sys_topics: true` to disable this feature.

### Endpoints
- `GET /`: Service information
- `GET /health`: Health check endpoint
- `GET /metrics`: Prometheus metrics endpoint

## Quick Start

### Docker Compose

```yaml
version: '3.8'
services:
  mqtt-exporter:
    image: ghcr.io/d0ugal/mqtt-exporter:v1.25.27
    ports:
      - "8080:8080"
    environment:
      - MQTT_EXPORTER_MQTT_BROKER=mqtt://localhost:1883
      - MQTT_EXPORTER_MQTT_TOPICS=#
    restart: unless-stopped
```

1. Update the MQTT broker URL and topics in the environment variables
2. Run: `docker-compose up -d`
3. Access metrics: `curl http://localhost:8080/metrics`

## Configuration

Create a `config.yaml` file to configure MQTT broker connection:

```yaml
server:
  host: "0.0.0.0"
  port: 8080

logging:
  level: "info"
  format: "json"

metrics:
  collection:
    default_interval: "30s"

mqtt:
  broker: "localhost:1883"
  client_id: "mqtt-exporter"
  username: ""
  password: ""
  topics:
    - "#"
  qos: 1
  clean_session: true
  keep_alive: 60
  connect_timeout: 30
  disable_sys_topics: false  # Set to true to disable $SYS topic monitoring
```

## Deployment

### Docker Compose (Environment Variables)

```yaml
version: '3.8'
services:
  mqtt-exporter:
    image: ghcr.io/d0ugal/mqtt-exporter:v1.25.27
    ports:
      - "8080:8080"
    environment:
      - MQTT_EXPORTER_MQTT_BROKER=mqtt://localhost:1883
      - MQTT_EXPORTER_MQTT_USERNAME=user
      - MQTT_EXPORTER_MQTT_PASSWORD=pass
      - MQTT_EXPORTER_MQTT_TOPICS=sensor/+/temperature,device/+/status
    restart: unless-stopped
```

### Kubernetes

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
        image: ghcr.io/d0ugal/mqtt-exporter:v1.25.27
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

## Prometheus Integration

Add to your `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: 'mqtt-exporter'
    static_configs:
      - targets: ['mqtt-exporter:8080']
```

## Environment Variable Configuration

For containerized deployments, you can configure the application entirely through environment variables:

### Environment Variable Format

All environment variables are prefixed with `MQTT_EXPORTER_`:

- `MQTT_EXPORTER_MQTT_BROKER` - MQTT broker address (required)
- `MQTT_EXPORTER_MQTT_CLIENT_ID` - MQTT client ID (default: "mqtt-exporter")
- `MQTT_EXPORTER_MQTT_USERNAME` - MQTT username (optional)
- `MQTT_EXPORTER_MQTT_PASSWORD` - MQTT password (optional)
- `MQTT_EXPORTER_MQTT_TOPICS` - Comma-separated list of topics (default: "#")
- `MQTT_EXPORTER_MQTT_QOS` - Quality of Service level (default: 1)
- `MQTT_EXPORTER_MQTT_DISABLE_SYS_TOPICS` - Disable $SYS topic monitoring (default: false)
- `MQTT_EXPORTER_SERVER_HOST` - Server host (default: "0.0.0.0")
- `MQTT_EXPORTER_SERVER_PORT` - Server port (default: 8080)
- `MQTT_EXPORTER_LOG_LEVEL` - Log level: debug, info, warn, error (default: "info")
- `MQTT_EXPORTER_LOG_FORMAT` - Log format: json, text (default: "json")

## Development

### Building

```bash
make build
```

### Testing

```bash
make test
```

### Linting

```bash
make lint
```

### Formatting

```bash
make fmt
```

## License

This project is licensed under the MIT License.