package teams

import (
	"sync"

	"github.com/samber/lo"

	teamspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/teams"
)

type Cache struct {
	mu           sync.RWMutex
	teams        []*teamspb.Team
	teamsByToken map[string]*teamspb.Team
}

func NewCache() *Cache {
	return &Cache{
		teamsByToken: make(map[string]*teamspb.Team),
	}
}

func (c *Cache) SetTeams(teams []*teamspb.Team) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.teams = teams
	c.teamsByToken = lo.KeyBy(teams, func(team *teamspb.Team) string {
		return team.GetToken()
	})
}

func (c *Cache) GetTeams() []*teamspb.Team {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.teams
}

func (c *Cache) GetTeamByToken(token string) *teamspb.Team {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.teamsByToken[token]
}
