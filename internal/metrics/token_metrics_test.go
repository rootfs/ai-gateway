package metrics

import (
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"

	"github.com/stretchr/testify/assert"
)

func TestNewTokenMetrics(t *testing.T) {
	tm := NewTokenMetrics()
	assert.NotNil(t, tm)
	assert.NotNil(t, tm.metrics)
	assert.False(t, tm.firstTokenSent)
}

func TestStartRequest(t *testing.T) {
	tm := NewTokenMetrics()
	before := time.Now()
	tm.StartRequest("test-backend", "test-model")
	after := time.Now()

	assert.False(t, tm.firstTokenSent)
	assert.True(t, tm.requestStart.After(before) || tm.requestStart.Equal(before))
	assert.True(t, tm.requestStart.Before(after) || tm.requestStart.Equal(after))
}

func TestUpdateTokenMetrics(t *testing.T) {
	// Reset the default registry to avoid conflicts with other tests
	prometheus.DefaultRegisterer = prometheus.NewRegistry()

	tm := NewTokenMetrics()
	tm.UpdateTokenMetrics("test-backend", "test-model", 10, 5, 15)

	// Get the current value of the metrics
	completion := getCounterValue(t, tm.metrics.TokensTotal, map[string]string{
		"backend": "test-backend",
		"model":   "test-model",
		"type":    "completion",
	})
	prompt := getCounterValue(t, tm.metrics.TokensTotal, map[string]string{
		"backend": "test-backend",
		"model":   "test-model",
		"type":    "prompt",
	})
	total := getCounterValue(t, tm.metrics.TokensTotal, map[string]string{
		"backend": "test-backend",
		"model":   "test-model",
		"type":    "total",
	})

	assert.Equal(t, float64(10), completion)
	assert.Equal(t, float64(5), prompt)
	assert.Equal(t, float64(15), total)
}

func TestUpdateLatencyMetrics(t *testing.T) {
	// Reset the default registry to avoid conflicts with other tests
	prometheus.DefaultRegisterer = prometheus.NewRegistry()

	tm := NewTokenMetrics()
	tm.StartRequest("test-backend", "test-model")

	// Test first token
	time.Sleep(10 * time.Millisecond)
	tm.UpdateLatencyMetrics("test-backend", "test-model", 1)
	assert.True(t, tm.firstTokenSent)

	firstTokenLatency := getHistogramValue(t, tm.metrics.FirstTokenLatency, map[string]string{
		"backend": "test-backend",
		"model":   "test-model",
	})
	assert.Greater(t, firstTokenLatency, 0.0)

	// Test subsequent tokens
	time.Sleep(10 * time.Millisecond)
	tm.UpdateLatencyMetrics("test-backend", "test-model", 2)

	interTokenLatency := getHistogramValue(t, tm.metrics.InterTokenLatency, map[string]string{
		"backend": "test-backend",
		"model":   "test-model",
	})
	assert.Greater(t, interTokenLatency, 0.0)

	// Test zero tokens case
	time.Sleep(10 * time.Millisecond)
	tm.UpdateLatencyMetrics("test-backend", "test-model", 0)
	assert.True(t, tm.firstTokenSent)

	// Test zero tokens case
	time.Sleep(10 * time.Millisecond)
	tm.UpdateLatencyMetrics("test-backend", "test-model", 0)
	assert.True(t, tm.firstTokenSent)
}

func TestUpdateRequestMetrics(t *testing.T) {
	// Reset the default registry to avoid conflicts with other tests
	prometheus.DefaultRegisterer = prometheus.NewRegistry()

	tm := NewTokenMetrics()
	tm.StartRequest("test-backend", "test-model")

	time.Sleep(10 * time.Millisecond)
	tm.UpdateRequestMetrics("test-backend", "test-model", "success")

	// Test requests counter
	requests := getCounterValue(t, tm.metrics.RequestsTotal, map[string]string{
		"backend": "test-backend",
		"model":   "test-model",
		"status":  "success",
	})
	assert.Equal(t, float64(1), requests)

	// Test total latency histogram
	totalLatency := getHistogramValue(t, tm.metrics.TotalLatency, map[string]string{
		"backend": "test-backend",
		"model":   "test-model",
		"status":  "success",
	})
	assert.Greater(t, totalLatency, 0.0)

	// Test failed request
	tm.UpdateRequestMetrics("test-backend", "test-model", "error")
	failedRequests := getCounterValue(t, tm.metrics.RequestsTotal, map[string]string{
		"backend": "test-backend",
		"model":   "test-model",
		"status":  "error",
	})
	assert.Equal(t, float64(1), failedRequests)
}

// Helper function to get the current value of a counter metric
func getCounterValue(t *testing.T, metric *prometheus.CounterVec, labels map[string]string) float64 {
	t.Helper()
	m, err := metric.GetMetricWith(labels)
	assert.NoError(t, err, "Error getting metric")

	metric_pb := &dto.Metric{}
	assert.NoError(t, m.(prometheus.Metric).Write(metric_pb), "Error writing metric")
	return metric_pb.Counter.GetValue()
}

// Helper function to get the current sum of a histogram metric
func getHistogramValue(t *testing.T, metric *prometheus.HistogramVec, labels map[string]string) float64 {
	t.Helper()
	m, err := metric.GetMetricWith(labels)
	assert.NoError(t, err, "Error getting metric")

	metric_pb := &dto.Metric{}
	assert.NoError(t, m.(prometheus.Metric).Write(metric_pb), "Error writing metric")
	return metric_pb.Histogram.GetSampleSum()
}
