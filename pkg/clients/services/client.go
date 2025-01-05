package services

import (
	"context"
	"fmt"
	"sync"

	servicespb "github.com/c4t-but-s4d/fastad/pkg/proto/data/services"
	versionpb "github.com/c4t-but-s4d/fastad/pkg/proto/data/version"
)

type Client struct {
	c servicespb.ServicesServiceClient

	refreshMu sync.Mutex
	version   *versionpb.Version

	cache *Cache
}

func NewClient(c servicespb.ServicesServiceClient) *Client {
	return &Client{c: c, cache: NewCache()}
}

func (c *Client) List(ctx context.Context) ([]*servicespb.Service, error) {
	if err := c.refresh(ctx); err != nil {
		return nil, fmt.Errorf("refreshing services: %w", err)
	}
	return c.cache.GetServices(), nil
}

func (c *Client) CreateBatch(ctx context.Context, services []*servicespb.Service) ([]*servicespb.Service, error) {
	resp, err := c.c.CreateBatch(ctx, &servicespb.CreateBatchRequest{Services: services})
	if err != nil {
		return nil, fmt.Errorf("making api request: %w", err)
	}

	if err := c.refresh(ctx); err != nil {
		return nil, fmt.Errorf("refreshing: %w", err)
	}

	return resp.GetServices(), nil
}

func (c *Client) refresh(ctx context.Context) error {
	c.refreshMu.Lock()
	defer c.refreshMu.Unlock()

	resp, err := c.c.List(ctx, &servicespb.ListRequest{Version: c.version})
	if err != nil {
		return fmt.Errorf("getting services: %w", err)
	}

	if c.version.EqualVT(resp.GetVersion()) {
		return nil
	}

	c.version = resp.GetVersion()

	c.cache.SetServices(resp.GetServices())

	return nil
}
