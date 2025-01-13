package Prometheus

import (
	prometheusMetrics "TODO-LIST/Handlers/metrics" // Import the metrics package where the ResponseWriter is defined
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"strconv"
)

// Middleware to count HTTP requests and capture status codes
func CountRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a custom responseWriter to capture the status code
		wrw := &prometheusMetrics.ResponseWriter{ResponseWriter: w}

		// Call the next handler in the chain
		next.ServeHTTP(wrw, r)

		// Increment the Prometheus counter with method, endpoint, and status code
		prometheusMetrics.HttpRequestCount.WithLabelValues(r.Method, r.URL.Path, strconv.Itoa(wrw.StatusCode)).Inc()
	})
}

// Expose Prometheus metrics
func ExposeMetrics() {
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":2112", nil)) // Expose Prometheus metrics on port 2112
}
