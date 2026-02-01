package cache

import (
	"context"
	"fmt"
	"sync"

	"github.com/c4t-but-s4d/fastad/pkg/clients/metrics"
	versionpb "github.com/c4t-but-s4d/fastad/pkg/proto/data/version"
)

// Fetcher defines how to fetch versioned data from the server.
// If the provided version matches the server's version, implementations should
// return (zero value, currentVersion, nil) to indicate a cache hit.
// If versions differ, return (data, newVersion, nil).
type Fetcher[T any] interface {
	Fetch(ctx context.Context, version *versionpb.Version) (T, *versionpb.Version, error)
}

// FetcherFunc is a function adapter for the Fetcher interface.
type FetcherFunc[T any] func(ctx context.Context, version *versionpb.Version) (T, *versionpb.Version, error)

func (f FetcherFunc[T]) Fetch(ctx context.Context, version *versionpb.Version) (T, *versionpb.Version, error) {
	return f(ctx, version)
}

// VersionedCache is a generic client-side cache with version-based invalidation.
// It uses optimistic versioning: the server returns data only when versions differ.
type VersionedCache[T any] struct {
	fetcher Fetcher[T]

	mu      sync.Mutex
	version *versionpb.Version
	data    T

	metrics *metrics.CacheMetrics
}

// NewVersionedCache creates a new versioned cache.
// The fetcher is called to retrieve data when the cache needs refreshing.
// Installation and clientName are used for metrics labels.
func NewVersionedCache[T any](fetcher Fetcher[T], installation, clientName string) *VersionedCache[T] {
	return &VersionedCache[T]{
		fetcher: fetcher,
		metrics: metrics.NewCacheMetrics(installation, clientName),
	}
}

// Get returns the cached data, refreshing from the server if needed.
// Thread-safe for concurrent access.
func (c *VersionedCache[T]) Get(ctx context.Context) (T, error) {
	if err := c.refresh(ctx); err != nil {
		var zero T
		return zero, fmt.Errorf("refreshing cache: %w", err)
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	return c.data, nil
}

// refresh fetches data from the server if the version has changed.
func (c *VersionedCache[T]) refresh(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	data, newVersion, err := c.fetcher.Fetch(ctx, c.version)
	if err != nil {
		return fmt.Errorf("fetching data: %w", err)
	}

	// Version matched - cache hit (server returned empty data)
	if c.version.EqualVT(newVersion) {
		c.metrics.Hit()
		return nil
	}

	// Version changed - cache miss (server returned fresh data)
	c.metrics.Miss()
	c.version = newVersion
	c.data = data
	return nil
}
