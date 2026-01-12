//nolint:dupl // Test functions intentionally follow similar patterns for consistency
package collectors

import (
	"testing"

	"github.com/d0ugal/mqtt-exporter/internal/config"
	"github.com/d0ugal/mqtt-exporter/internal/metrics"
	promexporter_metrics "github.com/d0ugal/promexporter/metrics"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestCollector creates a test MQTT collector with initialized metrics
func setupTestCollector(t *testing.T) *MQTTCollector {
	t.Helper()

	baseRegistry := promexporter_metrics.NewRegistry("mqtt_exporter_test")
	mqttRegistry := metrics.NewMQTTRegistry(baseRegistry)

	return &MQTTCollector{
		metrics: mqttRegistry,
		config: &config.Config{
			MQTT: config.MQTTConfig{
				Broker: "test-broker",
			},
		},
		sysCounters: make(map[string]float64),
	}
}

// getGaugeValue retrieves the value from a gauge metric
func getGaugeValue(t *testing.T, gauge *prometheus.GaugeVec, labels prometheus.Labels) float64 {
	t.Helper()

	metric := &dto.Metric{}
	gaugeWithLabels, err := gauge.GetMetricWith(labels)
	require.NoError(t, err)
	require.NoError(t, gaugeWithLabels.Write(metric))

	return metric.GetGauge().GetValue()
}

// getCounterValue retrieves the value from a counter metric
func getCounterValue(t *testing.T, counter *prometheus.CounterVec, labels prometheus.Labels) float64 {
	t.Helper()

	metric := &dto.Metric{}
	counterWithLabels, err := counter.GetMetricWith(labels)
	require.NoError(t, err)
	require.NoError(t, counterWithLabels.Write(metric))

	return metric.GetCounter().GetValue()
}

func TestProcessLoadAverageMetrics(t *testing.T) {
	tests := []struct {
		name          string
		topic         string
		broker        string
		value         float64
		expectedFound bool
		checkFunc     func(t *testing.T, mc *MQTTCollector)
	}{
		{
			name:          "load connections 1min",
			topic:         "$SYS/broker/load/connections/1min",
			broker:        "test-broker",
			value:         5.5,
			expectedFound: true,
			checkFunc: func(t *testing.T, mc *MQTTCollector) {
				labels := prometheus.Labels{"broker": "test-broker", "interval": "1min"}
				val := getGaugeValue(t, mc.metrics.MQTTSysBrokerLoadConnections, labels)
				assert.Equal(t, 5.5, val)
			},
		},
		{
			name:          "load bytes received 5min",
			topic:         "$SYS/broker/load/bytes/received/5min",
			broker:        "test-broker",
			value:         1024.0,
			expectedFound: true,
			checkFunc: func(t *testing.T, mc *MQTTCollector) {
				labels := prometheus.Labels{"broker": "test-broker", "interval": "5min"}
				val := getGaugeValue(t, mc.metrics.MQTTSysBrokerLoadBytesReceived, labels)
				assert.Equal(t, 1024.0, val)
			},
		},
		{
			name:          "load messages sent 15min",
			topic:         "$SYS/broker/load/messages/sent/15min",
			broker:        "test-broker",
			value:         100.0,
			expectedFound: true,
			checkFunc: func(t *testing.T, mc *MQTTCollector) {
				labels := prometheus.Labels{"broker": "test-broker", "interval": "15min"}
				val := getGaugeValue(t, mc.metrics.MQTTSysBrokerLoadMessagesSent, labels)
				assert.Equal(t, 100.0, val)
			},
		},
		{
			name:          "non-load topic",
			topic:         "$SYS/broker/clients/connected",
			broker:        "test-broker",
			value:         10.0,
			expectedFound: false,
		},
		{
			name:          "load topic without interval",
			topic:         "$SYS/broker/load/connections",
			broker:        "test-broker",
			value:         5.0,
			expectedFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := setupTestCollector(t)
			found := mc.processLoadAverageMetrics(tt.topic, tt.broker, tt.value)
			assert.Equal(t, tt.expectedFound, found)

			if tt.checkFunc != nil {
				tt.checkFunc(t, mc)
			}
		})
	}
}

