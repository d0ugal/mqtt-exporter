package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

// TestRegistry holds all the metrics for testing
type TestRegistry struct {
	// MQTT message counters
	MQTTMessageCount *prometheus.CounterVec
	MQTTMessageBytes *prometheus.CounterVec

	// MQTT connection metrics
	MQTTConnectionStatus *prometheus.GaugeVec
	MQTTConnectionErrors *prometheus.CounterVec

	// MQTT topic metrics
	MQTTTopicLastMessage *prometheus.GaugeVec
}

// NewTestRegistry creates a new test metrics registry
func NewTestRegistry() *TestRegistry {
	return &TestRegistry{
		MQTTMessageCount: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "mqtt_messages_total",
				Help: "Total number of MQTT messages received",
			},
			[]string{"topic"},
		),
		MQTTMessageBytes: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "mqtt_message_bytes_total",
				Help: "Total number of bytes received in MQTT messages",
			},
			[]string{"topic"},
		),
		MQTTConnectionStatus: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "mqtt_connection_status",
				Help: "MQTT connection status (1 = connected, 0 = disconnected)",
			},
			[]string{"broker"},
		),
		MQTTConnectionErrors: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "mqtt_connection_errors_total",
				Help: "Total number of MQTT connection errors",
			},
			[]string{"broker", "error_type"},
		),
		MQTTTopicLastMessage: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "mqtt_topic_last_message_timestamp",
				Help: "Timestamp of the last message received per topic",
			},
			[]string{"topic"},
		),
	}
}

func TestNewRegistry(t *testing.T) {
	registry := NewRegistry()

	// Test that all metrics are properly initialized
	if registry.MQTTMessageCount == nil {
		t.Error("MQTTMessageCount should not be nil")
	}

	if registry.MQTTMessageBytes == nil {
		t.Error("MQTTMessageBytes should not be nil")
	}

	if registry.MQTTConnectionStatus == nil {
		t.Error("MQTTConnectionStatus should not be nil")
	}

	if registry.MQTTConnectionErrors == nil {
		t.Error("MQTTConnectionErrors should not be nil")
	}

	if registry.MQTTTopicLastMessage == nil {
		t.Error("MQTTTopicLastMessage should not be nil")
	}
}

func TestMQTTMessageCount(t *testing.T) {
	registry := NewTestRegistry()

	// Test incrementing message count
	registry.MQTTMessageCount.WithLabelValues("test/topic").Inc()
	registry.MQTTMessageCount.WithLabelValues("test/topic").Inc()
	registry.MQTTMessageCount.WithLabelValues("another/topic").Inc()

	// Check the values
	expected := float64(2)

	actual := testutil.ToFloat64(registry.MQTTMessageCount.WithLabelValues("test/topic"))
	if actual != expected {
		t.Errorf("MQTTMessageCount for 'test/topic' = %v, want %v", actual, expected)
	}

	expected = float64(1)

	actual = testutil.ToFloat64(registry.MQTTMessageCount.WithLabelValues("another/topic"))
	if actual != expected {
		t.Errorf("MQTTMessageCount for 'another/topic' = %v, want %v", actual, expected)
	}
}

func TestMQTTMessageBytes(t *testing.T) {
	registry := NewTestRegistry()

	// Test adding message bytes
	registry.MQTTMessageBytes.WithLabelValues("test/topic").Add(100)
	registry.MQTTMessageBytes.WithLabelValues("test/topic").Add(50)
	registry.MQTTMessageBytes.WithLabelValues("another/topic").Add(200)

	// Check the values
	expected := float64(150)

	actual := testutil.ToFloat64(registry.MQTTMessageBytes.WithLabelValues("test/topic"))
	if actual != expected {
		t.Errorf("MQTTMessageBytes for 'test/topic' = %v, want %v", actual, expected)
	}

	expected = float64(200)

	actual = testutil.ToFloat64(registry.MQTTMessageBytes.WithLabelValues("another/topic"))
	if actual != expected {
		t.Errorf("MQTTMessageBytes for 'another/topic' = %v, want %v", actual, expected)
	}
}

func TestMQTTConnectionStatus(t *testing.T) {
	registry := NewTestRegistry()

	// Test setting connection status
	registry.MQTTConnectionStatus.WithLabelValues("localhost:1883").Set(1)
	registry.MQTTConnectionStatus.WithLabelValues("broker2:1883").Set(0)

	// Check the values
	expected := float64(1)

	actual := testutil.ToFloat64(registry.MQTTConnectionStatus.WithLabelValues("localhost:1883"))
	if actual != expected {
		t.Errorf("MQTTConnectionStatus for 'localhost:1883' = %v, want %v", actual, expected)
	}

	expected = float64(0)

	actual = testutil.ToFloat64(registry.MQTTConnectionStatus.WithLabelValues("broker2:1883"))
	if actual != expected {
		t.Errorf("MQTTConnectionStatus for 'broker2:1883' = %v, want %v", actual, expected)
	}
}

func TestMQTTConnectionErrors(t *testing.T) {
	registry := NewTestRegistry()

	// Test incrementing connection errors
	registry.MQTTConnectionErrors.WithLabelValues("localhost:1883", "connect").Inc()
	registry.MQTTConnectionErrors.WithLabelValues("localhost:1883", "connect").Inc()
	registry.MQTTConnectionErrors.WithLabelValues("localhost:1883", "subscribe").Inc()

	// Check the values
	expected := float64(2)

	actual := testutil.ToFloat64(registry.MQTTConnectionErrors.WithLabelValues("localhost:1883", "connect"))
	if actual != expected {
		t.Errorf("MQTTConnectionErrors for 'localhost:1883'/'connect' = %v, want %v", actual, expected)
	}

	expected = float64(1)

	actual = testutil.ToFloat64(registry.MQTTConnectionErrors.WithLabelValues("localhost:1883", "subscribe"))
	if actual != expected {
		t.Errorf("MQTTConnectionErrors for 'localhost:1883'/'subscribe' = %v, want %v", actual, expected)
	}
}

func TestMQTTTopicLastMessage(t *testing.T) {
	registry := NewTestRegistry()

	// Test setting last message timestamp
	timestamp1 := float64(1640995200) // 2022-01-01 00:00:00 UTC
	timestamp2 := float64(1640995260) // 2022-01-01 00:01:00 UTC

	registry.MQTTTopicLastMessage.WithLabelValues("test/topic").Set(timestamp1)
	registry.MQTTTopicLastMessage.WithLabelValues("another/topic").Set(timestamp2)

	// Check the values
	actual := testutil.ToFloat64(registry.MQTTTopicLastMessage.WithLabelValues("test/topic"))
	if actual != timestamp1 {
		t.Errorf("MQTTTopicLastMessage for 'test/topic' = %v, want %v", actual, timestamp1)
	}

	actual = testutil.ToFloat64(registry.MQTTTopicLastMessage.WithLabelValues("another/topic"))
	if actual != timestamp2 {
		t.Errorf("MQTTTopicLastMessage for 'another/topic' = %v, want %v", actual, timestamp2)
	}
}

func TestMetricsRegistration(t *testing.T) {
	registry := NewTestRegistry()

	// Test that metrics are properly created
	metrics := []prometheus.Collector{
		registry.MQTTMessageCount,
		registry.MQTTMessageBytes,
		registry.MQTTConnectionStatus,
		registry.MQTTConnectionErrors,
		registry.MQTTTopicLastMessage,
	}

	for _, metric := range metrics {
		// We're testing that the metrics are created without errors
		if metric == nil {
			t.Error("Metric should not be nil")
		}
	}
}
