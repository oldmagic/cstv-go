package main

import (
	"log"
	"net/http"

	"github.com/FlowingSPDG/gotv-plus-go/pkg/config"
	"github.com/FlowingSPDG/gotv-plus-go/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

var requestCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests",
	},
	[]string{"path"},
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logger.Init(cfg.LogLevel)

	// Register Prometheus metrics
	prometheus.MustRegister(requestCounter)

	// Create Fiber app
	app := fiber.New()

	// Middleware to count requests
	app.Use(func(c *fiber.Ctx) error {
		requestCounter.WithLabelValues(c.Path()).Inc()
		return c.Next()
	})

	// Prometheus endpoint
	http.Handle("/metrics", promhttp.Handler())

	// Start HTTP server for Prometheus metrics
	go func() {
		log.Fatal(http.ListenAndServe(":9091", nil))
	}()

	// Start Fiber server
	logrus.Infof("Starting server on port %s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		logrus.Fatalf("Error starting server: %v", err)
	}
}
