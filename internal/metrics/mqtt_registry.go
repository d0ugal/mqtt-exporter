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
}

// NewMQTTRegistry creates a new MQTT metrics registry
//
//nolint:maintidx // This function registers many metrics - splitting would reduce maintainability
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

	return mqtt
}
