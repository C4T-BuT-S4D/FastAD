package services

import (
	"context"
	"fmt"
	"sync"

	"github.com/c4t-but-s4d/fastad/internal/models"
	servicespb "github.com/c4t-but-s4d/fastad/pkg/proto/data/services"
	"github.com/samber/lo"
)

type Client struct {
	c servicespb.ServicesServiceClient

	refreshMu  sync.Mutex
	lastUpdate int64

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

	resp, err := c.c.List(ctx, &servicespb.ListRequest{LastUpdate: c.lastUpdate})
	if err != nil {
		return fmt.Errorf("getting services: %w", err)
	}

	if resp.LastUpdate == c.lastUpdate {
		return nil
	}

	c.lastUpdate = resp.LastUpdate

	serviceModels := lo.Map(resp.Services, func(service *servicespb.Service, _ int) *models.Service {
		return models.NewServiceFromProto(service)
	})
	c.cache.SetServices(serviceModels)

	return nil
}
