package handlers

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/samber/lo"
	"go.uber.org/zap"

	"github.com/c4t-but-s4d/fastad/internal/centutil"
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

type ScoreboardState struct {
	Teams      []*teamspb.Team
	Services   []*servicespb.Service
	Scoreboard *scoreboardpb.Scoreboard
}

type BoardBuilder struct {
	teamsClient    *teams.Client
	servicesClient *services.Client
	receiverClient receiverpb.ReceiverServiceClient
	slacClient     slacpb.SlacServiceClient

	producer centutil.Producer

	mu sync.RWMutex
	// Never modify these fields, only replace them with new values.
	// References to these fields are returned to the caller.
	teamsCache      []*teamspb.Team
	servicesCache   []*servicespb.Service
	scoreboardCache *scoreboardpb.Scoreboard

	logger *zap.Logger
}

func NewBoardBuilder(
	teamsClient *teams.Client,
	servicesClient *services.Client,
	receiverClient receiverpb.ReceiverServiceClient,
	slacClient slacpb.SlacServiceClient,
	producer centutil.Producer,
) *BoardBuilder {
	return &BoardBuilder{
		teamsClient:    teamsClient,
		servicesClient: servicesClient,
		receiverClient: receiverClient,
		slacClient:     slacClient,
		producer:       producer,

		logger: zap.L().Named("board_builder"),
	}
}

func (b *BoardBuilder) Run(ctx context.Context) {
	ticker := time.NewTicker(scoreboardRefreshInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			b.logger.Debug("scoreboard refreshed")
			if err := b.refresh(ctx); err != nil {
				b.logger.Error("refreshing scoreboard", zap.Error(err))
				break
			}
			b.logger.Debug("scoreboard refreshed")

			state, err := b.GetState()
			if err != nil {
				b.logger.Error("getting state", zap.Error(err))
				break
			}

			b.logger.Debug("publishing state")
			if err := b.producer.PublishProto(ctx, state.Scoreboard); err != nil {
				b.logger.Error("publishing state", zap.Error(err))
			}
		}
	}
}

func (b *BoardBuilder) GetState() (*ScoreboardState, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return &ScoreboardState{
		Teams:      b.teamsCache,
		Services:   b.servicesCache,
		Scoreboard: b.scoreboardCache,
	}, nil
}

func (b *BoardBuilder) refresh(ctx context.Context) error {
	teamsList, err := b.teamsClient.List(ctx)
	if err != nil {
		return fmt.Errorf("listing teams: %w", err)
	}

	servicesList, err := b.servicesClient.List(ctx)
	if err != nil {
		return fmt.Errorf("listing services: %w", err)
	}

	sb, err := b.buildScoreboard(ctx, teamsList, servicesList)
	if err != nil {
		return fmt.Errorf("building scoreboard state: %w", err)
	}

	b.mu.Lock()
	defer b.mu.Unlock()
	b.teamsCache = teamsList
	b.servicesCache = servicesList
	b.scoreboardCache = sb

	return nil
}

func (b *BoardBuilder) buildScoreboard(
	ctx context.Context,
	teams []*teamspb.Team,
	services []*servicespb.Service,
) (*scoreboardpb.Scoreboard, error) {
	sbMap := make(map[teamServiceKey]*scoreboardpb.Scoreboard_TeamServiceState)
	for _, team := range teams {
		for _, service := range services {
			key := teamServiceKey{TeamID: team.GetId(), ServiceID: service.GetId()}
			sbMap[key] = &scoreboardpb.Scoreboard_TeamServiceState{
				TeamId:    team.GetId(),
				ServiceId: service.GetId(),
				Points:    service.GetDefaultScore(),
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
