package gamestate

import (
	"sync"

	gspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/game_state"
)

type Cache struct {
	mu    sync.RWMutex
	state *gspb.GameState
}

func NewCache() *Cache {
	return &Cache{}
}

func (c *Cache) SetState(state *gspb.GameState) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.state = state
}

func (c *Cache) GetState() *gspb.GameState {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.state
}
