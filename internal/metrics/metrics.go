package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	registry *prometheus.Registry
	once     sync.Once

	// BackendLatency measures request duration for LLM backend calls
	BackendLatency *prometheus.HistogramVec

	// TokensTotal tracks token usage per model and type
	TokensTotal *prometheus.CounterVec

	// RequestsTotal tracks total requests by backend, model and status
	RequestsTotal *prometheus.CounterVec
)

// InitMetrics initializes and registers all metrics
func InitMetrics() {
	once.Do(func() {
		registry = prometheus.NewRegistry()

		BackendLatency = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "aigateway_backend_request_duration_seconds",
				Help:    "Time spent processing request by the LLM backend",
				Buckets: []float64{.1, .5, 1, 2.5, 5, 10, 20, 30, 60},
			},
			[]string{"backend", "model", "status"},
		)

		TokensTotal = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "aigateway_model_tokens_total",
				Help: "Total number of tokens processed by model and type",
			},
			[]string{"model", "type"},
		)

		RequestsTotal = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "aigateway_requests_total",
				Help: "Total number of requests processed",
			},
			[]string{"backend", "model", "status"},
		)

		// Register metrics
		registry.MustRegister(BackendLatency)
		registry.MustRegister(TokensTotal)
		registry.MustRegister(RequestsTotal)
	})
}

// GetRegistry returns the metrics registry
func GetRegistry() *prometheus.Registry {
	return registry
}
