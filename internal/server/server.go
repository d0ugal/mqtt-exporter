package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/d0ugal/mqtt-exporter/internal/config"
	"github.com/d0ugal/mqtt-exporter/internal/metrics"
	"github.com/d0ugal/mqtt-exporter/internal/version"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	config  *config.Config
	metrics *metrics.Registry
	server  *http.Server
	router  *gin.Engine
}

type MetricInfo struct {
	Name         string
	Help         string
	ExampleValue string
	Labels       map[string]string
}

func New(cfg *config.Config, metricsRegistry *metrics.Registry) *Server {
	router := gin.New()
	router.Use(gin.Recovery())

	server := &Server{
		config:  cfg,
		metrics: metricsRegistry,
		router:  router,
	}

	server.setupRoutes()

	return server
}

func (s *Server) setupRoutes() {
	// Root endpoint with HTML dashboard
	s.router.GET("/", s.handleRoot)

	// Metrics endpoint
	s.router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Health endpoint
	s.router.GET("/health", s.handleHealth)
}

func (s *Server) getMetricsInfo() []MetricInfo {
	metricsInfo := make([]MetricInfo, 0, 7)

	// Define all metrics manually since reflection approach is complex with Prometheus metrics
	metrics := []struct {
		name  string
		field string
	}{
		{"mqtt_exporter_info", "VersionInfo"},
		{"mqtt_messages_total", "MQTTMessageCount"},
		{"mqtt_message_bytes_total", "MQTTMessageBytes"},
		{"mqtt_connection_status", "MQTTConnectionStatus"},
		{"mqtt_connection_errors_total", "MQTTConnectionErrors"},
		{"mqtt_reconnects_total", "MQTTReconnectsTotal"},
		{"mqtt_topic_last_message_timestamp", "MQTTTopicLastMessage"},
	}

	for _, metric := range metrics {
		metricsInfo = append(metricsInfo, MetricInfo{
			Name:         metric.name,
			Help:         s.getMetricHelp(metric.field),
			ExampleValue: s.getExampleValue(metric.field),
			Labels:       s.getExampleLabels(metric.field),
		})
	}

	return metricsInfo
}

func (s *Server) getExampleLabels(metricName string) map[string]string {
	switch metricName {
	case "VersionInfo":
		return map[string]string{"version": "v1.5.0", "commit": "abc123", "build_date": "2024-01-01"}
	case "MQTTMessageCount", "MQTTMessageBytes", "MQTTTopicLastMessage":
		return map[string]string{"topic": "home/sensor/temperature"}
	case "MQTTConnectionStatus", "MQTTReconnectsTotal":
		return map[string]string{"broker": "10.10.10.2:1883"}
	case "MQTTConnectionErrors":
		return map[string]string{"broker": "10.10.10.2:1883", "error_type": "connection_refused"}
	default:
		return map[string]string{}
	}
}

func (s *Server) getExampleValue(metricName string) string {
	switch metricName {
	case "VersionInfo":
		return "1"
	case "MQTTMessageCount":
		return "42"
	case "MQTTMessageBytes":
		return "1024"
	case "MQTTConnectionStatus":
		return "1"
	case "MQTTConnectionErrors":
		return "3"
	case "MQTTReconnectsTotal":
		return "2"
	case "MQTTTopicLastMessage":
		return "1704067200"
	default:
		return "0"
	}
}

func (s *Server) getMetricHelp(metricName string) string {
	switch metricName {
	case "VersionInfo":
		return "Information about the MQTT exporter"
	case "MQTTMessageCount":
		return "Total number of MQTT messages received"
	case "MQTTMessageBytes":
		return "Total number of bytes received in MQTT messages"
	case "MQTTConnectionStatus":
		return "MQTT connection status (1 = connected, 0 = disconnected)"
	case "MQTTConnectionErrors":
		return "Total number of MQTT connection errors"
	case "MQTTReconnectsTotal":
		return "Total number of MQTT reconnection attempts"
	case "MQTTTopicLastMessage":
		return "Timestamp of the last message received per topic"
	default:
		return "MQTT exporter metric"
	}
}

