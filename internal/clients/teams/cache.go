package teams

import (
	"sync"

	"github.com/c4t-but-s4d/fastad/internal/models"
	"github.com/samber/lo"
)

type Cache struct {
	mu        sync.RWMutex
	teams     []*models.Team
	teamsByID map[int]*models.Team
}

func NewCache() *Cache {
	return &Cache{
		teamsByID: make(map[int]*models.Team),
	}
}

func (c *Cache) SetTeams(teams []*models.Team) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.teams = teams
	c.teamsByID = lo.KeyBy(teams, func(team *models.Team) int {
		return team.ID
	})
}

func (c *Cache) GetTeams() []*models.Team {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.teams
}

func (c *Cache) GetTeamByID(id int) *models.Team {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.teamsByID[id]
}
