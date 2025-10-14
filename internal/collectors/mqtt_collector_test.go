package collectors

import (
	"context"
	"testing"
	"time"

	"github.com/d0ugal/mqtt-exporter/internal/config"
	"github.com/d0ugal/mqtt-exporter/internal/metrics"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

// createTestCollector creates a new MQTTCollector for testing
func createTestCollector() *MQTTCollector {
	cfg := &config.Config{
		MQTT: config.MQTTConfig{
			Broker:   "localhost:1883",
			ClientID: "test-client",
		},
	}
	// Use a fresh metrics registry for each test to avoid duplicate registration
	metricsRegistry := &metrics.Registry{}

	return NewMQTTCollector(cfg, metricsRegistry)
}

func TestMQTTCollector_NewMQTTCollector(t *testing.T) {
	collector := createTestCollector()

	if collector == nil {
		t.Fatal("Expected collector to be created")
	}

	if collector.config == nil {
		t.Error("Expected config to be set")
	}

	if collector.metrics == nil {
		t.Error("Expected metrics registry to be set")
	}

	if collector.topics == nil {
		t.Error("Expected topics map to be initialized")
	}

	if collector.done == nil {
		t.Error("Expected done channel to be initialized")
	}
}

func TestMQTTCollector_Stop(t *testing.T) {
	cfg := &config.Config{
		MQTT: config.MQTTConfig{
			Broker:   "localhost:1883",
			ClientID: "test-client",
		},
	}
	// Create a minimal registry without initializing metrics to avoid duplicate registration
	metricsRegistry := &metrics.Registry{}
	collector := NewMQTTCollector(cfg, metricsRegistry)

	// Test that Stop doesn't panic
	collector.Stop()

	// Test that done channel is closed
	select {
	case <-collector.done:
		// Channel is closed, which is expected
	default:
		t.Error("Expected done channel to be closed")
	}
}

func TestMQTTCollector_monitorConnection_ContextCancelled(t *testing.T) {
	cfg := &config.Config{
		MQTT: config.MQTTConfig{
			Broker:   "localhost:1883",
			ClientID: "test-client",
		},
	}
	metricsRegistry := &metrics.Registry{}
	collector := NewMQTTCollector(cfg, metricsRegistry)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	// Should return false (not a connection loss) when context is cancelled
	result := collector.monitorConnection(ctx)
	if result {
		t.Error("Expected monitorConnection to return false when context is cancelled")
	}
}

func TestMQTTCollector_monitorConnection_ShutdownRequested(t *testing.T) {
	cfg := &config.Config{
		MQTT: config.MQTTConfig{
			Broker:   "localhost:1883",
			ClientID: "test-client",
		},
	}
	metricsRegistry := &metrics.Registry{}
	collector := NewMQTTCollector(cfg, metricsRegistry)

	ctx := context.Background()

	// Stop the collector to close the done channel
	collector.Stop()

	// Should return false (not a connection loss) when shutdown is requested
	result := collector.monitorConnection(ctx)
	if result {
		t.Error("Expected monitorConnection to return false when shutdown is requested")
	}
}

func TestMQTTCollector_monitorConnection_ConnectionLost(t *testing.T) {
	cfg := &config.Config{
		MQTT: config.MQTTConfig{
			Broker:   "localhost:1883",
			ClientID: "test-client",
		},
	}
	metricsRegistry := &metrics.Registry{}
	collector := NewMQTTCollector(cfg, metricsRegistry)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Set client to nil to simulate connection loss
	collector.client = nil

	// Should return true (connection lost) when client is nil
	result := collector.monitorConnection(ctx)
	if !result {
		t.Error("Expected monitorConnection to return true when client is nil")
	}
}

func TestMQTTCollector_monitorConnection_ClientDisconnected(t *testing.T) {
	cfg := &config.Config{
		MQTT: config.MQTTConfig{
			Broker:   "localhost:1883",
			ClientID: "test-client",
		},
	}
	metricsRegistry := &metrics.Registry{}
	collector := NewMQTTCollector(cfg, metricsRegistry)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Create a mock client that reports as disconnected
	mockClient := &mockMQTTClient{connected: false}
	collector.client = mockClient

	// Should return true (connection lost) when client is disconnected
	result := collector.monitorConnection(ctx)
	if !result {
		t.Error("Expected monitorConnection to return true when client is disconnected")
	}
}

func TestMQTTCollector_monitorConnection_ClientConnected(t *testing.T) {
	cfg := &config.Config{
		MQTT: config.MQTTConfig{
			Broker:   "localhost:1883",
			ClientID: "test-client",
		},
	}
	metricsRegistry := &metrics.Registry{}
	collector := NewMQTTCollector(cfg, metricsRegistry)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Create a mock client that reports as connected
	mockClient := &mockMQTTClient{connected: true}
	collector.client = mockClient

	// Should return false (not a connection loss) when client is connected
	// This will timeout after 200ms, which is expected behavior
	result := collector.monitorConnection(ctx)
	if result {
		t.Error("Expected monitorConnection to return false when client is connected")
	}
}

// mockMQTTClient is a mock implementation of MQTT.Client for testing
type mockMQTTClient struct {
	connected bool
}

func (m *mockMQTTClient) IsConnected() bool {
	return m.connected
}

func (m *mockMQTTClient) IsConnectionOpen() bool {
	return m.connected
}

func (m *mockMQTTClient) Connect() MQTT.Token {
	return &mockToken{}
}

func (m *mockMQTTClient) Disconnect(quiesce uint) {
	m.connected = false
}

func (m *mockMQTTClient) Publish(topic string, qos byte, retained bool, payload interface{}) MQTT.Token {
	return &mockToken{}
}

func (m *mockMQTTClient) Subscribe(topic string, qos byte, callback MQTT.MessageHandler) MQTT.Token {
	return &mockToken{}
}

func (m *mockMQTTClient) SubscribeMultiple(filters map[string]byte, callback MQTT.MessageHandler) MQTT.Token {
	return &mockToken{}
}

func (m *mockMQTTClient) Unsubscribe(topics ...string) MQTT.Token {
	return &mockToken{}
}

func (m *mockMQTTClient) AddRoute(topic string, callback MQTT.MessageHandler) {
}

func (m *mockMQTTClient) DeleteRoute(topic string) {
}

func (m *mockMQTTClient) OptionsReader() MQTT.ClientOptionsReader {
	return MQTT.ClientOptionsReader{}
}

// mockToken is a mock implementation of MQTT.Token for testing
type mockToken struct{}

func (m *mockToken) Wait() bool {
	return true
}

func (m *mockToken) WaitTimeout(time.Duration) bool {
	return true
}

func (m *mockToken) Done() <-chan struct{} {
	ch := make(chan struct{})
	close(ch)

	return ch
}

func (m *mockToken) Error() error {
	return nil
}

func TestMQTTCollector_run_ContextCancellation(t *testing.T) {
	cfg := &config.Config{
		MQTT: config.MQTTConfig{
			Broker:   "localhost:1883",
			ClientID: "test-client",
		},
	}
	// Use a properly initialized metrics registry
	metricsRegistry := metrics.NewRegistry()
	collector := NewMQTTCollector(cfg, metricsRegistry)

	// Create a context that will be cancelled
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Start the collector in a goroutine
	done := make(chan bool)

	go func() {
		collector.run(ctx)

		done <- true
	}()

	// Wait for either the goroutine to finish or timeout
	select {
	case <-done:
		// Context was cancelled and goroutine finished
	case <-time.After(200 * time.Millisecond):
		t.Error("Expected goroutine to finish when context is cancelled")
	}
}
