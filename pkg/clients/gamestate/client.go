package gamestate

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"google.golang.org/protobuf/proto"

	"github.com/c4t-but-s4d/fastad/internal/models"
	gspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/game_state"
	versionpb "github.com/c4t-but-s4d/fastad/pkg/proto/data/version"
)

var ErrStateUnavailable = errors.New("game state is unavailable")

type Client struct {
	c gspb.GameStateServiceClient

	refreshMu sync.Mutex
	version   *versionpb.Version

	cache *Cache
}

func NewClient(c gspb.GameStateServiceClient) *Client {
	return &Client{c: c, cache: NewCache()}
}

func (c *Client) Get(ctx context.Context) (*models.GameState, error) {
	if err := c.refresh(ctx); err != nil {
		return nil, fmt.Errorf("refreshing services: %w", err)
	}
	s := c.cache.GetState()
	if s == nil {
		return nil, ErrStateUnavailable
	}
	return c.cache.GetState(), nil
}

func (c *Client) Update(ctx context.Context, req *gspb.UpdateRequest) (*models.GameState, error) {
	resp, err := c.c.Update(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("updating state: %w", err)
	}
	return models.NewGameStateFromProto(resp.GameState), nil
}

func (c *Client) UpdateRound(ctx context.Context, req *gspb.UpdateRoundRequest) (*models.GameState, error) {
	resp, err := c.c.UpdateRound(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("updating round: %w", err)
	}
	return models.NewGameStateFromProto(resp.GameState), nil
}

func (c *Client) RawClient() gspb.GameStateServiceClient {
	return c.c
}

func (c *Client) refresh(ctx context.Context) error {
	c.refreshMu.Lock()
	defer c.refreshMu.Unlock()

	resp, err := c.c.Get(ctx, &gspb.GetRequest{Version: c.version})
	if err != nil {
		return fmt.Errorf("getting state: %w", err)
	}

	if proto.Equal(c.version, resp.Version) {
		return nil
	}

	c.version = resp.Version
	c.cache.SetState(models.NewGameStateFromProto(resp.GameState))

	return nil
}
