package main

import (
	"log"
	"net/http"
	"time"

	"github.com/StepOne-ai/pvz_avito/internal/db"
	"github.com/StepOne-ai/pvz_avito/internal/logger"
	"github.com/StepOne-ai/pvz_avito/internal/routes"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	responseTimeHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_time_seconds",
			Help:    "Histogram of response times for HTTP requests",
			Buckets: []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
		[]string{"method", "endpoint"},
	)
)

var Logger *logrus.Logger

func main() {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(responseTimeHistogram)

	Logger = logger.InitializeLogger()

	if db.InitDB("database.db") != nil {
		log.Fatal("Error accessing db")
	}
	Logger.Info("Database initialized successfully")

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()

		httpRequestsTotal.WithLabelValues(c.Request.Method, c.FullPath(), http.StatusText(c.Writer.Status())).Inc()

		duration := time.Since(start).Seconds()
		responseTimeHistogram.WithLabelValues(c.Request.Method, c.FullPath()).Observe(duration)
	})

	routes.SetupRoutes(r)

	r.GET("/metrics", func(c *gin.Context) {
		promhttp.Handler().ServeHTTP(c.Writer, c.Request)
	})

	r.Run(":8080")
	Logger.Info("Application started")
}
