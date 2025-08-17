package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Registry holds all the metrics for the MQTT exporter
type Registry struct {
	// Version info metric
	VersionInfo *prometheus.GaugeVec

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

// NewRegistry creates a new metrics registry
func NewRegistry() *Registry {
	return &Registry{
		VersionInfo: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "mqtt_exporter_info",
				Help: "Information about the MQTT exporter",
			},
			[]string{"version", "commit", "build_date"},
		),
		MQTTMessageCount: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "mqtt_messages_total",
				Help: "Total number of MQTT messages received",
			},
			[]string{"topic"},
		),
		MQTTMessageBytes: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "mqtt_message_bytes_total",
				Help: "Total number of bytes received in MQTT messages",
			},
			[]string{"topic"},
		),
		MQTTConnectionStatus: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "mqtt_connection_status",
				Help: "MQTT connection status (1 = connected, 0 = disconnected)",
			},
			[]string{"broker"},
		),
		MQTTConnectionErrors: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "mqtt_connection_errors_total",
				Help: "Total number of MQTT connection errors",
			},
			[]string{"broker", "error_type"},
		),
		MQTTReconnectsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "mqtt_reconnects_total",
				Help: "Total number of MQTT reconnection attempts",
			},
			[]string{"broker"},
		),
		MQTTTopicLastMessage: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "mqtt_topic_last_message_timestamp",
				Help: "Timestamp of the last message received per topic",
			},
			[]string{"topic"},
		),
	}
}