func TestProcessClientMetrics(t *testing.T) {
	tests := []struct {
		name          string
		topic         string
		value         float64
		expectedFound bool
		checkFunc     func(t *testing.T, mc *MQTTCollector)
	}{
		{
			name:          "clients connected",
			topic:         "$SYS/broker/clients/connected",
			value:         10.0,
			expectedFound: true,
			checkFunc: func(t *testing.T, mc *MQTTCollector) {
				labels := prometheus.Labels{"broker": "test-broker"}
				val := getGaugeValue(t, mc.metrics.MQTTSysBrokerClientsConnected, labels)
				assert.Equal(t, 10.0, val)
			},
		},
		{
			name:          "clients disconnected",
			topic:         "$SYS/broker/clients/disconnected",
			value:         3.0,
			expectedFound: true,
			checkFunc: func(t *testing.T, mc *MQTTCollector) {
				labels := prometheus.Labels{"broker": "test-broker"}
				val := getGaugeValue(t, mc.metrics.MQTTSysBrokerClientsDisconnected, labels)
				assert.Equal(t, 3.0, val)
			},
		},
		{
			name:          "clients total",
			topic:         "$SYS/broker/clients/total",
			value:         13.0,
			expectedFound: true,
			checkFunc: func(t *testing.T, mc *MQTTCollector) {
				labels := prometheus.Labels{"broker": "test-broker"}
				val := getGaugeValue(t, mc.metrics.MQTTSysBrokerClientsTotal, labels)
				assert.Equal(t, 13.0, val)
			},
		},
		{
			name:          "non-client metric",
			topic:         "$SYS/broker/messages/received",
			value:         100.0,
			expectedFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := setupTestCollector(t)
			labels := prometheus.Labels{"broker": "test-broker"}
			found := mc.processClientMetrics(tt.topic, labels, tt.value)
			assert.Equal(t, tt.expectedFound, found)

			if tt.checkFunc != nil {
				tt.checkFunc(t, mc)
			}
		})
	}
}

func TestProcessMessageMetrics(t *testing.T) {
	tests := []struct {
		name          string
		topic         string
		value         float64
		expectedFound bool
		checkFunc     func(t *testing.T, mc *MQTTCollector)
	}{
		{
			name:          "messages received",
			topic:         "$SYS/broker/messages/received",
			value:         500.0,
			expectedFound: true,
			checkFunc: func(t *testing.T, mc *MQTTCollector) {
				labels := prometheus.Labels{"broker": "test-broker"}
				val := getCounterValue(t, mc.metrics.MQTTSysBrokerMessagesReceived, labels)
				assert.Equal(t, 500.0, val)
			},
		},
		{
			name:          "messages sent",
			topic:         "$SYS/broker/messages/sent",
			value:         450.0,
			expectedFound: true,
			checkFunc: func(t *testing.T, mc *MQTTCollector) {
				labels := prometheus.Labels{"broker": "test-broker"}
				val := getCounterValue(t, mc.metrics.MQTTSysBrokerMessagesSent, labels)
				assert.Equal(t, 450.0, val)
			},
		},
		{
			name:          "messages inflight",
			topic:         "$SYS/broker/messages/inflight",
			value:         5.0,
			expectedFound: true,
			checkFunc: func(t *testing.T, mc *MQTTCollector) {
				labels := prometheus.Labels{"broker": "test-broker"}
				val := getGaugeValue(t, mc.metrics.MQTTSysBrokerMessagesInflight, labels)
				assert.Equal(t, 5.0, val)
			},
		},
		{
			name:          "store messages count",
			topic:         "$SYS/broker/store/messages/count",
			value:         25.0,
			expectedFound: true,
			checkFunc: func(t *testing.T, mc *MQTTCollector) {
				labels := prometheus.Labels{"broker": "test-broker"}
				val := getGaugeValue(t, mc.metrics.MQTTSysBrokerStoreMessagesCount, labels)
				assert.Equal(t, 25.0, val)
			},
		},
		{
			name:          "non-message metric",
			topic:         "$SYS/broker/clients/connected",
			value:         10.0,
			expectedFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := setupTestCollector(t)
			labels := prometheus.Labels{"broker": "test-broker"}
			// For counter metrics, call twice: first to establish baseline, second to record delta
			mc.processMessageMetrics(tt.topic, labels, 0)
			found := mc.processMessageMetrics(tt.topic, labels, tt.value)
			assert.Equal(t, tt.expectedFound, found)

			if tt.checkFunc != nil {
				tt.checkFunc(t, mc)
			}
		})
	}
}

