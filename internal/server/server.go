package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
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

func (s *Server) handleRoot(c *gin.Context) {
	versionInfo := version.Get()
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
        .endpoint {
            background: #f8f9fa;
            border: 1px solid #e9ecef;
            border-radius: 8px;
            padding: 1rem;
            margin: 1rem 0;
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
    </style>
</head>
<body>
    <h1>MQTT Exporter<span class="version">` + versionInfo.Version + `</span></h1>
    
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

    <div class="service-status">
        <h3>Service Status</h3>
        <p><strong>Status:</strong> <span class="status ready">Ready</span></p>
        <p><strong>MQTT Connection:</strong> <span class="status connected">Connected</span></p>
        <p><strong>Message Monitoring:</strong> <span class="status ready">Active</span></p>
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

    <div class="metrics-info">
        <h3>Available Metrics</h3>
        <ul>
            <li><strong>mqtt_messages_total:</strong> Total messages received per topic</li>
            <li><strong>mqtt_message_bytes_total:</strong> Total bytes received per topic</li>
            <li><strong>mqtt_connection_status:</strong> MQTT connection status</li>
            <li><strong>mqtt_connection_errors_total:</strong> Connection error tracking</li>
            <li><strong>mqtt_topic_last_message_timestamp:</strong> Last message timestamp per topic</li>
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
