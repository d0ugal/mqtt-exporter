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
	tracer := mc.app.GetTracer()

	var (
		span    *tracing.CollectorSpan
		spanCtx context.Context
	)

	if tracer != nil && tracer.IsEnabled() {
		span = tracer.NewCollectorSpan(ctx, "mqtt-collector", "connect")

		span.SetAttributes(
			attribute.String("mqtt.broker", mc.config.MQTT.Broker),
			attribute.String("mqtt.client_id", mc.config.MQTT.ClientID),
			attribute.Bool("mqtt.clean_session", mc.config.MQTT.CleanSession),
			attribute.Int64("mqtt.keep_alive_seconds", int64(mc.config.MQTT.KeepAlive.Duration.Seconds())),
			attribute.Int64("mqtt.connect_timeout_seconds", int64(mc.config.MQTT.ConnectTimeout.Duration.Seconds())),
			attribute.Bool("mqtt.has_username", mc.config.MQTT.Username != ""),
		)

		spanCtx = span.Context()
		defer span.End()
	} else {
		spanCtx = ctx
	}

	configStart := time.Now()

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

	configDuration := time.Since(configStart)

	if span != nil {
		span.SetAttributes(
			attribute.Float64("config.duration_seconds", configDuration.Seconds()),
		)
		span.AddEvent("config_created")
	}

	// Set up connection handlers
	opts.SetOnConnectHandler(mc.onConnect)
	opts.SetConnectionLostHandler(mc.onConnectionLost)
	opts.SetDefaultPublishHandler(mc.onMessageReceived)

	mc.client = MQTT.NewClient(opts)

	// Connect to broker
	connectStart := time.Now()

	// Create context with timeout using span context if available
	timeoutCtx, cancel := context.WithTimeout(spanCtx, mc.config.MQTT.ConnectTimeout.Duration)
	defer cancel()

	if token := mc.client.Connect(); token.Wait() && token.Error() != nil {
		connectDuration := time.Since(connectStart)

		if span != nil {
			span.SetAttributes(
				attribute.Float64("connect.duration_seconds", connectDuration.Seconds()),
				attribute.Bool("connect.success", false),
			)
			span.RecordError(token.Error(), attribute.String("operation", "mqtt_connect"))
		}

		_ = timeoutCtx // Use timeoutCtx for future operations if needed

		return fmt.Errorf("failed to connect: %w", token.Error())
	}

	connectDuration := time.Since(connectStart)

	if span != nil {
		span.SetAttributes(
			attribute.Float64("connect.duration_seconds", connectDuration.Seconds()),
			attribute.Bool("connect.success", true),
		)
		span.AddEvent("connection_established",
			attribute.String("broker", mc.config.MQTT.Broker),
		)
	}

	_ = timeoutCtx // Use timeoutCtx for future operations if needed

	return nil
}

