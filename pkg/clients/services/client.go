package services

import (
	"context"
	"fmt"

	"github.com/c4t-but-s4d/fastad/pkg/clients/cache"
	servicespb "github.com/c4t-but-s4d/fastad/pkg/proto/data/services"
	versionpb "github.com/c4t-but-s4d/fastad/pkg/proto/data/version"
)

type Client struct {
	grpc  servicespb.ServicesServiceClient
	cache *cache.VersionedCache[[]*servicespb.Service]
}

func NewClient(c servicespb.ServicesServiceClient, installation string) *Client {
	client := &Client{grpc: c}
	client.cache = cache.NewVersionedCache(
		cache.FetcherFunc[[]*servicespb.Service](client.fetch),
		installation,
		"services",
	)
	return client
}

func (c *Client) fetch(ctx context.Context, version *versionpb.Version) ([]*servicespb.Service, *versionpb.Version, error) {
	resp, err := c.grpc.List(ctx, &servicespb.ListRequest{Version: version})
	if err != nil {
		return nil, nil, fmt.Errorf("getting services: %w", err)
	}
	return resp.GetServices(), resp.GetVersion(), nil
}

func (c *Client) List(ctx context.Context) ([]*servicespb.Service, error) {
	services, err := c.cache.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting cached services: %w", err)
	}
	return services, nil
}

func (c *Client) CreateBatch(ctx context.Context, services []*servicespb.Service) ([]*servicespb.Service, error) {
	resp, err := c.grpc.CreateBatch(ctx, &servicespb.CreateBatchRequest{Services: services})
	if err != nil {
		return nil, fmt.Errorf("making api request: %w", err)
	}

	// Force cache refresh after write
	if _, err := c.cache.Get(ctx); err != nil {
		return nil, fmt.Errorf("refreshing cache: %w", err)
	}

	return resp.GetServices(), nil
}

func (c *Client) Update(ctx context.Context, req *servicespb.UpdateRequest) (*servicespb.Service, error) {
	resp, err := c.grpc.Update(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("making api request: %w", err)
	}

	// Force cache refresh after write
	if _, err := c.cache.Get(ctx); err != nil {
		return nil, fmt.Errorf("refreshing cache: %w", err)
	}

	return resp.GetService(), nil
}
