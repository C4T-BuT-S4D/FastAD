package cache_test

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/c4t-but-s4d/fastad/pkg/clients/cache"
	versionpb "github.com/c4t-but-s4d/fastad/pkg/proto/data/version"
)

func newVersion(v int64) *versionpb.Version {
	return &versionpb.Version{Version: v}
}

func TestVersionedCache_Get_CacheMiss(t *testing.T) {
	fetchCount := 0
	fetcher := cache.FetcherFunc[string](func(_ context.Context, _ *versionpb.Version) (string, *versionpb.Version, error) {
		fetchCount++
		return "data", newVersion(1), nil
	})

	c := cache.NewVersionedCache(fetcher, "test", "cache_miss")

	result, err := c.Get(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "data", result)
	assert.Equal(t, 1, fetchCount)
}

func TestVersionedCache_Get_CacheHit(t *testing.T) {
	fetchCount := 0
	fetcher := cache.FetcherFunc[string](func(_ context.Context, version *versionpb.Version) (string, *versionpb.Version, error) {
		fetchCount++
		// First call: return data with version 1
		// Subsequent calls: version matches, return same version (simulating server behavior)
		if version == nil {
			return "data", newVersion(1), nil
		}
		// Server returns same version = cache hit (no data needed)
		return "", version, nil
	})

	c := cache.NewVersionedCache(fetcher, "test", "cache_hit")

	// First call - cache miss
	result1, err := c.Get(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "data", result1)
	assert.Equal(t, 1, fetchCount)

	// Second call - cache hit (version matches)
	result2, err := c.Get(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "data", result2) // Returns cached data
	assert.Equal(t, 2, fetchCount)   // Fetcher was called but returned same version
}

func TestVersionedCache_Get_VersionChange(t *testing.T) {
	currentVersion := int64(1)
	fetcher := cache.FetcherFunc[string](func(_ context.Context, version *versionpb.Version) (string, *versionpb.Version, error) {
		v := atomic.LoadInt64(&currentVersion)
		if version != nil && version.GetVersion() == v {
			return "", version, nil
		}
		return "data-v" + string(rune('0'+v)), newVersion(v), nil
	})

	c := cache.NewVersionedCache(fetcher, "test", "version_change")

	// First fetch
	result1, err := c.Get(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "data-v1", result1)

	// Same version - should return cached
	result2, err := c.Get(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "data-v1", result2)

	// Bump version
	atomic.StoreInt64(&currentVersion, 2)

	// Should get new data
	result3, err := c.Get(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "data-v2", result3)
}

func TestVersionedCache_Get_FetchError(t *testing.T) {
	expectedErr := errors.New("fetch failed")
	fetcher := cache.FetcherFunc[string](func(_ context.Context, _ *versionpb.Version) (string, *versionpb.Version, error) {
		return "", nil, expectedErr
	})

	c := cache.NewVersionedCache(fetcher, "test", "fetch_error")

	result, err := c.Get(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "fetch failed")
	assert.Empty(t, result)
}

func TestVersionedCache_Get_Concurrent(t *testing.T) {
	var fetchCount atomic.Int32
	fetcher := cache.FetcherFunc[int](func(_ context.Context, version *versionpb.Version) (int, *versionpb.Version, error) {
		fetchCount.Add(1)
		if version != nil && version.GetVersion() == 1 {
			return 0, version, nil
		}
		return 42, newVersion(1), nil
	})

	c := cache.NewVersionedCache(fetcher, "test", "concurrent")

	const goroutines = 100
	var wg sync.WaitGroup
	wg.Add(goroutines)

	results := make([]int, 0, goroutines)
	errs := make([]error, 0, goroutines)
	var mu sync.Mutex

	for range goroutines {
		go func() {
			defer wg.Done()
			result, err := c.Get(context.Background())
			mu.Lock()
			results = append(results, result)
			errs = append(errs, err)
			mu.Unlock()
		}()
	}

	wg.Wait()

	// All should succeed with same value
	for i := range goroutines {
		require.NoError(t, errs[i])
		assert.Equal(t, 42, results[i])
	}
}

func TestVersionedCache_Get_NilData(t *testing.T) {
	fetcher := cache.FetcherFunc[*string](func(_ context.Context, _ *versionpb.Version) (*string, *versionpb.Version, error) {
		return nil, newVersion(1), nil
	})

	c := cache.NewVersionedCache(fetcher, "test", "nil_data")

	result, err := c.Get(context.Background())
	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestVersionedCache_Get_SliceData(t *testing.T) {
	fetcher := cache.FetcherFunc[[]string](func(_ context.Context, version *versionpb.Version) ([]string, *versionpb.Version, error) {
		if version != nil && version.GetVersion() == 1 {
			return nil, version, nil
		}
		return []string{"a", "b", "c"}, newVersion(1), nil
	})

	c := cache.NewVersionedCache(fetcher, "test", "slice_data")

	result, err := c.Get(context.Background())
	require.NoError(t, err)
	assert.Equal(t, []string{"a", "b", "c"}, result)

	// Second call should return cached slice
	result2, err := c.Get(context.Background())
	require.NoError(t, err)
	assert.Equal(t, []string{"a", "b", "c"}, result2)
}

func TestVersionedCache_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	fetcher := cache.FetcherFunc[string](func(ctx context.Context, _ *versionpb.Version) (string, *versionpb.Version, error) {
		if ctx.Err() != nil {
			return "", nil, ctx.Err()
		}
		return "data", newVersion(1), nil
	})

	c := cache.NewVersionedCache(fetcher, "test", "ctx_cancel")

	_, err := c.Get(ctx)
	require.Error(t, err)
	assert.ErrorIs(t, err, context.Canceled)
}
