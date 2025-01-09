package main

import (
	"log"
	"net/http"
	"time"

	"github.com/oldmagic/cstv-go/pkg/config"
	"github.com/oldmagic/cstv-go/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

// Prometheus Metrics
var requestCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests",
	},
	[]string{"path"},
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	logger.Init(cfg.LogLevel)
	prometheus.MustRegister(requestCounter)

	// Optimized Fiber App
	app := fiber.New(fiber.Config{
		IdleTimeout:  30 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Prefork:      true, // Enables multiple processes for better performance
	})

	// Middleware
	app.Use(compress.New()) // Enable Gzip compression
	app.Use(recover.New())  // Graceful panic recovery
	app.Use(limiter.New(limiter.Config{
		Max:               1000, // Allows 1000 requests per second per IP
		Expiration:        1 * time.Second,
		LimitReached: func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusTooManyRequests)
		},
	}))
	app.Use(func(c *fiber.Ctx) error {
		requestCounter.WithLabelValues(c.Path()).Inc()
		return c.Next()
	})

	// Prometheus Metrics Endpoint
	http.Handle("/metrics", promhttp.Handler())

	go func() {
		log.Fatal(http.ListenAndServe(":9091", nil))
	}()

	// Start Server
	logrus.Infof("Starting server on port %s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		logrus.Fatalf("Error starting server: %v", err)
	}
}
