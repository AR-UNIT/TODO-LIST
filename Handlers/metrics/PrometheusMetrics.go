package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
)

// Define a request counter
var (
	HttpRequestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method", "endpoint", "status_code"},
	)
)

func init() {
	// Register the metrics
	prometheus.MustRegister(HttpRequestCount)
}

// Custom response writer to capture status code
type ResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func (rw *ResponseWriter) WriteHeader(code int) {
	rw.StatusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
