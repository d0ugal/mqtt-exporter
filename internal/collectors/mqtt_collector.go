package collectors

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
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
	// Track previous values for cumulative $SYS counters to calculate deltas
	sysCounters map[string]float64
}

func NewMQTTCollector(cfg *config.Config, metricsRegistry *metrics.MQTTRegistry, app *app.App) *MQTTCollector {
	return &MQTTCollector{
		config:         cfg,
		metrics:        metricsRegistry,
		app:            app,
		topics:         make(map[string]int64),
		done:           make(chan struct{}),
		connectionLost: make(chan struct{}, 1),
		sysCounters:    make(map[string]float64),
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
		spanCtx context.Context //nolint:contextcheck // Extracting context from span for child operations
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
		spanCtx context.Context //nolint:contextcheck // Extracting context from span for child operations
	)

	// Build list of topics to subscribe to
	topics := make([]string, len(mc.config.MQTT.Topics))
	copy(topics, mc.config.MQTT.Topics)

	// Add $SYS/# topic unless disabled
	if !mc.config.MQTT.DisableSysTopics {
		topics = append(topics, "$SYS/#")

		slog.Info("$SYS topic monitoring enabled", "topic", "$SYS/#")
	} else {
		slog.Info("$SYS topic monitoring disabled")
	}

	if tracer != nil && tracer.IsEnabled() {
		span = tracer.NewCollectorSpan(ctx, "mqtt-collector", "subscribe-to-topics")

		span.SetAttributes(
			attribute.Int("mqtt.topics_count", len(topics)),
			attribute.Int("mqtt.qos", int(mc.config.MQTT.QoS)),
			attribute.Bool("mqtt.sys_topics_enabled", !mc.config.MQTT.DisableSysTopics),
		)

		spanCtx = span.Context()
		defer span.End()
	} else {
		spanCtx = ctx
	}

	var topicsSubscribed, topicsFailed int

	for _, topic := range topics {
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

	if span != nil {
		span.SetAttributes(
			attribute.Int("subscribe.topics_subscribed", topicsSubscribed),
			attribute.Int("subscribe.topics_failed", topicsFailed),
		)
		span.AddEvent("all_subscriptions_completed",
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

	// Check if this is a $SYS topic
	if len(topic) >= 5 && topic[:5] == "$SYS/" {
		mc.processSysMessage(metricsCtx, topic, payload)
	} else {
		mc.updateMetrics(metricsCtx, topic, payload)
	}

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

// addCounterDelta calculates the delta from the previous value and adds it to the counter.
// This handles cumulative values from the broker's $SYS topics correctly.
// Returns true if the counter was updated, false if this is the first value (no previous baseline).
func (mc *MQTTCollector) addCounterDelta(counterKey string, newValue float64, counterVec *prometheus.CounterVec, labels prometheus.Labels) bool {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	previousValue, exists := mc.sysCounters[counterKey]
	mc.sysCounters[counterKey] = newValue

	if !exists {
		// First time seeing this counter, store the baseline but don't update the metric
		// This prevents a large spike on startup
		return false
	}

	// Calculate delta
	delta := newValue - previousValue

	// Only add positive deltas (handle counter resets gracefully)
	if delta > 0 {
		counterVec.With(labels).Add(delta)
		return true
	} else if delta < 0 {
		// Counter reset detected (broker restart), reset our baseline
		slog.Debug("Counter reset detected",
			"counter", counterKey,
			"previous", previousValue,
			"new", newValue,
		)

		return false
	}

	// delta == 0, no change
	return false
}

// processLoadAverageMetrics handles load average metrics with time intervals
func (mc *MQTTCollector) processLoadAverageMetrics(topic, broker string, value float64) bool {
	// Extract interval from topic (1min, 5min, 15min)
	var interval string
	if strings.HasSuffix(topic, "/1min") {
		interval = "1min"
	} else if strings.HasSuffix(topic, "/5min") {
		interval = "5min"
	} else if strings.HasSuffix(topic, "/15min") {
		interval = "15min"
	}

	if interval == "" {
		return false
	}

	loadLabels := prometheus.Labels{"broker": broker, "interval": interval}

	switch {
	case strings.Contains(topic, "/load/connections/"):
		mc.metrics.MQTTSysBrokerLoadConnections.With(loadLabels).Set(value)
		return true
	case strings.Contains(topic, "/load/bytes/received/"):
		mc.metrics.MQTTSysBrokerLoadBytesReceived.With(loadLabels).Set(value)
		return true
	case strings.Contains(topic, "/load/bytes/sent/"):
		mc.metrics.MQTTSysBrokerLoadBytesSent.With(loadLabels).Set(value)
		return true
	case strings.Contains(topic, "/load/messages/received/"):
		mc.metrics.MQTTSysBrokerLoadMessagesReceived.With(loadLabels).Set(value)
		return true
	case strings.Contains(topic, "/load/messages/sent/"):
		mc.metrics.MQTTSysBrokerLoadMessagesSent.With(loadLabels).Set(value)
		return true
	case strings.Contains(topic, "/load/publish/received/"):
		mc.metrics.MQTTSysBrokerLoadPublishReceived.With(loadLabels).Set(value)
		return true
	case strings.Contains(topic, "/load/publish/sent/"):
		mc.metrics.MQTTSysBrokerLoadPublishSent.With(loadLabels).Set(value)
		return true
	case strings.Contains(topic, "/load/publish/dropped/"):
		mc.metrics.MQTTSysBrokerLoadPublishDropped.With(loadLabels).Set(value)
		return true
	case strings.Contains(topic, "/load/sockets/"):
		mc.metrics.MQTTSysBrokerLoadSockets.With(loadLabels).Set(value)
		return true
	}

	return false
}

// processClientMetrics handles client connection metrics
func (mc *MQTTCollector) processClientMetrics(topic string, labels prometheus.Labels, value float64) bool {
	broker := labels["broker"]

	switch {
	case strings.HasSuffix(topic, "/broker/clients/connected") || strings.HasSuffix(topic, "/clients/connected") || strings.HasSuffix(topic, "/clients/active"):
		mc.metrics.MQTTSysBrokerClientsConnected.With(labels).Set(value)
		return true
	case strings.HasSuffix(topic, "/broker/clients/disconnected") || strings.HasSuffix(topic, "/clients/disconnected") || strings.HasSuffix(topic, "/clients/inactive"):
		mc.metrics.MQTTSysBrokerClientsDisconnected.With(labels).Set(value)
		return true
	case strings.HasSuffix(topic, "/broker/clients/expired") || strings.HasSuffix(topic, "/clients/expired"):
		counterKey := broker + ":clients_expired"
		mc.addCounterDelta(counterKey, value, mc.metrics.MQTTSysBrokerClientsExpired, labels)

		return true
	case strings.HasSuffix(topic, "/broker/clients/total") || strings.HasSuffix(topic, "/clients/total"):
		mc.metrics.MQTTSysBrokerClientsTotal.With(labels).Set(value)
		return true
	case strings.HasSuffix(topic, "/broker/clients/maximum") || strings.HasSuffix(topic, "/clients/maximum"):
		mc.metrics.MQTTSysBrokerClientsMaximum.With(labels).Set(value)
		return true
	}

	return false
}

// processMessageMetrics handles message-related metrics
func (mc *MQTTCollector) processMessageMetrics(topic string, labels prometheus.Labels, value float64) bool {
	broker := labels["broker"]

	switch {
	case strings.HasSuffix(topic, "/broker/messages/received") || strings.HasSuffix(topic, "/messages/received"):
		counterKey := broker + ":messages_received"
		mc.addCounterDelta(counterKey, value, mc.metrics.MQTTSysBrokerMessagesReceived, labels)

		return true
	case strings.HasSuffix(topic, "/broker/messages/sent") || strings.HasSuffix(topic, "/messages/sent"):
		counterKey := broker + ":messages_sent"
		mc.addCounterDelta(counterKey, value, mc.metrics.MQTTSysBrokerMessagesSent, labels)

		return true
	case strings.HasSuffix(topic, "/broker/messages/inflight") || strings.HasSuffix(topic, "/messages/inflight"):
		mc.metrics.MQTTSysBrokerMessagesInflight.With(labels).Set(value)
		return true
	case strings.HasSuffix(topic, "/broker/messages/stored") || strings.HasSuffix(topic, "/messages/stored"):
		mc.metrics.MQTTSysBrokerStoreMessagesCount.With(labels).Set(value)
		return true
	case strings.HasSuffix(topic, "/broker/store/messages/count") || strings.HasSuffix(topic, "/store/messages/count"):
		mc.metrics.MQTTSysBrokerStoreMessagesCount.With(labels).Set(value)
		return true
	case strings.HasSuffix(topic, "/broker/store/messages/bytes") || strings.HasSuffix(topic, "/store/messages/bytes"):
		mc.metrics.MQTTSysBrokerStoreMessagesBytes.With(labels).Set(value)
		return true
	}

	return false
}

// processByteMetrics handles byte transfer metrics
func (mc *MQTTCollector) processByteMetrics(topic string, labels prometheus.Labels, value float64) bool {
	broker := labels["broker"]

	switch {
	case strings.HasSuffix(topic, "/broker/bytes/received") || strings.HasSuffix(topic, "/bytes/received"):
		counterKey := broker + ":bytes_received"
		mc.addCounterDelta(counterKey, value, mc.metrics.MQTTSysBrokerBytesReceived, labels)

		return true
	case strings.HasSuffix(topic, "/broker/bytes/sent") || strings.HasSuffix(topic, "/bytes/sent"):
		counterKey := broker + ":bytes_sent"
		mc.addCounterDelta(counterKey, value, mc.metrics.MQTTSysBrokerBytesSent, labels)

		return true
	}

	return false
}

// processPublishMetrics handles publish-related metrics
func (mc *MQTTCollector) processPublishMetrics(topic string, labels prometheus.Labels, value float64) bool {
	broker := labels["broker"]

	switch {
	case strings.HasSuffix(topic, "/broker/publish/messages/dropped") || strings.HasSuffix(topic, "/publish/messages/dropped"):
		counterKey := broker + ":publish_dropped"
		mc.addCounterDelta(counterKey, value, mc.metrics.MQTTSysBrokerPublishDropped, labels)

		return true
	case strings.HasSuffix(topic, "/broker/publish/messages/received") || strings.HasSuffix(topic, "/publish/messages/received"):
		counterKey := broker + ":publish_received"
		mc.addCounterDelta(counterKey, value, mc.metrics.MQTTSysBrokerPublishReceived, labels)

		return true
	case strings.HasSuffix(topic, "/broker/publish/messages/sent") || strings.HasSuffix(topic, "/publish/messages/sent"):
		counterKey := broker + ":publish_sent"
		mc.addCounterDelta(counterKey, value, mc.metrics.MQTTSysBrokerPublishSent, labels)

		return true
	}

	return false
}

// processSubscriptionMetrics handles subscription and retained message metrics
func (mc *MQTTCollector) processSubscriptionMetrics(topic string, labels prometheus.Labels, value float64) bool {
	switch {
	case strings.HasSuffix(topic, "/broker/subscriptions/count") || strings.HasSuffix(topic, "/subscriptions/count"):
		mc.metrics.MQTTSysBrokerSubscriptionsCount.With(labels).Set(value)
		return true
	case strings.HasSuffix(topic, "/broker/retained/messages/count") || strings.HasSuffix(topic, "/retained messages/count"):
		mc.metrics.MQTTSysBrokerRetainedMessagesCount.With(labels).Set(value)
		return true
	}

	return false
}

// processHeapMetrics handles heap memory metrics
func (mc *MQTTCollector) processHeapMetrics(topic string, labels prometheus.Labels, value float64) bool {
	switch {
	case strings.HasSuffix(topic, "/broker/heap/current") || strings.HasSuffix(topic, "/heap/current size"):
		mc.metrics.MQTTSysBrokerHeapCurrentBytes.With(labels).Set(value)
		return true
	case strings.HasSuffix(topic, "/broker/heap/maximum") || strings.HasSuffix(topic, "/heap/maximum size"):
		mc.metrics.MQTTSysBrokerHeapMaximumBytes.With(labels).Set(value)
		return true
	}

	return false
}

// processSysMessage processes $SYS topic messages and maps them to specific metrics
func (mc *MQTTCollector) processSysMessage(ctx context.Context, topic string, payload []byte) {
	tracer := mc.app.GetTracer()

	var span *tracing.CollectorSpan

	if tracer != nil && tracer.IsEnabled() {
		span = tracer.NewCollectorSpan(ctx, "mqtt-collector", "process-sys-message")

		span.SetAttributes(
			attribute.String("mqtt.topic", topic),
			attribute.Int("mqtt.payload_length", len(payload)),
		)

		defer span.End()
	}

	updateStart := time.Now()

	// Try to parse the payload as a numeric value
	payloadStr := string(payload)

	value, err := parseNumericValue(payloadStr)
	if err != nil {
		// Not a numeric value, log and skip
		slog.Debug("Skipping non-numeric $SYS topic",
			"topic", topic,
			"payload", payloadStr,
			"error", err,
		)

		if span != nil {
			span.SetAttributes(
				attribute.Bool("sys.is_numeric", false),
				attribute.String("sys.payload", payloadStr),
			)
		}

		return
	}

	broker := mc.config.MQTT.Broker
	labels := prometheus.Labels{"broker": broker}

	// Map $SYS topics to specific metrics
	// This supports common Mosquito topic patterns
	updated := false

	// Handle version metric specially (it's a string, not numeric at this point)
	if strings.HasSuffix(topic, "/broker/version") {
		// Version is a string, set gauge to 1 with version as label
		versionLabels := prometheus.Labels{"broker": broker, "version": payloadStr}
		mc.metrics.MQTTSysBrokerVersion.With(versionLabels).Set(1)

		updated = true
	}

	// Handle load average metrics with time intervals
	if !updated && strings.Contains(topic, "/broker/load/") {
		updated = mc.processLoadAverageMetrics(topic, broker, value)
	}

	// Handle regular metrics by trying each category
	if !updated {
		updated = mc.processClientMetrics(topic, labels, value)
	}

	if !updated {
		updated = mc.processMessageMetrics(topic, labels, value)
	}

	if !updated {
		updated = mc.processByteMetrics(topic, labels, value)
	}

	if !updated {
		updated = mc.processPublishMetrics(topic, labels, value)
	}

	if !updated {
		updated = mc.processSubscriptionMetrics(topic, labels, value)
	}

	if !updated {
		updated = mc.processHeapMetrics(topic, labels, value)
	}

	if updated {
		slog.Debug("Updated $SYS metric",
			"topic", topic,
			"value", value,
		)

		if span != nil {
			span.SetAttributes(
				attribute.Bool("sys.is_numeric", true),
				attribute.Bool("sys.metric_mapped", true),
				attribute.Float64("sys.value", value),
				attribute.Float64("metrics.update_duration_seconds", time.Since(updateStart).Seconds()),
			)
			span.AddEvent("sys_metric_updated",
				attribute.String("topic", topic),
				attribute.Float64("value", value),
			)
		}
	} else {
		slog.Debug("No metric mapping found for $SYS topic",
			"topic", topic,
			"value", value,
		)

		if span != nil {
			span.SetAttributes(
				attribute.Bool("sys.is_numeric", true),
				attribute.Bool("sys.metric_mapped", false),
				attribute.String("sys.topic", topic),
			)
		}
	}
}

// parseNumericValue attempts to parse a string as a numeric value
func parseNumericValue(s string) (float64, error) {
	// Try to parse as float64
	value, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse as numeric: %w", err)
	}

	return value, nil
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