func TestProcessByteMetrics(t *testing.T) {
	tests := []struct {
		name          string
		topic         string
		value         float64
		expectedFound bool
		checkFunc     func(t *testing.T, mc *MQTTCollector)
	}{
		{
			name:          "bytes received",
			topic:         "$SYS/broker/bytes/received",
			value:         1024000.0,
			expectedFound: true,
			checkFunc: func(t *testing.T, mc *MQTTCollector) {
				labels := prometheus.Labels{"broker": "test-broker"}
				val := getCounterValue(t, mc.metrics.MQTTSysBrokerBytesReceived, labels)
				assert.Equal(t, 1024000.0, val)
			},
		},
		{
			name:          "bytes sent",
			topic:         "$SYS/broker/bytes/sent",
			value:         2048000.0,
			expectedFound: true,
			checkFunc: func(t *testing.T, mc *MQTTCollector) {
				labels := prometheus.Labels{"broker": "test-broker"}
				val := getCounterValue(t, mc.metrics.MQTTSysBrokerBytesSent, labels)
				assert.Equal(t, 2048000.0, val)
			},
		},
		{
			name:          "non-byte metric",
			topic:         "$SYS/broker/messages/received",
			value:         100.0,
			expectedFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := setupTestCollector(t)
			labels := prometheus.Labels{"broker": "test-broker"}
			// For counter metrics, call twice: first to establish baseline, second to record delta
			mc.processByteMetrics(tt.topic, labels, 0)
			found := mc.processByteMetrics(tt.topic, labels, tt.value)
			assert.Equal(t, tt.expectedFound, found)

			if tt.checkFunc != nil {
				tt.checkFunc(t, mc)
			}
		})
	}
}

func TestProcessPublishMetrics(t *testing.T) {
	tests := []struct {
		name          string
		topic         string
		value         float64
		expectedFound bool
		checkFunc     func(t *testing.T, mc *MQTTCollector)
	}{
		{
			name:          "publish messages dropped",
			topic:         "$SYS/broker/publish/messages/dropped",
			value:         2.0,
			expectedFound: true,
			checkFunc: func(t *testing.T, mc *MQTTCollector) {
				labels := prometheus.Labels{"broker": "test-broker"}
				val := getCounterValue(t, mc.metrics.MQTTSysBrokerPublishDropped, labels)
				assert.Equal(t, 2.0, val)
			},
		},
		{
			name:          "publish messages received",
			topic:         "$SYS/broker/publish/messages/received",
			value:         250.0,
			expectedFound: true,
			checkFunc: func(t *testing.T, mc *MQTTCollector) {
				labels := prometheus.Labels{"broker": "test-broker"}
				val := getCounterValue(t, mc.metrics.MQTTSysBrokerPublishReceived, labels)
				assert.Equal(t, 250.0, val)
			},
		},
		{
			name:          "publish messages sent",
			topic:         "$SYS/broker/publish/messages/sent",
			value:         240.0,
			expectedFound: true,
			checkFunc: func(t *testing.T, mc *MQTTCollector) {
				labels := prometheus.Labels{"broker": "test-broker"}
				val := getCounterValue(t, mc.metrics.MQTTSysBrokerPublishSent, labels)
				assert.Equal(t, 240.0, val)
			},
		},
		{
			name:          "non-publish metric",
			topic:         "$SYS/broker/messages/received",
			value:         100.0,
			expectedFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := setupTestCollector(t)
			labels := prometheus.Labels{"broker": "test-broker"}
			// For counter metrics, call twice: first to establish baseline, second to record delta
			mc.processPublishMetrics(tt.topic, labels, 0)
			found := mc.processPublishMetrics(tt.topic, labels, tt.value)
			assert.Equal(t, tt.expectedFound, found)

			if tt.checkFunc != nil {
				tt.checkFunc(t, mc)
			}
		})
	}
}