func (s *Server) handleRoot(c *gin.Context) {
	versionInfo := version.Get()
	metricsInfo := s.getMetricsInfo()

	// Generate metrics HTML dynamically
	metricsHTML := ""

	for i, metric := range metricsInfo {
		labelsStr := ""

		if len(metric.Labels) > 0 {
			var labelPairs []string
			for k, v := range metric.Labels {
				labelPairs = append(labelPairs, fmt.Sprintf(`%s="%s"`, k, v))
			}

			labelsStr = "{" + strings.Join(labelPairs, ", ") + "}"
		}

		// Create clickable metric with hidden details
		metricsHTML += fmt.Sprintf(`
            <div class="metric-item" onclick="toggleMetricDetails(%d)">
                <div class="metric-header">
                    <span class="metric-name">%s</span>
                    <span class="metric-toggle">▼</span>
                </div>
                <div class="metric-details" id="metric-%d">
                    <div class="metric-help"><strong>Description:</strong> %s</div>
                    <div class="metric-example"><strong>Example:</strong> %s = %s</div>
                    <div class="metric-labels"><strong>Labels:</strong> %s</div>
                </div>
            </div>`,
			i,
			metric.Name,
			i,
			metric.Help,
			metric.Name,
			metric.ExampleValue,
			labelsStr)
	}

	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>MQTT Exporter ` + versionInfo.Version + `</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            max-width: 800px;
            margin: 0 auto;
            padding: 2rem;
            line-height: 1.6;
            color: #333;
        }
        h1 {
            color: #2c3e50;
            border-bottom: 2px solid #3498db;
            padding-bottom: 0.5rem;
        }
        h1 .version {
            font-size: 0.6em;
            color: #6c757d;
            font-weight: normal;
            margin-left: 0.5rem;
        }
        .endpoints-grid {
            display: grid;
            grid-template-columns: repeat(2, 1fr);
            gap: 1rem;
            margin: 1rem 0;
        }
        .endpoint {
            background: #f8f9fa;
            border: 1px solid #e9ecef;
            border-radius: 8px;
            padding: 1rem;
            text-align: center;
            transition: all 0.2s ease;
        }
        .endpoint:hover {
            border-color: #007bff;
            background-color: #e3f2fd;
        }
        .endpoint h3 {
            margin: 0 0 0.5rem 0;
            color: #495057;
        }
        .endpoint a {
            color: #007bff;
            text-decoration: none;
            font-weight: 500;
        }
        .endpoint a:hover {
            text-decoration: underline;
        }
        .description {
            color: #6c757d;
            font-size: 0.9rem;
            margin-bottom: 0.5rem;
        }
        .status {
            display: inline-block;
            padding: 0.25rem 0.5rem;
            border-radius: 4px;
            font-size: 0.8rem;
            font-weight: 500;
        }
        .status.healthy {
            background: #d4edda;
            color: #155724;
        }
        .status.metrics {
            background: #d1ecf1;
            color: #0c5460;
        }
        .status.ready {
            background: #d4edda;
            color: #155724;
        }
        .status.connected {
            background: #d4edda;
            color: #155724;
        }
        .status.disconnected {
            background: #f8d7da;
            color: #721c24;
        }
        .service-status {
            background: #e9ecef;
            border: 1px solid #dee2e6;
            border-radius: 8px;
            padding: 1rem;
            margin: 1rem 0;
        }
        .service-status h3 {
            margin: 0 0 0.5rem 0;
            color: #495057;
        }
        .service-status p {
            margin: 0.25rem 0;
            color: #6c757d;
        }
        .metrics-info {
            background: #e9ecef;
            border: 1px solid #dee2e6;
            border-radius: 8px;
            padding: 1rem;
            margin: 1rem 0;
        }
        .metrics-info h3 {
            margin: 0 0 0.5rem 0;
            color: #495057;
        }
        .metrics-info ul {
            margin: 0.5rem 0;
            padding-left: 1.5rem;
        }
        .metrics-info li {
            margin: 0.25rem 0;
            color: #6c757d;
        }
        .footer {
            margin-top: 2rem;
            padding-top: 1rem;
            border-top: 1px solid #dee2e6;
            text-align: center;
            color: #6c757d;
            font-size: 0.9rem;
        }
        .footer a {
            color: #007bff;
            text-decoration: none;
        }
        .footer a:hover {
            text-decoration: underline;
        }
        .metrics-list {
            margin: 0.5rem 0;
        }
        .metric-item {
            border: 1px solid #dee2e6;
            border-radius: 6px;
            margin: 0.5rem 0;
            cursor: pointer;
            transition: all 0.2s ease;
        }
        .metric-item:hover {
            border-color: #007bff;
            background-color: #f8f9fa;
        }
        .metric-header {
            padding: 0.75rem;
            display: flex;
            justify-content: space-between;
            align-items: center;
            font-weight: 500;
            color: #495057;
        }
        .metric-name {
            font-family: 'Courier New', monospace;
            font-size: 0.9rem;
        }
        .metric-toggle {
            font-size: 0.8rem;
            color: #6c757d;
            transition: transform 0.2s ease;
        }
        .metric-details {
            display: none;
            padding: 0.75rem;
            border-top: 1px solid #dee2e6;
            background-color: #f8f9fa;
            font-size: 0.85rem;
            line-height: 1.4;
        }
        .metric-details.show {
            display: block;
        }
        .metric-help, .metric-example, .metric-labels {
            margin: 0.5rem 0;
        }
        .metric-example {
            font-family: 'Courier New', monospace;
            background-color: #e9ecef;
            padding: 0.25rem 0.5rem;
            border-radius: 3px;
        }
        .metric-labels {
            color: #6c757d;
        }
    </style>
    <script>
        function toggleMetricDetails(id) {
            const details = document.getElementById('metric-' + id);
            const toggle = details.previousElementSibling.querySelector('.metric-toggle');
            
            if (details.classList.contains('show')) {
                details.classList.remove('show');
                toggle.textContent = '▼';
            } else {
                details.classList.add('show');
                toggle.textContent = '▲';
            }
        }
    </script>
