// Copyright Envoy AI Gateway Authors
// SPDX-License-Identifier: Apache-2.0
// The full text of the Apache license is available in the LICENSE file at
// the root of the repo.

package metrics

import "time"

// TokenMetrics tracks token latency metrics for LLM responses
type TokenMetrics struct {
	metrics        *Metrics
	firstTokenSent bool
	requestStart   time.Time
	lastTokenTime  time.Time
}

// NewTokenMetrics creates a new TokenMetrics instance
func NewTokenMetrics() *TokenMetrics {
	m := GetOrCreate()
	return &TokenMetrics{
		metrics: m,
	}
}

// StartRequest initializes timing for a new request
func (t *TokenMetrics) StartRequest(backendName, modelName string) {
	t.requestStart = time.Now()
	t.firstTokenSent = false
}

// UpdateTokenMetrics updates the token metrics
func (t *TokenMetrics) UpdateTokenMetrics(backendName, modelName string, outputTokens, inputTokens, totalTokens uint32) {
	t.metrics.TokensTotal.WithLabelValues(backendName, modelName, "completion").Add(float64(outputTokens))
	t.metrics.TokensTotal.WithLabelValues(backendName, modelName, "prompt").Add(float64(inputTokens))
	t.metrics.TokensTotal.WithLabelValues(backendName, modelName, "total").Add(float64(totalTokens))
}

// UpdateLatencyMetrics updates the latency metrics
func (t *TokenMetrics) UpdateLatencyMetrics(backendName, modelName string, outputTokens uint32) {
	now := time.Now()
	if !t.firstTokenSent {
		t.firstTokenSent = true
		t.metrics.FirstTokenLatency.WithLabelValues(backendName, modelName).
			Observe(now.Sub(t.requestStart).Seconds())
	} else {
		// Calculate the time between tokens.
		// Since we are only interested in the time between tokens, and streaming is by chunk,
		// we can calculate the time between tokens by the time between the last token and the current token, divided by the number of tokens.
		// And in some cases, the number of tokens can be 0, so we need to check for that.
		div := outputTokens
		if div == 0 {
			div = 1
		}
		itl := now.Sub(t.lastTokenTime).Seconds() / float64(div)
		t.metrics.InterTokenLatency.WithLabelValues(backendName, modelName).
			Observe(itl)
	}
	t.lastTokenTime = now
}

// UpdateRequestMetrics updates request-related metrics
func (t *TokenMetrics) UpdateRequestMetrics(backendName, modelName string, status string) {
	t.metrics.RequestsTotal.WithLabelValues(backendName, modelName, status).Inc()
	if status == "success" {
		t.metrics.TotalLatency.WithLabelValues(backendName, modelName, status).
			Observe(time.Since(t.requestStart).Seconds())
	}
}
