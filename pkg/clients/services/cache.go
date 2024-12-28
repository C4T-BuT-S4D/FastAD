package services

import (
	"sync"

	servicespb "github.com/c4t-but-s4d/fastad/pkg/proto/data/services"
)

type Cache struct {
	mu       sync.RWMutex
	services []*servicespb.Service
}

func NewCache() *Cache {
	return &Cache{}
}

func (c *Cache) SetServices(services []*servicespb.Service) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.services = services
}

func (c *Cache) GetServices() []*servicespb.Service {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.services
}