func (mc *MQTTCollector) subscribeToTopics(ctx context.Context) error {
	tracer := mc.app.GetTracer()

	var (
		span    *tracing.CollectorSpan
		spanCtx context.Context
	)

	if tracer != nil && tracer.IsEnabled() {
		span = tracer.NewCollectorSpan(ctx, "mqtt-collector", "subscribe-to-topics")

		span.SetAttributes(
			attribute.Int("mqtt.topics_count", len(mc.config.MQTT.Topics)),
			attribute.Int("mqtt.qos", int(mc.config.MQTT.QoS)),
		)

		spanCtx = span.Context()
		defer span.End()
	} else {
		spanCtx = ctx
	}

	var topicsSubscribed, topicsFailed int

	for _, topic := range mc.config.MQTT.Topics {
		subscribeStart := time.Now()

		var topicSpan *tracing.CollectorSpan

		if tracer != nil && tracer.IsEnabled() {
			topicSpan = tracer.NewCollectorSpan(spanCtx, "mqtt-collector", "subscribe-topic")

			topicSpan.SetAttributes(
				attribute.String("mqtt.topic", topic),
				attribute.Int("mqtt.qos", int(mc.config.MQTT.QoS)),
			)

			defer topicSpan.End()
		}

		if token := mc.client.Subscribe(topic, byte(mc.config.MQTT.QoS), nil); token.Wait() && token.Error() != nil {
			subscribeDuration := time.Since(subscribeStart)
			topicsFailed++

			if topicSpan != nil {
				topicSpan.SetAttributes(
					attribute.Float64("subscribe.duration_seconds", subscribeDuration.Seconds()),
					attribute.Bool("subscribe.success", false),
				)
				topicSpan.RecordError(token.Error(), attribute.String("operation", "mqtt_subscribe"))
			}

			return fmt.Errorf("failed to subscribe to topic %s: %w", topic, token.Error())
		}

		subscribeDuration := time.Since(subscribeStart)
		topicsSubscribed++

		if topicSpan != nil {
			topicSpan.SetAttributes(
				attribute.Float64("subscribe.duration_seconds", subscribeDuration.Seconds()),
				attribute.Bool("subscribe.success", true),
			)
			topicSpan.AddEvent("subscription_completed",
				attribute.String("topic", topic),
			)
		}

		slog.Info("Subscribed to topic", "topic", topic)
	}

	if parentSpan != nil {
		parentSpan.SetAttributes(
			attribute.Int("subscribe.topics_subscribed", topicsSubscribed),
			attribute.Int("subscribe.topics_failed", topicsFailed),
		)
		parentSpan.AddEvent("all_subscriptions_completed",
			attribute.Int("topics_subscribed", topicsSubscribed),
		)
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

	// Update topic counter with tracing
	updateCounterStart := time.Now()

	mc.mu.Lock()
	mc.topics[topic]++
	messageCount := mc.topics[topic]
	mc.mu.Unlock()

	updateCounterDuration := time.Since(updateCounterStart)

	if messageSpan != nil {
		messageSpan.SetAttributes(
			attribute.Float64("update.counter_duration_seconds", updateCounterDuration.Seconds()),
			attribute.Int64("update.message_count", messageCount),
		)
		messageSpan.AddEvent("counter_updated")
	}

	// Update metrics with tracing
	updateMetricsStart := time.Now()

	var metricsCtx context.Context

	if messageSpan != nil {
		metricsCtx = messageSpan.Context()
	} else {
		metricsCtx = context.Background()
	}

	mc.updateMetrics(metricsCtx, topic, payload)

	updateMetricsDuration := time.Since(updateMetricsStart)

	slog.Debug("Updated metrics for topic", "topic", topic)

	// Add event to span
	if messageSpan != nil {
		messageSpan.SetAttributes(
			attribute.Float64("update.metrics_duration_seconds", updateMetricsDuration.Seconds()),
		)
		messageSpan.AddEvent("metrics_updated",
			attribute.String("topic", topic),
			attribute.Int("message_count", int(messageCount)),
		)
	}
}

// updateMetrics updates Prometheus metrics with tracing
func (mc *MQTTCollector) updateMetrics(ctx context.Context, topic string, payload []byte) {
	tracer := mc.app.GetTracer()

	var span *tracing.CollectorSpan

	if tracer != nil && tracer.IsEnabled() {
		span = tracer.NewCollectorSpan(ctx, "mqtt-collector", "update-metrics")

		span.SetAttributes(
			attribute.String("mqtt.topic", topic),
			attribute.Int("mqtt.payload_length", len(payload)),
		)

		defer span.End()
	}

	updateStart := time.Now()

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

	if span != nil {
		span.SetAttributes(
			attribute.Float64("metrics.update_duration_seconds", time.Since(updateStart).Seconds()),
			attribute.Int("metrics.count", 3),
		)
		span.AddEvent("metrics_updated",
			attribute.String("topic", topic),
			attribute.Int("payload_length", len(payload)),
		)
	}

	_ = spanCtx // Use spanCtx for future operations if needed
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
