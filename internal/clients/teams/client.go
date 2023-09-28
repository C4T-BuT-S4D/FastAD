package teams

import (
	"context"
	"fmt"
	"sync"

	"github.com/c4t-but-s4d/fastad/internal/models"
	teamspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/teams"
	"github.com/samber/lo"
)

type Client struct {
	c teamspb.TeamsServiceClient

	refreshMu  sync.Mutex
	lastUpdate int64

	cache *Cache
}

func NewClient(c teamspb.TeamsServiceClient) *Client {
	return &Client{c: c, cache: NewCache()}
}

func (c *Client) List(ctx context.Context) ([]*models.Team, error) {
	if err := c.refresh(ctx); err != nil {
		return nil, fmt.Errorf("refreshing teams: %w", err)
	}
	return c.cache.GetTeams(), nil
}

func (c *Client) GetByID(ctx context.Context, id int) (*models.Team, error) {
	if err := c.refresh(ctx); err != nil {
		return nil, fmt.Errorf("refreshing teams: %w", err)
	}
	return c.cache.GetTeamByID(id), nil
}

func (c *Client) CreateBatch(ctx context.Context, teams []*teamspb.Team) ([]*models.Team, error) {
	resp, err := c.c.CreateBatch(ctx, &teamspb.CreateBatchRequest{Teams: teams})
	if err != nil {
		return nil, fmt.Errorf("making api request: %w", err)
	}

	teamModels := lo.Map(resp.Teams, func(team *teamspb.Team, _ int) *models.Team {
		return models.NewTeamFromProto(team)
	})
	c.cache.SetTeams(teamModels)

	return teamModels, nil
}

func (c *Client) refresh(ctx context.Context) error {
	c.refreshMu.Lock()
	defer c.refreshMu.Unlock()

	resp, err := c.c.List(ctx, &teamspb.ListRequest{LastUpdate: c.lastUpdate})
	if err != nil {
		return fmt.Errorf("getting teams: %w", err)
	}

	if resp.LastUpdate == c.lastUpdate {
		return nil
	}

	c.lastUpdate = resp.LastUpdate

	teamModels := lo.Map(resp.Teams, func(team *teamspb.Team, _ int) *models.Team {
		return models.NewTeamFromProto(team)
	})
	c.cache.SetTeams(teamModels)

	return nil
}
