package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/uptrace/bun"
	"go.temporal.io/sdk/client"
	"go.uber.org/atomic"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/c4t-but-s4d/fastad/pkg/clients/gamestate"
	"github.com/c4t-but-s4d/fastad/pkg/clients/services"
	"github.com/c4t-but-s4d/fastad/pkg/clients/teams"
	checkerpb "github.com/c4t-but-s4d/fastad/pkg/proto/checker"
	gspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/game_state"
	servicespb "github.com/c4t-but-s4d/fastad/pkg/proto/data/services"
	teamspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/teams"
)

const CheckManagerRefreshInterval = 5 * time.Second

type teamServiceKey struct {
	teamID    int64
	serviceID int64
	action    checkerpb.Action
}

func (k teamServiceKey) logFields() []zap.Field {
	return []zap.Field{
		zap.Int64("team_id", k.teamID),
		zap.Int64("service_id", k.serviceID),
		zap.String("action", k.action.String()),
	}
}

type CheckManager struct {
	activeSchedulers map[teamServiceKey]*CheckScheduler
	cancellers       map[teamServiceKey]context.CancelFunc
	dones            map[teamServiceKey]chan struct{}

	db              *bun.DB
	gameStateClient *gamestate.Client
	teamsClient     *teams.Client
	servicesClient  *services.Client
	temporalClient  client.Client

	gameState *atomic.Pointer[gspb.GameState]
	teams     []*teamspb.Team
	services  []*servicespb.Service

	logger *zap.Logger
}

func NewCheckManager(
	db *bun.DB,
	gameStateClient *gamestate.Client,
	teamsClient *teams.Client,
	servicesClient *services.Client,
	temporalClient client.Client,
) *CheckManager {
	return &CheckManager{
		activeSchedulers: make(map[teamServiceKey]*CheckScheduler),
		cancellers:       make(map[teamServiceKey]context.CancelFunc),
		dones:            make(map[teamServiceKey]chan struct{}),

		db:              db,
		gameStateClient: gameStateClient,
		teamsClient:     teamsClient,
		servicesClient:  servicesClient,
		temporalClient:  temporalClient,

		gameState: atomic.NewPointer[gspb.GameState](nil),

		logger: zap.L().With(zap.String("component", "check_manager")),
	}
}

func (m *CheckManager) Run(ctx context.Context) error {
	if err := m.refreshData(ctx); err != nil {
		return fmt.Errorf("refreshing data: %w", err)
	}

	refreshTicker := time.NewTicker(CheckManagerRefreshInterval)
	defer refreshTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-refreshTicker.C:
			if err := m.refreshData(ctx); err != nil {
				m.logger.Error("refreshing data", zap.Error(err))
			}
		}
	}
}

func (m *CheckManager) syncSchedulers(ctx context.Context) error {
	needKeys := make(map[teamServiceKey]struct{})

	for _, team := range m.teams {
		for _, service := range m.services {
			for _, action := range []checkerpb.Action{checkerpb.Action_ACTION_CHECK, checkerpb.Action_ACTION_GET} {
				key := teamServiceKey{
					teamID:    team.Id,
					serviceID: service.Id,
					action:    action,
				}
				needKeys[key] = struct{}{}

				logger := m.logger.With(key.logFields()...)

				if current, ok := m.activeSchedulers[key]; ok {
					logger.Debug("scheduler already exists")

					if current.Team.EqualVT(team) && current.Service.EqualVT(service) {
						continue
					}

					logger.Debug("scheduler needs update")
					m.cancellers[key]()
					delete(m.activeSchedulers, key)
					delete(m.cancellers, key)

					logger.Debug("waiting for scheduler to finish")
					select {
					case <-m.dones[key]:
					case <-ctx.Done():
						return nil
					}
					logger.Debug("scheduler finished")
					delete(m.dones, key)
				}

				checkScheduler, err := NewCheckScheduler(
					team,
					service,
					action,
					m.gameState,
					m.db,
					m.temporalClient,
				)
				if err != nil {
					return fmt.Errorf("creating check scheduler: %w", err)
				}
				m.activeSchedulers[key] = checkScheduler

				runCtx, cancel := context.WithCancel(ctx)
				m.cancellers[key] = cancel
				m.dones[key] = make(chan struct{})

				go func() {
					defer close(m.dones[key])
					checkScheduler.Run(runCtx)
				}()
				logger.Info("scheduler started")
			}
		}
	}

	for key := range m.activeSchedulers {
		if _, ok := needKeys[key]; !ok {
			logger := m.logger.With(key.logFields()...)
			logger.Debug("scheduler needs to be removed")

			m.cancellers[key]()
			delete(m.activeSchedulers, key)
			delete(m.cancellers, key)

			logger.Debug("waiting for scheduler to finish")
			select {
			case <-m.dones[key]:
			case <-ctx.Done():
				return nil
			}
			logger.Debug("scheduler finished")
			delete(m.dones, key)
		}
	}

	return nil
}

func (m *CheckManager) refreshData(ctx context.Context) error {
	gs, err := m.gameStateClient.Get(ctx)
	if err != nil {
		if status.Code(err) == codes.Unavailable {
			m.gameState.Store(nil)
			return nil
		}
		return fmt.Errorf("getting game state: %w", err)
	}
	m.gameState.Store(gs)

	teamsList, err := m.teamsClient.List(ctx)
	if err != nil {
		return fmt.Errorf("getting teams: %w", err)
	}
	m.teams = teamsList

	servicesList, err := m.servicesClient.List(ctx)
	if err != nil {
		return fmt.Errorf("getting services: %w", err)
	}
	m.services = servicesList

	if err := m.syncSchedulers(ctx); err != nil {
		return fmt.Errorf("syncing schedulers: %w", err)
	}

	return nil
}
