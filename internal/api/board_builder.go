package api

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/samber/lo"

	"github.com/c4t-but-s4d/fastad/pkg/clients/services"
	"github.com/c4t-but-s4d/fastad/pkg/clients/teams"
	"github.com/c4t-but-s4d/fastad/pkg/httpext"
	servicespb "github.com/c4t-but-s4d/fastad/pkg/proto/data/services"
	teamspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/teams"
	receiverpb "github.com/c4t-but-s4d/fastad/pkg/proto/receiver"
	scoreboardpb "github.com/c4t-but-s4d/fastad/pkg/proto/scoreboard"
	slacpb "github.com/c4t-but-s4d/fastad/pkg/proto/slac"
)

const scoreboardRefreshInterval = 2 * time.Second

type BoardBuilder struct {
	teamsClient    *teams.Client
	servicesClient *services.Client
	receiverClient receiverpb.ReceiverServiceClient
	slacClient     slacpb.SlacServiceClient

	mu          sync.RWMutex
	lastRefresh time.Time
	// Never modify these fields, only replace them with new values.
	// References to these fields are returned to the caller.
	teamsCache      []*teamspb.Team
	servicesCache   []*servicespb.Service
	scoreboardCache *scoreboardpb.Scoreboard
}

func NewBoardBuilder(
	teamsClient *teams.Client,
	servicesClient *services.Client,
	receiverClient receiverpb.ReceiverServiceClient,
	slacClient slacpb.SlacServiceClient,
) *BoardBuilder {
	return &BoardBuilder{
		teamsClient:    teamsClient,
		servicesClient: servicesClient,
		receiverClient: receiverClient,
		slacClient:     slacClient,
	}
}

func (b *BoardBuilder) GetScoreboard(ctx context.Context) (*scoreboardpb.Scoreboard, error) {
	b.mu.RLock()
	board := b.scoreboardCache
	lastRefresh := b.lastRefresh
	b.mu.RUnlock()

	if board != nil && time.Since(lastRefresh) < scoreboardRefreshInterval {
		return board, nil
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	if board != nil && time.Since(b.lastRefresh) < scoreboardRefreshInterval {
		return b.scoreboardCache, nil
	}

	if err := b.refreshUnlocked(ctx); err != nil {
		return nil, fmt.Errorf("refreshing scoreboard: %w", err)
	}

	return b.scoreboardCache, nil
}

func (b *BoardBuilder) GetTeams(ctx context.Context) ([]*teamspb.Team, error) {
	b.mu.RLock()
	cache := b.teamsCache
	lastRefresh := b.lastRefresh
	b.mu.RUnlock()

	if len(cache) > 0 && time.Since(lastRefresh) < scoreboardRefreshInterval {
		return cache, nil
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	if len(cache) > 0 && time.Since(lastRefresh) < scoreboardRefreshInterval {
		return cache, nil
	}

	if err := b.refreshUnlocked(ctx); err != nil {
		return nil, fmt.Errorf("refreshing scoreboard: %w", err)
	}

	return b.teamsCache, nil
}

func (b *BoardBuilder) refreshUnlocked(ctx context.Context) error {
	if err := b.getTeamsUnlocked(ctx); err != nil {
		return fmt.Errorf("getting teams: %w", err)
	}

	if err := b.getServicesUnlocked(ctx); err != nil {
		return fmt.Errorf("getting services: %w", err)
	}

	sb, err := b.buildScoreboardStateUnlocked(ctx)
	if err != nil {
		return fmt.Errorf("building scoreboard state: %w", err)
	}

	b.scoreboardCache = sb
	b.lastRefresh = time.Now()

	return nil
}

func (b *BoardBuilder) getTeamsUnlocked(ctx context.Context) error {
	teamsList, err := b.teamsClient.List(ctx)
	if err != nil {
		return fmt.Errorf("listing teams: %w", err)
	}

	b.teamsCache = teamsList
	return nil
}

func (b *BoardBuilder) getServicesUnlocked(ctx context.Context) error {
	servicesList, err := b.servicesClient.List(ctx)
	if err != nil {
		return fmt.Errorf("listing services: %w", err)
	}

	b.servicesCache = servicesList
	return nil
}

func (b *BoardBuilder) buildScoreboardStateUnlocked(ctx context.Context) (*scoreboardpb.Scoreboard, error) {
	sbMap := make(map[teamServiceKey]*scoreboardpb.Scoreboard_TeamServiceState)
	for _, team := range b.teamsCache {
		for _, service := range b.servicesCache {
			sbMap[teamServiceKey{TeamID: team.Id, ServiceID: service.Id}] = &scoreboardpb.Scoreboard_TeamServiceState{
				TeamId:    team.Id,
				ServiceId: service.Id,
				Points:    service.DefaultScore,
			}
		}
	}

	slaState, err := b.slacClient.GetState(ctx, &slacpb.GetStateRequest{})
	if err != nil {
		return nil, httpext.NewErrorFromStatus(err, "getting slac state")
	}

	receiverState, err := b.receiverClient.GetState(ctx, &receiverpb.GetStateRequest{})
	if err != nil {
		return nil, httpext.NewErrorFromStatus(err, "getting receiver state")
	}

	for _, tss := range slaState.GetState().GetTeamServiceStates() {
		key := teamServiceKey{TeamID: tss.GetTeamId(), ServiceID: tss.GetServiceId()}
		if sbs, ok := sbMap[key]; ok {
			sbs.ChecksTotal = tss.GetChecksTotal()
			sbs.ChecksPassed = tss.GetChecksPassed()
			sbs.CheckStatuses = tss.GetCheckStatuses()
		}
	}

	for _, tss := range receiverState.GetState().GetTeamServices() {
		key := teamServiceKey{TeamID: tss.GetTeamId(), ServiceID: tss.GetServiceId()}
		if sbs, ok := sbMap[key]; ok {
			sbs.FlagsStolen = tss.GetFlagsStolen()
			sbs.FlagsLost = tss.GetFlagsLost()
			sbs.Points = tss.GetPoints()
		}
	}

	return &scoreboardpb.Scoreboard{TeamServiceStates: lo.Values(sbMap)}, nil
}

type teamServiceKey struct {
	TeamID    int64
	ServiceID int64
}
