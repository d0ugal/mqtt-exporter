package collectors

import (
	"testing"

	"github.com/d0ugal/mqtt-exporter/internal/metrics"
	promexporter_metrics "github.com/d0ugal/promexporter/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

// TestMQTTConnectionErrors_LabelsMatchRegistry guards against the label-name
// mismatch between the registry definition (broker, error_type) and the call
// sites in run() / onConnectionLost. If the labels drift again,
// prometheus.Labels.With() panics here.
func TestMQTTConnectionErrors_LabelsMatchRegistry(t *testing.T) {
	baseRegistry := promexporter_metrics.NewRegistry("mqtt_exporter_info_test")
	mqttMetrics := metrics.NewMQTTRegistry(baseRegistry)

	assert.NotPanics(t, func() {
		mqttMetrics.MQTTConnectionErrors.With(prometheus.Labels{
			"broker":     "tcp://localhost:1883",
			"error_type": "connect",
		}).Inc()
		mqttMetrics.MQTTConnectionErrors.With(prometheus.Labels{
			"broker":     "tcp://localhost:1883",
			"error_type": "subscribe",
		}).Inc()
		mqttMetrics.MQTTConnectionErrors.With(prometheus.Labels{
			"broker":     "tcp://localhost:1883",
			"error_type": "connection_lost",
		}).Inc()
	})
}
