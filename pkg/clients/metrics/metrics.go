package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type CacheMetrics struct {
	hits   prometheus.Counter
	misses prometheus.Counter
}

var (
	metricsCache = make(map[string]*CacheMetrics)
	metricsMu    sync.Mutex
)

func NewCacheMetrics(installation, client string) *CacheMetrics {
	key := installation + "/" + client

	metricsMu.Lock()
	defer metricsMu.Unlock()

	if existing, ok := metricsCache[key]; ok {
		return existing
	}

	constLabels := prometheus.Labels{"client": client}
	if installation != "" {
		constLabels["installation"] = installation
	}

	m := &CacheMetrics{
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

	metricsCache[key] = m
	return m
}

// Hit records a cache hit.
func (m *CacheMetrics) Hit() {
	m.hits.Inc()
}

// Miss records a cache miss.
func (m *CacheMetrics) Miss() {
	m.misses.Inc()
}
