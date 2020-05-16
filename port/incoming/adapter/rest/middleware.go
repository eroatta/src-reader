package rest

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

var (
	requests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "requests",
			Help: "A counter for the requests Golden Signal",
		},
		[]string{"code", "method"},
	)

	latency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "latency",
			Help:    "A histogram of latencies for the latency Golden Signal",
			Buckets: []float64{0.5, 1, 1.5, 2, 2.5, 3, 10, 30},
		},
		[]string{"code", "method"},
	)
)

func setMetricsCollectors(r *gin.Engine) {
	collectors := []prometheus.Collector{requests, latency}
	for _, coll := range collectors {
		err := prometheus.Register(coll)
		if err != nil {
			logrus.WithError(err).Warn("unable to setup metric collector")
		}
	}

	r.Use(goldenSignalsHandler)
}

func goldenSignalsHandler(c *gin.Context) {
	if c.Request.URL.Path == "/metrics" {
		c.Next()
		return
	}

	// metric by endpoint, status code, time, p95, p99
	defer func(t time.Time, path string) {
		code := strconv.Itoa(c.Writer.Status())
		// register metrics
		requests.WithLabelValues(code, path).Inc()
		latency.WithLabelValues(code, path).Observe(time.Since(t).Seconds())
	}(time.Now(), c.Request.URL.Path)

	c.Next()
}