</head>
<body>
    <h1>MQTT Exporter<span class="version">` + versionInfo.Version + `</span></h1>
    
    <div class="endpoints-grid">
        <div class="endpoint">
            <h3><a href="/metrics">📊 Metrics</a></h3>
            <p class="description">Prometheus metrics endpoint</p>
            <span class="status metrics">Available</span>
        </div>

        <div class="endpoint">
            <h3><a href="/health">❤️ Health Check</a></h3>
            <p class="description">Service health status</p>
            <span class="status healthy">Healthy</span>
        </div>
    </div>

    <div class="service-status">
        <h3>Service Status</h3>
        <p><strong>Status:</strong> <span class="status ready">Ready</span></p>
        <p><strong>MQTT Connection:</strong> <span class="status connected">Connected</span></p>
        <p><strong>Message Monitoring:</strong> <span class="status ready">Active</span></p>
    </div>

    <div class="metrics-info">
        <h3>Available Metrics</h3>
        <div class="metrics-list">` + metricsHTML + `
        </div>
    </div>

    <div class="metrics-info">
        <h3>Version Information</h3>
        <ul>
            <li><strong>Version:</strong> ` + versionInfo.Version + `</li>
            <li><strong>Commit:</strong> ` + versionInfo.Commit + `</li>
            <li><strong>Build Date:</strong> ` + versionInfo.BuildDate + `</li>
        </ul>
    </div>

    <div class="metrics-info">
        <h3>Configuration</h3>
        <ul>
            <li><strong>MQTT Broker:</strong> ` + s.config.MQTT.Broker + `</li>
            <li><strong>Client ID:</strong> ` + s.config.MQTT.ClientID + `</li>
            <li><strong>Topics:</strong> ` + fmt.Sprintf("%d", len(s.config.MQTT.Topics)) + ` configured</li>
            <li><strong>QoS Level:</strong> ` + fmt.Sprintf("%d", s.config.MQTT.QoS) + `</li>
        </ul>
    </div>

    <div class="footer">
        <p>Copyright © 2025 Dougal Matthews. Licensed under <a href="https://opensource.org/licenses/MIT" target="_blank">MIT License</a>.</p>
        <p><a href="https://github.com/d0ugal/mqtt-exporter" target="_blank">GitHub Repository</a> | <a href="https://github.com/d0ugal/mqtt-exporter/issues" target="_blank">Report Issues</a></p>
    </div>
</body>
</html>`

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}

func (s *Server) handleHealth(c *gin.Context) {
	versionInfo := version.Get()
	c.JSON(http.StatusOK, gin.H{
		"status":     "healthy",
		"timestamp":  time.Now().Unix(),
		"service":    "mqtt-exporter",
		"version":    versionInfo.Version,
		"commit":     versionInfo.Commit,
		"build_date": versionInfo.BuildDate,
	})
}

func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port)

	s.server = &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	slog.Info("Starting HTTP server", "address", addr)

	return s.server.ListenAndServe()
}

func (s *Server) Shutdown() error {
	slog.Info("Shutting down HTTP server")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return s.server.Shutdown(ctx)
}
