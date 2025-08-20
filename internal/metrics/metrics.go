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

	// Metric information for UI
	metricInfo []MetricInfo
}

// MetricInfo contains information about a metric for the UI
type MetricInfo struct {
	Name         string
	Help         string
	Labels       []string
	ExampleValue string
}

// addMetricInfo adds metric information to the registry
func (r *Registry) addMetricInfo(name, help string, labels []string) {
	r.metricInfo = append(r.metricInfo, MetricInfo{
		Name:         name,
		Help:         help,
		Labels:       labels,
		ExampleValue: "",
	})
}

// NewRegistry creates a new metrics registry
func NewRegistry() *Registry {
	r := &Registry{}

	r.VersionInfo = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "mqtt_exporter_info",
			Help: "Information about the MQTT exporter",
		},
		[]string{"version", "commit", "build_date"},
	)
	r.addMetricInfo("mqtt_exporter_info", "Information about the MQTT exporter", []string{"version", "commit", "build_date"})

	r.MQTTMessageCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "mqtt_messages_total",
			Help: "Total number of MQTT messages received",
		},
		[]string{"topic"},
	)
	r.addMetricInfo("mqtt_messages_total", "Total number of MQTT messages received", []string{"topic"})

	r.MQTTMessageBytes = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "mqtt_message_bytes_total",
			Help: "Total number of bytes received in MQTT messages",
		},
		[]string{"topic"},
	)
	r.addMetricInfo("mqtt_message_bytes_total", "Total number of bytes received in MQTT messages", []string{"topic"})

	r.MQTTConnectionStatus = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "mqtt_connection_status",
			Help: "MQTT connection status (1 = connected, 0 = disconnected)",
		},
		[]string{"broker"},
	)
	r.addMetricInfo("mqtt_connection_status", "MQTT connection status (1 = connected, 0 = disconnected)", []string{"broker"})

	r.MQTTConnectionErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "mqtt_connection_errors_total",
			Help: "Total number of MQTT connection errors",
		},
		[]string{"broker", "error_type"},
	)
	r.addMetricInfo("mqtt_connection_errors_total", "Total number of MQTT connection errors", []string{"broker", "error_type"})

	r.MQTTReconnectsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "mqtt_reconnects_total",
			Help: "Total number of MQTT reconnection attempts",
		},
		[]string{"broker"},
	)
	r.addMetricInfo("mqtt_reconnects_total", "Total number of MQTT reconnection attempts", []string{"broker"})

	r.MQTTTopicLastMessage = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "mqtt_topic_last_message_timestamp",
			Help: "Timestamp of the last message received per topic",
		},
		[]string{"topic"},
	)
	r.addMetricInfo("mqtt_topic_last_message_timestamp", "Timestamp of the last message received per topic", []string{"topic"})

	return r
}

// GetMetricsInfo returns information about all metrics for the UI
func (r *Registry) GetMetricsInfo() []MetricInfo {
	return r.metricInfo
}
