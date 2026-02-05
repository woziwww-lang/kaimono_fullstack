package metrics

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method", "path", "status"},
	)
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latency in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)
)

func Init() {
	prometheus.MustRegister(requestCount, requestDuration)
}

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		status := c.Writer.Status()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		requestCount.WithLabelValues(c.Request.Method, path, strconv.Itoa(status)).Inc()
		requestDuration.WithLabelValues(c.Request.Method, path, strconv.Itoa(status)).
			Observe(time.Since(start).Seconds())
	}
}

func Handler() gin.HandlerFunc {
	return gin.WrapH(promhttp.Handler())
}
