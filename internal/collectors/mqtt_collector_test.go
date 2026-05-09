package collectors

import (
	"testing"

	"github.com/d0ugal/mqtt-exporter/internal/metrics"
	promexporter_config "github.com/d0ugal/promexporter/config"
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

// TestMQTTPassword_RedactsInString locks in that the broker password is
// stored as a SensitiveString whose String() returns "[REDACTED]" rather
// than the raw value. Without this, anything that prints the config would
// leak the broker password.
func TestMQTTPassword_RedactsInString(t *testing.T) {
	const secret = "supersecretmqttpw"

	pw := promexporter_config.NewSensitiveString(secret)

	if pw.Value() != secret {
		t.Fatalf("Value() did not round-trip: want %q, got %q", secret, pw.Value())
	}

	if got := pw.String(); got == secret {
		t.Fatalf("String() leaked the raw password: %q", got)
	}

	if got := pw.String(); got != "[REDACTED]" {
		t.Fatalf("String() unexpected: want [REDACTED], got %q", got)
	}
}
