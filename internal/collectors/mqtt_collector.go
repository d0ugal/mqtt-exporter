package collectors

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/d0ugal/mqtt-exporter/internal/config"
	"github.com/d0ugal/mqtt-exporter/internal/metrics"
	"github.com/d0ugal/promexporter/app"
	"github.com/d0ugal/promexporter/tracing"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/attribute"
)

type MQTTCollector struct {
	config         *config.Config
	metrics        *metrics.MQTTRegistry
	app            *app.App
	client         MQTT.Client
	mu             sync.RWMutex
	topics         map[string]int64
	done           chan struct{}
	connectionLost chan struct{}
}

func NewMQTTCollector(cfg *config.Config, metricsRegistry *metrics.MQTTRegistry, app *app.App) *MQTTCollector {
	return &MQTTCollector{
		config:         cfg,
		metrics:        metricsRegistry,
		app:            app,
		topics:         make(map[string]int64),
		done:           make(chan struct{}),
		connectionLost: make(chan struct{}, 1),
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
		// Create a span for each connection attempt
		tracer := mc.app.GetTracer()

		var collectorSpan *tracing.CollectorSpan

		spanCtx := context.Background()

		if tracer != nil && tracer.IsEnabled() {
			collectorSpan = tracer.NewCollectorSpan(ctx, "mqtt-collector", "connection-attempt")
			spanCtx = collectorSpan.Context()
		}

		select {
		case <-spanCtx.Done():
			slog.Info("Shutting down MQTT collector")

			if collectorSpan != nil {
				collectorSpan.AddEvent("shutdown_requested")
				collectorSpan.End()
			}

			if mc.client != nil && mc.client.IsConnected() {
				mc.client.Disconnect(250)
			}

			return
		case <-mc.done:
			slog.Info("Stopping MQTT collector")

			if collectorSpan != nil {
				collectorSpan.AddEvent("stop_requested")
				collectorSpan.End()
			}

			if mc.client != nil && mc.client.IsConnected() {
				mc.client.Disconnect(250)
			}

			return
		default:
		}

		if err := mc.connect(spanCtx); err != nil { //nolint:contextcheck
			slog.Error("Failed to connect to MQTT broker",
				"broker", mc.config.MQTT.Broker,
				"error", err,
			)

			if collectorSpan != nil {
				collectorSpan.RecordError(err, attribute.String("broker", mc.config.MQTT.Broker))
				collectorSpan.End()
			}

			mc.metrics.MQTTConnectionStatus.With(prometheus.Labels{
				"broker": mc.config.MQTT.Broker,
			}).Set(0)
			mc.metrics.MQTTConnectionErrors.With(prometheus.Labels{
				"broker": mc.config.MQTT.Broker,
				"reason": "connect",
			}).Inc()

			select {
			case <-spanCtx.Done():
				if collectorSpan != nil {
					collectorSpan.End()
				}

				return
			case <-time.After(reconnectDelay):
				reconnectDelay = minDuration(reconnectDelay*2, maxReconnectDelay)
				if collectorSpan != nil {
					collectorSpan.AddEvent("reconnect_scheduled",
						attribute.String("delay", reconnectDelay.String()))
					collectorSpan.End()
				}

				continue
			}
		}

		// Reset reconnect delay on successful connection
		reconnectDelay = time.Second

		slog.Info("Connected to MQTT broker", "broker", mc.config.MQTT.Broker)

		if collectorSpan != nil {
			collectorSpan.AddEvent("connected", attribute.String("broker", mc.config.MQTT.Broker))
		}

		mc.metrics.MQTTConnectionStatus.With(prometheus.Labels{
			"broker": mc.config.MQTT.Broker,
		}).Set(1)

		// Subscribe to topics
		if err := mc.subscribeToTopics(spanCtx); err != nil { //nolint:contextcheck
			slog.Error("Failed to subscribe to topics", "error", err)

			if collectorSpan != nil {
				collectorSpan.RecordError(err, attribute.String("operation", "subscribe"))
				collectorSpan.End()
			}

			mc.metrics.MQTTConnectionErrors.With(prometheus.Labels{
				"broker": mc.config.MQTT.Broker,
				"reason": "subscribe",
			}).Inc()

			// Disconnect and retry
			if mc.client != nil {
				mc.client.Disconnect(250)
			}

			continue
		}

		// Connection and subscription successful - end the connection span
		if collectorSpan != nil {
			collectorSpan.AddEvent("subscription_completed")
			collectorSpan.End()
		}

		// Wait for connection to be lost or context cancellation
		select {
		case <-spanCtx.Done():
			slog.Info("Shutting down MQTT collector")

			if collectorSpan != nil {
				collectorSpan.AddEvent("shutdown_requested")
				collectorSpan.End()
			}

			if mc.client != nil && mc.client.IsConnected() {
				mc.client.Disconnect(250)
			}

			return
		case <-mc.connectionLost:
			slog.Info("Connection lost, attempting to reconnect", "broker", mc.config.MQTT.Broker)

			if collectorSpan != nil {
				collectorSpan.AddEvent("connection_lost", attribute.String("broker", mc.config.MQTT.Broker))
				collectorSpan.End()
			}
			// Clean up the old client
			if mc.client != nil {
				mc.client.Disconnect(250)
			}
			// Continue the loop to reconnect
		}
	}
}

