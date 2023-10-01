package services

import (
	"context"
	"fmt"
	"sync"

	"github.com/c4t-but-s4d/fastad/internal/models"
	servicespb "github.com/c4t-but-s4d/fastad/pkg/proto/data/services"
	versionpb "github.com/c4t-but-s4d/fastad/pkg/proto/data/version"
	"github.com/samber/lo"
	"google.golang.org/protobuf/proto"
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

func (c *Client) List(ctx context.Context) ([]*models.Service, error) {
	if err := c.refresh(ctx); err != nil {
		return nil, fmt.Errorf("refreshing services: %w", err)
	}
	return c.cache.GetServices(), nil
}

func (c *Client) GetByID(ctx context.Context, id int) (*models.Service, error) {
	if err := c.refresh(ctx); err != nil {
		return nil, fmt.Errorf("refreshing services: %w", err)
	}
	return c.cache.GetServiceByID(id), nil
}

func (c *Client) CreateBatch(ctx context.Context, services []*servicespb.Service) ([]*models.Service, error) {
	resp, err := c.c.CreateBatch(ctx, &servicespb.CreateBatchRequest{Services: services})
	if err != nil {
		return nil, fmt.Errorf("making api request: %w", err)
	}

	serviceModels := lo.Map(resp.Services, func(service *servicespb.Service, _ int) *models.Service {
		return models.NewServiceFromProto(service)
	})
	c.cache.SetServices(serviceModels)

	return serviceModels, nil
}

func (c *Client) refresh(ctx context.Context) error {
	c.refreshMu.Lock()
	defer c.refreshMu.Unlock()

	resp, err := c.c.List(ctx, &servicespb.ListRequest{Version: c.version})
	if err != nil {
		return fmt.Errorf("getting services: %w", err)
	}

	if proto.Equal(c.version, resp.Version) {
		return nil
	}

	c.version = resp.Version

	serviceModels := lo.Map(resp.Services, func(service *servicespb.Service, _ int) *models.Service {
		return models.NewServiceFromProto(service)
	})
	c.cache.SetServices(serviceModels)

	return nil
}
