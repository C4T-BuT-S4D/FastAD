package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// CacheMetrics tracks cache hit/miss for a specific client.
// The client name is set as a const label, simplifying usage.
type CacheMetrics struct {
	hits   prometheus.Counter
	misses prometheus.Counter
}

// NewCacheMetrics creates a CacheMetrics instance for a specific client.
// The client name and installation are set as const labels.
func NewCacheMetrics(installation, client string) *CacheMetrics {
	constLabels := prometheus.Labels{"client": client}
	if installation != "" {
		constLabels["installation"] = installation
	}

	return &CacheMetrics{
		hits: promauto.NewCounter(
			prometheus.CounterOpts{
				Name:        "fastad_client_cache_hits_total",
				Help:        "Total number of client cache hits (version matched, no data fetch)",
				ConstLabels: constLabels,
			},
		),
		misses: promauto.NewCounter(
			prometheus.CounterOpts{
				Name:        "fastad_client_cache_misses_total",
				Help:        "Total number of client cache misses (version mismatch, data fetched)",
				ConstLabels: constLabels,
			},
		),
	}
}

// Hit records a cache hit.
func (m *CacheMetrics) Hit() {
	if m == nil {
		return
	}
	m.hits.Inc()
}

// Miss records a cache miss.
func (m *CacheMetrics) Miss() {
	if m == nil {
		return
	}
	m.misses.Inc()
}
