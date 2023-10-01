package services

import (
	"sync"

	"github.com/c4t-but-s4d/fastad/internal/models"
)

type Cache struct {
	mu    sync.RWMutex
	state *models.GameState
}

func NewCache() *Cache {
	return &Cache{}
}

func (c *Cache) SetState(state *models.GameState) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.state = state
}

func (c *Cache) GetState() *models.GameState {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.state
}
