package main

import (
	"strconv"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "HTTP requests by route, method and status.",
	}, []string{"method", "path", "status"})

	httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "HTTP request latency by route.",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "path"})

	clicksTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "clicker_clicks_total",
		Help: "Button clicks handled by this pod.",
	})
)

// metricsMiddleware records count and latency per route template (not raw URL, to keep cardinality low).
func metricsMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		start := time.Now()
		err := next(c)
		path := c.Path()
		if path == "/metrics" {
			return err
		}
		status := 0
		if r, uErr := echo.UnwrapResponse(c.Response()); uErr == nil {
			status = r.Status
		}
		method := c.Request().Method
		httpRequests.WithLabelValues(method, path, strconv.Itoa(status)).Inc()
		httpDuration.WithLabelValues(method, path).Observe(time.Since(start).Seconds())
		return err
	}
}