func (mc *MQTTCollector) connect(ctx context.Context) error {
	opts := MQTT.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s", mc.config.MQTT.Broker))
	opts.SetClientID(mc.config.MQTT.ClientID)

	if mc.config.MQTT.Username != "" {
		opts.SetUsername(mc.config.MQTT.Username)
		opts.SetPassword(mc.config.MQTT.Password)
	}

	opts.SetCleanSession(mc.config.MQTT.CleanSession)
	opts.SetKeepAlive(mc.config.MQTT.KeepAlive.Duration)
	opts.SetConnectTimeout(mc.config.MQTT.ConnectTimeout.Duration)

	// Enhanced connection robustness settings
	opts.SetAutoReconnect(true)
	opts.SetConnectRetry(true)
	opts.SetConnectRetryInterval(5 * time.Second)
	opts.SetMaxReconnectInterval(30 * time.Second)
	opts.SetPingTimeout(10 * time.Second)
	opts.SetWriteTimeout(10 * time.Second)
	opts.SetResumeSubs(true) // Resume subscriptions after reconnection

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

func (mc *MQTTCollector) subscribeToTopics(ctx context.Context) error {
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
	mc.metrics.MQTTConnectionStatus.With(prometheus.Labels{
		"broker": mc.config.MQTT.Broker,
	}).Set(1)
}

func (mc *MQTTCollector) onConnectionLost(client MQTT.Client, err error) {
	slog.Error("MQTT connection lost",
		"broker", mc.config.MQTT.Broker,
		"error", err,
	)
	mc.metrics.MQTTConnectionStatus.With(prometheus.Labels{
		"broker": mc.config.MQTT.Broker,
	}).Set(0)
	mc.metrics.MQTTConnectionErrors.With(prometheus.Labels{
		"broker": mc.config.MQTT.Broker,
		"reason": "connection_lost",
	}).Inc()
	mc.metrics.MQTTReconnectsTotal.With(prometheus.Labels{
		"broker": mc.config.MQTT.Broker,
	}).Inc()

	// Signal that connection was lost to trigger reconnection
	select {
	case mc.connectionLost <- struct{}{}:
	default:
		// Channel is full, connection lost signal already pending
	}

	slog.Info("MQTT reconnection attempt initiated", "broker", mc.config.MQTT.Broker)
}

func (mc *MQTTCollector) onMessageReceived(client MQTT.Client, msg MQTT.Message) {
	topic := msg.Topic()
	payload := msg.Payload()

	slog.Debug("Received MQTT message",
		"topic", topic,
		"payload_length", len(payload),
		"qos", msg.Qos(),
	)

	// Create a span for each message processing
	tracer := mc.app.GetTracer()

	var messageSpan *tracing.CollectorSpan

	if tracer != nil && tracer.IsEnabled() {
		messageSpan = tracer.NewCollectorSpan(context.Background(), "mqtt-collector", "process-message")

		// Add message attributes to the span
		messageSpan.SetAttributes(
			attribute.String("mqtt.topic", topic),
			attribute.Int("mqtt.payload_length", len(payload)),
			attribute.Int("mqtt.qos", int(msg.Qos())),
		)

		defer messageSpan.End()
	}

	// Update metrics
	mc.mu.Lock()
	mc.topics[topic]++
	mc.mu.Unlock()

	// Increment counters
	mc.metrics.MQTTMessageCount.With(prometheus.Labels{
		"topic": topic,
	}).Inc()
	mc.metrics.MQTTMessageBytes.With(prometheus.Labels{
		"topic": topic,
	}).Add(float64(len(payload)))
	mc.metrics.MQTTTopicLastMessage.With(prometheus.Labels{
		"topic": topic,
	}).Set(float64(time.Now().Unix()))

	slog.Debug("Updated metrics for topic", "topic", topic)

	// Add event to span
	if messageSpan != nil {
		messageSpan.AddEvent("metrics_updated",
			attribute.String("topic", topic),
			attribute.Int("message_count", int(mc.topics[topic])),
		)
	}
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