func TestProcessSubscriptionMetrics(t *testing.T) {
	tests := []struct {
		name          string
		topic         string
		value         float64
		expectedFound bool
		checkFunc     func(t *testing.T, mc *MQTTCollector)
	}{
		{
			name:          "subscriptions count",
			topic:         "$SYS/broker/subscriptions/count",
			value:         15.0,
			expectedFound: true,
			checkFunc: func(t *testing.T, mc *MQTTCollector) {
				labels := prometheus.Labels{"broker": "test-broker"}
				val := getGaugeValue(t, mc.metrics.MQTTSysBrokerSubscriptionsCount, labels)
				assert.Equal(t, 15.0, val)
			},
		},
		{
			name:          "retained messages count",
			topic:         "$SYS/broker/retained/messages/count",
			value:         8.0,
			expectedFound: true,
			checkFunc: func(t *testing.T, mc *MQTTCollector) {
				labels := prometheus.Labels{"broker": "test-broker"}
				val := getGaugeValue(t, mc.metrics.MQTTSysBrokerRetainedMessagesCount, labels)
				assert.Equal(t, 8.0, val)
			},
		},
		{
			name:          "non-subscription metric",
			topic:         "$SYS/broker/messages/received",
			value:         100.0,
			expectedFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := setupTestCollector(t)
			labels := prometheus.Labels{"broker": "test-broker"}
			found := mc.processSubscriptionMetrics(tt.topic, labels, tt.value)
			assert.Equal(t, tt.expectedFound, found)

			if tt.checkFunc != nil {
				tt.checkFunc(t, mc)
			}
		})
	}
}

func TestProcessHeapMetrics(t *testing.T) {
	tests := []struct {
		name          string
		topic         string
		value         float64
		expectedFound bool
		checkFunc     func(t *testing.T, mc *MQTTCollector)
	}{
		{
			name:          "heap current",
			topic:         "$SYS/broker/heap/current",
			value:         52428800.0,
			expectedFound: true,
			checkFunc: func(t *testing.T, mc *MQTTCollector) {
				labels := prometheus.Labels{"broker": "test-broker"}
				val := getGaugeValue(t, mc.metrics.MQTTSysBrokerHeapCurrentBytes, labels)
				assert.Equal(t, 52428800.0, val)
			},
		},
		{
			name:          "heap maximum",
			topic:         "$SYS/broker/heap/maximum",
			value:         104857600.0,
			expectedFound: true,
			checkFunc: func(t *testing.T, mc *MQTTCollector) {
				labels := prometheus.Labels{"broker": "test-broker"}
				val := getGaugeValue(t, mc.metrics.MQTTSysBrokerHeapMaximumBytes, labels)
				assert.Equal(t, 104857600.0, val)
			},
		},
		{
			name:          "non-heap metric",
			topic:         "$SYS/broker/messages/received",
			value:         100.0,
			expectedFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := setupTestCollector(t)
			labels := prometheus.Labels{"broker": "test-broker"}
			found := mc.processHeapMetrics(tt.topic, labels, tt.value)
			assert.Equal(t, tt.expectedFound, found)

			if tt.checkFunc != nil {
				tt.checkFunc(t, mc)
			}
		})
	}
}
