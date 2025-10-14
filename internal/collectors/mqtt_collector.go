package collectors

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/d0ugal/mqtt-exporter/internal/config"
	"github.com/d0ugal/mqtt-exporter/internal/metrics"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type MQTTCollector struct {
	config  *config.Config
	metrics *metrics.Registry
	client  MQTT.Client
	mu      sync.RWMutex
	topics  map[string]int64
	done    chan struct{}
}

func NewMQTTCollector(cfg *config.Config, metricsRegistry *metrics.Registry) *MQTTCollector {
	return &MQTTCollector{
		config:  cfg,
		metrics: metricsRegistry,
		topics:  make(map[string]int64),
		done:    make(chan struct{}),
	}
}

func (mc *MQTTCollector) Start(ctx context.Context) {
	go mc.run(ctx)
}

// run handles the main connection loop with automatic reconnection
func (mc *MQTTCollector) run(ctx context.Context) {
	reconnectDelay := time.Second
	maxReconnectDelay := time.Minute

	for {
		select {
		case <-ctx.Done():
			return
		case <-mc.done:
			return
		default:
		}

		if err := mc.connect(); err != nil {
			slog.Error("Failed to connect to MQTT broker",
				"broker", mc.config.MQTT.Broker,
				"error", err,
			)
			mc.metrics.MQTTConnectionStatus.WithLabelValues(mc.config.MQTT.Broker).Set(0)
			mc.metrics.MQTTConnectionErrors.WithLabelValues(mc.config.MQTT.Broker, "connect").Inc()

			select {
			case <-ctx.Done():
				return
			case <-time.After(reconnectDelay):
				reconnectDelay = minDuration(reconnectDelay*2, maxReconnectDelay)
				continue
			}
		}

		// Reset reconnect delay on successful connection
		reconnectDelay = time.Second

		slog.Info("Connected to MQTT broker", "broker", mc.config.MQTT.Broker)
		mc.metrics.MQTTConnectionStatus.WithLabelValues(mc.config.MQTT.Broker).Set(1)

		// Subscribe to topics
		if err := mc.subscribeToTopics(); err != nil {
			slog.Error("Failed to subscribe to topics", "error", err)
			mc.metrics.MQTTConnectionErrors.WithLabelValues(mc.config.MQTT.Broker, "subscribe").Inc()

			// Disconnect and retry
			if mc.client != nil {
				mc.client.Disconnect(250)
			}

			continue
		}

		// Monitor connection status and reconnect if lost
		connectionLost := mc.monitorConnection(ctx)
		if connectionLost {
			slog.Info("Connection lost, attempting to reconnect", "broker", mc.config.MQTT.Broker)
			// Disconnect current client before reconnecting
			if mc.client != nil {
				mc.client.Disconnect(250)
			}

			continue
		}

		// Check if context was cancelled
		select {
		case <-ctx.Done():
			slog.Info("Context cancelled, shutting down MQTT collector")

			if mc.client != nil && mc.client.IsConnected() {
				mc.client.Disconnect(250)
			}

			return
		default:
			// This shouldn't happen in normal operation, but handle gracefully
			slog.Warn("monitorConnection returned false but context not cancelled, continuing...")
			continue
		}
	}
}

func (mc *MQTTCollector) connect() error {
	opts := MQTT.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s", mc.config.MQTT.Broker))
	opts.SetClientID(mc.config.MQTT.ClientID)

	if mc.config.MQTT.Username != "" {
		opts.SetUsername(mc.config.MQTT.Username)
		opts.SetPassword(mc.config.MQTT.Password)
	}

	opts.SetCleanSession(mc.config.MQTT.CleanSession)
	opts.SetKeepAlive(time.Duration(mc.config.MQTT.KeepAlive) * time.Second)
	opts.SetConnectTimeout(time.Duration(mc.config.MQTT.ConnectTimeout) * time.Second)

	// Set up connection handlers
	opts.SetOnConnectHandler(mc.onConnect)
	opts.SetConnectionLostHandler(mc.onConnectionLost)
	opts.SetDefaultPublishHandler(mc.onMessageReceived)

	mc.client = MQTT.NewClient(opts)

	// Connect to broker
	if token := mc.client.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to connect: %w", token.Error())
	}

	return nil
}

func (mc *MQTTCollector) subscribeToTopics() error {
	for _, topic := range mc.config.MQTT.Topics {
		if token := mc.client.Subscribe(topic, byte(mc.config.MQTT.QoS), nil); token.Wait() && token.Error() != nil {
			return fmt.Errorf("failed to subscribe to topic %s: %w", topic, token.Error())
		}

		slog.Info("Subscribed to topic", "topic", topic)
	}

	return nil
}

func (mc *MQTTCollector) onConnect(client MQTT.Client) {
	slog.Info("MQTT connection established", "broker", mc.config.MQTT.Broker)
	mc.metrics.MQTTConnectionStatus.WithLabelValues(mc.config.MQTT.Broker).Set(1)
}

func (mc *MQTTCollector) onConnectionLost(client MQTT.Client, err error) {
	slog.Error("MQTT connection lost",
		"broker", mc.config.MQTT.Broker,
		"error", err,
	)
	mc.metrics.MQTTConnectionStatus.WithLabelValues(mc.config.MQTT.Broker).Set(0)
	mc.metrics.MQTTConnectionErrors.WithLabelValues(mc.config.MQTT.Broker, "connection_lost").Inc()
	mc.metrics.MQTTReconnectsTotal.WithLabelValues(mc.config.MQTT.Broker).Inc()

	slog.Info("MQTT reconnection attempt initiated", "broker", mc.config.MQTT.Broker)
}

// monitorConnection monitors the MQTT connection and returns true if connection is lost
func (mc *MQTTCollector) monitorConnection(ctx context.Context) bool {
	// Use a shorter interval for more responsive connection monitoring
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return false // Context cancelled, not a connection loss
		case <-mc.done:
			return false // Shutdown requested, not a connection loss
		case <-ticker.C:
			// Check if client is still connected
			if mc.client == nil || !mc.client.IsConnected() {
				slog.Debug("Connection check failed - client disconnected", "broker", mc.config.MQTT.Broker)
				return true // Connection lost
			}
		}
	}
}

func (mc *MQTTCollector) onMessageReceived(client MQTT.Client, msg MQTT.Message) {
	topic := msg.Topic()
	payload := msg.Payload()

	slog.Debug("Received MQTT message",
		"topic", topic,
		"payload_length", len(payload),
		"qos", msg.Qos(),
	)

	// Update metrics
	mc.mu.Lock()
	mc.topics[topic]++
	mc.mu.Unlock()

	// Increment counters
	mc.metrics.MQTTMessageCount.WithLabelValues(topic).Inc()
	mc.metrics.MQTTMessageBytes.WithLabelValues(topic).Add(float64(len(payload)))
	mc.metrics.MQTTTopicLastMessage.WithLabelValues(topic).Set(float64(time.Now().Unix()))

	slog.Debug("Updated metrics for topic", "topic", topic)
}

// Stop stops the collector
func (mc *MQTTCollector) Stop() {
	close(mc.done)

	if mc.client != nil && mc.client.IsConnected() {
		mc.client.Disconnect(250)
	}
}

// minDuration returns the minimum of two time.Duration values
func minDuration(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}

	return b
}
