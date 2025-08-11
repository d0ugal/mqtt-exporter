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
}

func NewMQTTCollector(cfg *config.Config, metricsRegistry *metrics.Registry) *MQTTCollector {
	return &MQTTCollector{
		config:  cfg,
		metrics: metricsRegistry,
		topics:  make(map[string]int64),
	}
}

func (mc *MQTTCollector) Start(ctx context.Context) {
	go mc.connect(ctx)
}

func (mc *MQTTCollector) connect(ctx context.Context) {
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
		slog.Error("Failed to connect to MQTT broker",
			"broker", mc.config.MQTT.Broker,
			"error", token.Error(),
		)
		mc.metrics.MQTTConnectionErrors.WithLabelValues(mc.config.MQTT.Broker, "connect").Inc()

		return
	}

	slog.Info("Connected to MQTT broker", "broker", mc.config.MQTT.Broker)
	mc.metrics.MQTTConnectionStatus.WithLabelValues(mc.config.MQTT.Broker).Set(1)

	// Subscribe to topics
	for _, topic := range mc.config.MQTT.Topics {
		if token := mc.client.Subscribe(topic, byte(mc.config.MQTT.QoS), nil); token.Wait() && token.Error() != nil {
			slog.Error("Failed to subscribe to topic",
				"topic", topic,
				"error", token.Error(),
			)
			mc.metrics.MQTTConnectionErrors.WithLabelValues(mc.config.MQTT.Broker, "subscribe").Inc()
		} else {
			slog.Info("Subscribed to topic", "topic", topic)
		}
	}

	// Keep connection alive
	<-ctx.Done()
	slog.Info("Shutting down MQTT collector")

	if mc.client != nil && mc.client.IsConnected() {
		mc.client.Disconnect(250)
	}
}

func (mc *MQTTCollector) onConnect(client MQTT.Client) {
	slog.Info("MQTT connection established")
	mc.metrics.MQTTConnectionStatus.WithLabelValues(mc.config.MQTT.Broker).Set(1)
}

func (mc *MQTTCollector) onConnectionLost(client MQTT.Client, err error) {
	slog.Error("MQTT connection lost", "error", err)
	mc.metrics.MQTTConnectionStatus.WithLabelValues(mc.config.MQTT.Broker).Set(0)
	mc.metrics.MQTTConnectionErrors.WithLabelValues(mc.config.MQTT.Broker, "connection_lost").Inc()
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
