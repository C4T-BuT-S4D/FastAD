package gamestate

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/c4t-but-s4d/fastad/pkg/clients/cache"
	gspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/game_state"
	versionpb "github.com/c4t-but-s4d/fastad/pkg/proto/data/version"
)

var ErrStateUnavailable = status.Error(codes.Unavailable, "game state is unavailable")

type Client struct {
	grpc  gspb.GameStateServiceClient
	cache *cache.VersionedCache[*gspb.GameState]
}

func NewClient(c gspb.GameStateServiceClient, installation string) *Client {
	client := &Client{grpc: c}
	client.cache = cache.NewVersionedCache(
		cache.FetcherFunc[*gspb.GameState](client.fetch),
		installation,
		"gamestate",
	)
	return client
}

func (c *Client) fetch(ctx context.Context, version *versionpb.Version) (*gspb.GameState, *versionpb.Version, error) {
	resp, err := c.grpc.Get(ctx, &gspb.GetRequest{Version: version})
	if err != nil {
		return nil, nil, fmt.Errorf("getting state: %w", err)
	}
	return resp.GetGameState(), resp.GetVersion(), nil
}

func (c *Client) Get(ctx context.Context) (*gspb.GameState, error) {
	state, err := c.cache.Get(ctx)
	if err != nil {
		return nil, err
	}
	if state == nil {
		return nil, ErrStateUnavailable
	}
	return state, nil
}

func (c *Client) Create(ctx context.Context, req *gspb.CreateRequest) (*gspb.GameState, error) {
	resp, err := c.grpc.Create(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("creating state: %w", err)
	}
	return resp.GetGameState(), nil
}

func (c *Client) Update(ctx context.Context, req *gspb.UpdateRequest) (*gspb.GameState, error) {
	resp, err := c.grpc.Update(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("updating state: %w", err)
	}
	return resp.GetGameState(), nil
}

func (c *Client) UpdateRound(ctx context.Context, req *gspb.UpdateRoundRequest) (*gspb.GameState, error) {
	resp, err := c.grpc.UpdateRound(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("updating round: %w", err)
	}
	return resp.GetGameState(), nil
}

func (c *Client) FinishGame(ctx context.Context) (*gspb.GameState, error) {
	resp, err := c.grpc.FinishGame(ctx, &gspb.FinishGameRequest{})
	if err != nil {
		return nil, fmt.Errorf("finishing game: %w", err)
	}
	return resp.GetGameState(), nil
}

func (c *Client) RawClient() gspb.GameStateServiceClient {
	return c.grpc
}
