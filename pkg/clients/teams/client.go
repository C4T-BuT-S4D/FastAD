package teams

import (
	"context"
	"fmt"
	"sync"

	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/c4t-but-s4d/fastad/pkg/clients/cache"
	teamspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/teams"
	versionpb "github.com/c4t-but-s4d/fastad/pkg/proto/data/version"
)

type Client struct {
	grpc  teamspb.TeamsServiceClient
	cache *cache.VersionedCache[[]*teamspb.Team]

	// Secondary index for fast token lookup.
	// Updated lazily when List is called.
	indexMu      sync.RWMutex
	teamsByToken map[string]*teamspb.Team
}

func NewClient(c teamspb.TeamsServiceClient, installation string) *Client {
	client := &Client{
		grpc:         c,
		teamsByToken: make(map[string]*teamspb.Team),
	}
	client.cache = cache.NewVersionedCache(
		cache.FetcherFunc[[]*teamspb.Team](client.fetch),
		installation,
		"teams",
	)
	return client
}

func (c *Client) fetch(ctx context.Context, version *versionpb.Version) ([]*teamspb.Team, *versionpb.Version, error) {
	resp, err := c.grpc.List(ctx, &teamspb.ListRequest{Version: version})
	if err != nil {
		return nil, nil, fmt.Errorf("getting teams: %w", err)
	}
	return resp.GetTeams(), resp.GetVersion(), nil
}

func (c *Client) List(ctx context.Context) ([]*teamspb.Team, error) {
	teams, err := c.cache.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting cached teams: %w", err)
	}
	c.updateIndex(teams)
	return teams, nil
}

func (c *Client) GetByToken(ctx context.Context, token string) (*teamspb.Team, error) {
	// Ensure cache is fresh and index is updated
	if _, err := c.List(ctx); err != nil {
		return nil, err
	}

	c.indexMu.RLock()
	team := c.teamsByToken[token]
	c.indexMu.RUnlock()

	if team == nil {
		return nil, status.Error(codes.NotFound, "team not found")
	}
	return team, nil
}

func (c *Client) CreateBatch(ctx context.Context, teams []*teamspb.Team) ([]*teamspb.Team, error) {
	resp, err := c.grpc.CreateBatch(ctx, &teamspb.CreateBatchRequest{Teams: teams})
	if err != nil {
		return nil, fmt.Errorf("making api request: %w", err)
	}

	// Force cache refresh after write
	if _, err := c.List(ctx); err != nil {
		return nil, fmt.Errorf("refreshing cache: %w", err)
	}

	return resp.GetTeams(), nil
}

func (c *Client) Update(ctx context.Context, req *teamspb.UpdateRequest) (*teamspb.Team, error) {
	resp, err := c.grpc.Update(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("making api request: %w", err)
	}

	// Force cache refresh after write
	if _, err := c.List(ctx); err != nil {
		return nil, fmt.Errorf("refreshing cache: %w", err)
	}

	return resp.GetTeam(), nil
}

func (c *Client) updateIndex(teams []*teamspb.Team) {
	c.indexMu.Lock()
	defer c.indexMu.Unlock()
	c.teamsByToken = lo.KeyBy(teams, func(team *teamspb.Team) string {
		return team.GetToken()
	})
}
