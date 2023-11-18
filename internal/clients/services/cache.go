package services

import (
	"sync"

	"github.com/samber/lo"

	"github.com/c4t-but-s4d/fastad/internal/models"
)

type Cache struct {
	mu           sync.RWMutex
	services     []*models.Service
	servicesByID map[int]*models.Service
}

func NewCache() *Cache {
	return &Cache{
		servicesByID: make(map[int]*models.Service),
	}
}

func (c *Cache) SetServices(services []*models.Service) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.services = services
	c.servicesByID = lo.KeyBy(services, func(service *models.Service) int {
		return service.ID
	})
}

func (c *Cache) GetServices() []*models.Service {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.services
}

func (c *Cache) GetServiceByID(id int) *models.Service {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.servicesByID[id]
}
