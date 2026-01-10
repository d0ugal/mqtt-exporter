package metrics

import (
	promexporter_metrics "github.com/d0ugal/promexporter/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// MQTTRegistry wraps the promexporter registry with MQTT-specific metrics
type MQTTRegistry struct {
	*promexporter_metrics.Registry

	// MQTT message counters
	MQTTMessageCount *prometheus.CounterVec
	MQTTMessageBytes *prometheus.CounterVec

	// MQTT connection metrics
	MQTTConnectionStatus *prometheus.GaugeVec
	MQTTConnectionErrors *prometheus.CounterVec
	MQTTReconnectsTotal  *prometheus.CounterVec

	// MQTT topic metrics
	MQTTTopicLastMessage *prometheus.GaugeVec

	// MQTT $SYS broker metrics - common metrics across brokers
	MQTTSysBrokerClientsConnected      *prometheus.GaugeVec
	MQTTSysBrokerClientsDisconnected   *prometheus.GaugeVec
	MQTTSysBrokerClientsExpired        *prometheus.CounterVec
	MQTTSysBrokerClientsTotal          *prometheus.GaugeVec
	MQTTSysBrokerClientsMaximum        *prometheus.GaugeVec
	MQTTSysBrokerMessagesReceived      *prometheus.CounterVec
	MQTTSysBrokerMessagesSent          *prometheus.CounterVec
	MQTTSysBrokerMessagesInflight      *prometheus.GaugeVec
	MQTTSysBrokerBytesReceived         *prometheus.CounterVec
	MQTTSysBrokerBytesSent             *prometheus.CounterVec
	MQTTSysBrokerStoreMessagesCount    *prometheus.GaugeVec
	MQTTSysBrokerStoreMessagesBytes    *prometheus.GaugeVec
	MQTTSysBrokerSubscriptionsCount    *prometheus.GaugeVec
	MQTTSysBrokerRetainedMessagesCount *prometheus.GaugeVec
	MQTTSysBrokerHeapCurrentBytes      *prometheus.GaugeVec
	MQTTSysBrokerHeapMaximumBytes      *prometheus.GaugeVec
	MQTTSysBrokerPublishDropped        *prometheus.CounterVec
	MQTTSysBrokerPublishReceived       *prometheus.CounterVec
	MQTTSysBrokerPublishSent           *prometheus.CounterVec
	MQTTSysBrokerVersion               *prometheus.GaugeVec
	MQTTSysBrokerLoadConnections       *prometheus.GaugeVec
	MQTTSysBrokerLoadBytesReceived     *prometheus.GaugeVec
	MQTTSysBrokerLoadBytesSent         *prometheus.GaugeVec
	MQTTSysBrokerLoadMessagesReceived  *prometheus.GaugeVec
	MQTTSysBrokerLoadMessagesSent      *prometheus.GaugeVec
	MQTTSysBrokerLoadPublishReceived   *prometheus.GaugeVec
	MQTTSysBrokerLoadPublishSent       *prometheus.GaugeVec
	MQTTSysBrokerLoadPublishDropped    *prometheus.GaugeVec
	MQTTSysBrokerLoadSockets           *prometheus.GaugeVec
}

// NewMQTTRegistry creates a new MQTT metrics registry
func NewMQTTRegistry(baseRegistry *promexporter_metrics.Registry) *MQTTRegistry {
	// Get the underlying Prometheus registry
	promRegistry := baseRegistry.GetRegistry()
	factory := promauto.With(promRegistry)

	mqtt := &MQTTRegistry{
		Registry: baseRegistry,
	}

	// MQTT message counters
	mqtt.MQTTMessageCount = factory.NewCounterVec(
		prometheus.CounterOpts{
			Name: "mqtt_messages_total",
			Help: "Total number of MQTT messages received",
		},
		[]string{"topic"},
	)
	baseRegistry.AddMetricInfo("mqtt_messages_total", "Total number of MQTT messages received", []string{"topic"})

	mqtt.MQTTMessageBytes = factory.NewCounterVec(
		prometheus.CounterOpts{
			Name: "mqtt_message_bytes_total",
			Help: "Total number of bytes received in MQTT messages",
		},
		[]string{"topic"},
	)
	baseRegistry.AddMetricInfo("mqtt_message_bytes_total", "Total number of bytes received in MQTT messages", []string{"topic"})

	// MQTT connection metrics
	mqtt.MQTTConnectionStatus = factory.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "mqtt_connection_status",
			Help: "MQTT connection status (1 = connected, 0 = disconnected)",
		},
		[]string{"broker"},
	)
	baseRegistry.AddMetricInfo("mqtt_connection_status", "MQTT connection status (1 = connected, 0 = disconnected)", []string{"broker"})

	mqtt.MQTTConnectionErrors = factory.NewCounterVec(
		prometheus.CounterOpts{
			Name: "mqtt_connection_errors_total",
			Help: "Total number of MQTT connection errors",
		},
		[]string{"broker", "error_type"},
	)
	baseRegistry.AddMetricInfo("mqtt_connection_errors_total", "Total number of MQTT connection errors", []string{"broker", "error_type"})

	mqtt.MQTTReconnectsTotal = factory.NewCounterVec(
		prometheus.CounterOpts{
			Name: "mqtt_reconnects_total",
			Help: "Total number of MQTT reconnection attempts",
		},
		[]string{"broker"},
	)
	baseRegistry.AddMetricInfo("mqtt_reconnects_total", "Total number of MQTT reconnection attempts", []string{"broker"})

	// MQTT topic metrics
	mqtt.MQTTTopicLastMessage = factory.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "mqtt_topic_last_message_timestamp",
			Help: "Unix timestamp of the last message received on each topic",
		},
		[]string{"topic"},
	)
	baseRegistry.AddMetricInfo("mqtt_topic_last_message_timestamp", "Unix timestamp of the last message received on each topic", []string{"topic"})

	// MQTT $SYS broker metrics - Client metrics
	mqtt.MQTTSysBrokerClientsConnected = factory.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "mqtt_sys_broker_clients_connected",
			Help: "Number of currently connected clients",
		},
		[]string{"broker"},
	)
	baseRegistry.AddMetricInfo("mqtt_sys_broker_clients_connected", "Number of currently connected clients", []string{"broker"})

	mqtt.MQTTSysBrokerClientsDisconnected = factory.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "mqtt_sys_broker_clients_disconnected",
			Help: "Total number of persistent clients currently disconnected",
		},
		[]string{"broker"},
	)
	baseRegistry.AddMetricInfo("mqtt_sys_broker_clients_disconnected", "Total number of persistent clients currently disconnected", []string{"broker"})

	mqtt.MQTTSysBrokerClientsExpired = factory.NewCounterVec(
		prometheus.CounterOpts{
			Name: "mqtt_sys_broker_clients_expired_total",
			Help: "Number of disconnected persistent clients that have been expired and removed",
		},
		[]string{"broker"},
	)
	baseRegistry.AddMetricInfo("mqtt_sys_broker_clients_expired_total", "Number of disconnected persistent clients that have been expired and removed", []string{"broker"})

	mqtt.MQTTSysBrokerClientsTotal = factory.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "mqtt_sys_broker_clients_total",
			Help: "Total number of active and inactive clients",
		},
		[]string{"broker"},
	)
	baseRegistry.AddMetricInfo("mqtt_sys_broker_clients_total", "Total number of active and inactive clients", []string{"broker"})

	mqtt.MQTTSysBrokerClientsMaximum = factory.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "mqtt_sys_broker_clients_maximum",
			Help: "Maximum number of clients connected simultaneously",
		},
		[]string{"broker"},
	)
	baseRegistry.AddMetricInfo("mqtt_sys_broker_clients_maximum", "Maximum number of clients connected simultaneously", []string{"broker"})

	// Message metrics
	mqtt.MQTTSysBrokerMessagesReceived = factory.NewCounterVec(
		prometheus.CounterOpts{
			Name: "mqtt_sys_broker_messages_received_total",
			Help: "Total number of messages of any type received since startup",
		},
		[]string{"broker"},
	)
	baseRegistry.AddMetricInfo("mqtt_sys_broker_messages_received_total", "Total number of messages of any type received since startup", []string{"broker"})

	mqtt.MQTTSysBrokerMessagesSent = factory.NewCounterVec(
		prometheus.CounterOpts{
			Name: "mqtt_sys_broker_messages_sent_total",
			Help: "Total number of messages of any type sent since startup",
		},
		[]string{"broker"},
	)
	baseRegistry.AddMetricInfo("mqtt_sys_broker_messages_sent_total", "Total number of messages of any type sent since startup", []string{"broker"})

	mqtt.MQTTSysBrokerMessagesInflight = factory.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "mqtt_sys_broker_messages_inflight",
			Help: "Number of messages with QoS>0 awaiting acknowledgments",
		},
		[]string{"broker"},
	)
	baseRegistry.AddMetricInfo("mqtt_sys_broker_messages_inflight", "Number of messages with QoS>0 awaiting acknowledgments", []string{"broker"})

	// Byte metrics
	mqtt.MQTTSysBrokerBytesReceived = factory.NewCounterVec(
		prometheus.CounterOpts{
			Name: "mqtt_sys_broker_bytes_received_total",
			Help: "Total number of bytes received since the broker started",
		},
		[]string{"broker"},
	)
	baseRegistry.AddMetricInfo("mqtt_sys_broker_bytes_received_total", "Total number of bytes received since the broker started", []string{"broker"})

	mqtt.MQTTSysBrokerBytesSent = factory.NewCounterVec(
		prometheus.CounterOpts{
			Name: "mqtt_sys_broker_bytes_sent_total",
			Help: "Total number of bytes sent since the broker started",
		},
		[]string{"broker"},
	)
	baseRegistry.AddMetricInfo("mqtt_sys_broker_bytes_sent_total", "Total number of bytes sent since the broker started", []string{"broker"})

	// Store metrics
	mqtt.MQTTSysBrokerStoreMessagesCount = factory.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "mqtt_sys_broker_store_messages_count",
			Help: "Number of messages currently held in the message store",
		},
		[]string{"broker"},
	)
	baseRegistry.AddMetricInfo("mqtt_sys_broker_store_messages_count", "Number of messages currently held in the message store", []string{"broker"})

	mqtt.MQTTSysBrokerStoreMessagesBytes = factory.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "mqtt_sys_broker_store_messages_bytes",
			Help: "Bytes currently held by message payloads in the message store",
		},
		[]string{"broker"},
	)
	baseRegistry.AddMetricInfo("mqtt_sys_broker_store_messages_bytes", "Bytes currently held by message payloads in the message store", []string{"broker"})

	// Subscription and retained message metrics
	mqtt.MQTTSysBrokerSubscriptionsCount = factory.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "mqtt_sys_broker_subscriptions_count",
			Help: "Total number of subscriptions active on the broker",
		},
		[]string{"broker"},
	)
	baseRegistry.AddMetricInfo("mqtt_sys_broker_subscriptions_count", "Total number of subscriptions active on the broker", []string{"broker"})

	mqtt.MQTTSysBrokerRetainedMessagesCount = factory.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "mqtt_sys_broker_retained_messages_count",
			Help: "Total number of retained messages active on the broker",
		},
		[]string{"broker"},
	)
	baseRegistry.AddMetricInfo("mqtt_sys_broker_retained_messages_count", "Total number of retained messages active on the broker", []string{"broker"})

	// Heap memory metrics
	mqtt.MQTTSysBrokerHeapCurrentBytes = factory.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "mqtt_sys_broker_heap_current_bytes",
			Help: "Current heap memory size in bytes",
		},
		[]string{"broker"},
	)
	baseRegistry.AddMetricInfo("mqtt_sys_broker_heap_current_bytes", "Current heap memory size in bytes", []string{"broker"})

	mqtt.MQTTSysBrokerHeapMaximumBytes = factory.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "mqtt_sys_broker_heap_maximum_bytes",
			Help: "Largest amount of heap memory used by the broker",
		},
		[]string{"broker"},
	)
	baseRegistry.AddMetricInfo("mqtt_sys_broker_heap_maximum_bytes", "Largest amount of heap memory used by the broker", []string{"broker"})

	// Publish metrics
	mqtt.MQTTSysBrokerPublishDropped = factory.NewCounterVec(
		prometheus.CounterOpts{
			Name: "mqtt_sys_broker_publish_dropped_total",
			Help: "Total number of publish messages dropped due to inflight/queuing limits",
		},
		[]string{"broker"},
	)
	baseRegistry.AddMetricInfo("mqtt_sys_broker_publish_dropped_total", "Total number of publish messages dropped due to inflight/queuing limits", []string{"broker"})

	mqtt.MQTTSysBrokerPublishReceived = factory.NewCounterVec(
		prometheus.CounterOpts{
			Name: "mqtt_sys_broker_publish_received_total",
			Help: "Total number of PUBLISH messages received since startup",
		},
		[]string{"broker"},
	)
	baseRegistry.AddMetricInfo("mqtt_sys_broker_publish_received_total", "Total number of PUBLISH messages received since startup", []string{"broker"})

	mqtt.MQTTSysBrokerPublishSent = factory.NewCounterVec(
		prometheus.CounterOpts{
			Name: "mqtt_sys_broker_publish_sent_total",
			Help: "Total number of PUBLISH messages sent since startup",
		},
		[]string{"broker"},
	)
	baseRegistry.AddMetricInfo("mqtt_sys_broker_publish_sent_total", "Total number of PUBLISH messages sent since startup", []string{"broker"})

	// Version metric (set to 1, actual version in label)
	mqtt.MQTTSysBrokerVersion = factory.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "mqtt_sys_broker_version_info",
			Help: "Broker version information",
		},
		[]string{"broker", "version"},
	)
	baseRegistry.AddMetricInfo("mqtt_sys_broker_version_info", "Broker version information", []string{"broker", "version"})

	// Load average metrics (1min, 5min, 15min)
	mqtt.MQTTSysBrokerLoadConnections = factory.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "mqtt_sys_broker_load_connections",
			Help: "Moving average of CONNECT packets received",
		},
		[]string{"broker", "interval"},
	)
	baseRegistry.AddMetricInfo("mqtt_sys_broker_load_connections", "Moving average of CONNECT packets received", []string{"broker", "interval"})

	mqtt.MQTTSysBrokerLoadBytesReceived = factory.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "mqtt_sys_broker_load_bytes_received",
			Help: "Moving average of bytes received",
		},
		[]string{"broker", "interval"},
	)
	baseRegistry.AddMetricInfo("mqtt_sys_broker_load_bytes_received", "Moving average of bytes received", []string{"broker", "interval"})

	mqtt.MQTTSysBrokerLoadBytesSent = factory.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "mqtt_sys_broker_load_bytes_sent",
			Help: "Moving average of bytes sent",
		},
		[]string{"broker", "interval"},
	)
	baseRegistry.AddMetricInfo("mqtt_sys_broker_load_bytes_sent", "Moving average of bytes sent", []string{"broker", "interval"})

	mqtt.MQTTSysBrokerLoadMessagesReceived = factory.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "mqtt_sys_broker_load_messages_received",
			Help: "Moving average of all MQTT messages received",
		},
		[]string{"broker", "interval"},
	)
	baseRegistry.AddMetricInfo("mqtt_sys_broker_load_messages_received", "Moving average of all MQTT messages received", []string{"broker", "interval"})

	mqtt.MQTTSysBrokerLoadMessagesSent = factory.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "mqtt_sys_broker_load_messages_sent",
			Help: "Moving average of all MQTT messages sent",
		},
		[]string{"broker", "interval"},
	)
	baseRegistry.AddMetricInfo("mqtt_sys_broker_load_messages_sent", "Moving average of all MQTT messages sent", []string{"broker", "interval"})

	mqtt.MQTTSysBrokerLoadPublishReceived = factory.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "mqtt_sys_broker_load_publish_received",
			Help: "Moving average of PUBLISH messages received",
		},
		[]string{"broker", "interval"},
	)
	baseRegistry.AddMetricInfo("mqtt_sys_broker_load_publish_received", "Moving average of PUBLISH messages received", []string{"broker", "interval"})

	mqtt.MQTTSysBrokerLoadPublishSent = factory.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "mqtt_sys_broker_load_publish_sent",
			Help: "Moving average of PUBLISH messages sent",
		},
		[]string{"broker", "interval"},
	)
	baseRegistry.AddMetricInfo("mqtt_sys_broker_load_publish_sent", "Moving average of PUBLISH messages sent", []string{"broker", "interval"})

	mqtt.MQTTSysBrokerLoadPublishDropped = factory.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "mqtt_sys_broker_load_publish_dropped",
			Help: "Moving average of dropped PUBLISH messages",
		},
		[]string{"broker", "interval"},
	)
	baseRegistry.AddMetricInfo("mqtt_sys_broker_load_publish_dropped", "Moving average of dropped PUBLISH messages", []string{"broker", "interval"})

	mqtt.MQTTSysBrokerLoadSockets = factory.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "mqtt_sys_broker_load_sockets",
			Help: "Moving average of socket connections opened to the broker",
		},
		[]string{"broker", "interval"},
	)
	baseRegistry.AddMetricInfo("mqtt_sys_broker_load_sockets", "Moving average of socket connections opened to the broker", []string{"broker", "interval"})

	return mqtt
}
