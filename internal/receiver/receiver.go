package receiver

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/uptrace/bun"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/c4t-but-s4d/fastad/internal/centutil"
	"github.com/c4t-but-s4d/fastad/internal/models"
	"github.com/c4t-but-s4d/fastad/pkg/clients/gamestate"
	"github.com/c4t-but-s4d/fastad/pkg/clients/services"
	"github.com/c4t-but-s4d/fastad/pkg/clients/teams"
	gspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/game_state"
	servicespb "github.com/c4t-but-s4d/fastad/pkg/proto/data/services"
	teamspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/teams"
	receiverpb "github.com/c4t-but-s4d/fastad/pkg/proto/receiver"
)

const maxFlagsInRequest = 100

const (
	invalidFlagMessage    = "invalid flag"
	oldFlagMessage        = "flag is too old"
	ownFlagMessage        = "flag is your own"
	serviceInvalidMessage = "service is invalid or disabled"
	duplicateFlagMessage  = "flag already submitted"
	notReadyFlagMessage   = "flag is not ready"
)

type Service struct {
	receiverpb.UnimplementedReceiverServiceServer

	db              *bun.DB
	teamsClient     *teams.Client
	servicesClient  *services.Client
	gameStateClient *gamestate.Client
	producer        centutil.Producer
	metrics         *Metrics

	stateMu sync.Mutex
	state   *State
}

func New(
	db *bun.DB,
	teamsClient *teams.Client,
	servicesClient *services.Client,
	gameStateClient *gamestate.Client,
	producer centutil.Producer,
	metrics *Metrics,
) *Service {
	return &Service{
		db:              db,
		teamsClient:     teamsClient,
		servicesClient:  servicesClient,
		gameStateClient: gameStateClient,
		state:           NewState(),
		producer:        producer,
		metrics:         metrics,
	}
}

func (s *Service) SubmitFlags(ctx context.Context, req *receiverpb.SubmitFlagsRequest) (*receiverpb.SubmitFlagsResponse, error) {
	start := time.Now()
	zap.L().Debug("Receiver/SubmitFlags", zap.Any("request", req))

	if len(req.GetFlags()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "no flags")
	}
	if len(req.GetFlags()) > maxFlagsInRequest {
		return nil, status.Errorf(codes.InvalidArgument, "too many flags (max %d)", maxFlagsInRequest)
	}
	if req.GetTeamToken() == "" {
		return nil, status.Error(codes.InvalidArgument, "team token is required")
	}

	gameState, err := s.gameStateClient.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetching game state: %w", err)
	}

	switch gameState.GetStatus() {
	case gspb.GameStatus_GAME_STATUS_PAUSED:
		return nil, status.Error(codes.Unavailable, "game is paused")
	case gspb.GameStatus_GAME_STATUS_FINISHED:
		return nil, status.Error(codes.FailedPrecondition, "game is finished")
	case gspb.GameStatus_GAME_STATUS_NOT_STARTED:
		return nil, status.Error(codes.FailedPrecondition, "game has not started")
	}
	// Also check end time even if status is RUNNING
	if gameState.GetEndTime() != nil && time.Now().After(gameState.GetEndTime().AsTime()) {
		return nil, status.Error(codes.FailedPrecondition, "game is finished")
	}

	serviceList, err := s.servicesClient.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetching services: %w", err)
	}
	serviceByID := lo.KeyBy(serviceList, func(srv *servicespb.Service) int {
		return int(srv.GetId())
	})

	teamList, err := s.teamsClient.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetching teams: %w", err)
	}
	teamByID := lo.KeyBy(teamList, func(t *teamspb.Team) int {
		return int(t.GetId())
	})

	attacker, ok := lo.Find(teamList, func(t *teamspb.Team) bool {
		return t.GetToken() == req.GetTeamToken()
	})
	if !ok {
		return nil, status.Error(codes.NotFound, "team not found")
	}

	attacksRequestID := uuid.NewString()

	var addedAttacks []*models.Attack

	uniqueFlags := lo.Uniq(req.GetFlags())
	resp := &receiverpb.SubmitFlagsResponse{}
	if err := s.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		var flagModels []*models.Flag
		if err := tx.
			NewSelect().
			Model(&models.Flag{}).
			ExcludeColumn("public", "private").
			Where("flag IN (?)", bun.In(uniqueFlags)).
			Scan(ctx, &flagModels); err != nil {
			return fmt.Errorf("selecting flags: %w", err)
		}

		flagByID := lo.KeyBy(flagModels, func(flag *models.Flag) int {
			return flag.ID
		})

		haveFlags := lo.Map(flagModels, func(item *models.Flag, _ int) string {
			return item.Flag
		})
		missingFlags, _ := lo.Difference(uniqueFlags, haveFlags)
		for _, flag := range missingFlags {
			resp.Responses = append(resp.Responses, &receiverpb.FlagResponse{
				Flag:    flag,
				Verdict: receiverpb.FlagResponse_VERDICT_INVALID,
				Message: invalidFlagMessage,
			})
		}

		attacksToAdd := make([]*models.Attack, 0, len(flagModels))
		attackFlagIDs := make([]int, 0, len(flagModels))
		for _, flag := range flagModels {
			baseResponse := &receiverpb.FlagResponse{
				Flag:      flag.Flag,
				ServiceId: int64(flag.ServiceID),
				VictimId:  int64(flag.TeamID),
			}

			if service, ok := serviceByID[flag.ServiceID]; !ok || service.GetDisabled() {
				baseResponse.Verdict = receiverpb.FlagResponse_VERDICT_INVALID
				baseResponse.Message = serviceInvalidMessage
				resp.Responses = append(resp.Responses, baseResponse)
				continue
			}

			if flag.TeamID == int(attacker.GetId()) {
				baseResponse.Verdict = receiverpb.FlagResponse_VERDICT_OWN
				baseResponse.Message = ownFlagMessage
				resp.Responses = append(resp.Responses, baseResponse)
				continue
			}

			if gameState.GetRunningRound()-flag.Round > gameState.GetFlagLifetimeRounds() {
				baseResponse.Verdict = receiverpb.FlagResponse_VERDICT_OLD
				baseResponse.Message = oldFlagMessage
				resp.Responses = append(resp.Responses, baseResponse)
				continue
			}

			// Allow players to resubmit the flag later when the put is finished.
			if !flag.PutFinished {
				baseResponse.Verdict = receiverpb.FlagResponse_VERDICT_FLAG_NOT_READY
				baseResponse.Message = notReadyFlagMessage
				resp.Responses = append(resp.Responses, baseResponse)
				continue
			}

			attacksToAdd = append(attacksToAdd, &models.Attack{
				ServiceID:  flag.ServiceID,
				AttackerID: int(attacker.GetId()),
				VictimID:   flag.TeamID,
				FlagID:     flag.ID,
				RequestID:  attacksRequestID,
			})
			attackFlagIDs = append(attackFlagIDs, flag.ID)
		}

		if len(attacksToAdd) == 0 {
			return nil
		}

		if _, err := tx.
			NewInsert().
			Model(&attacksToAdd).
			On("CONFLICT (attacker_id, flag_id) DO NOTHING").
			Returning("id, flag_id").
			Exec(ctx); err != nil {
			return fmt.Errorf("inserting attacks: %w", err)
		}

		insertedAttackFlagIDs := lo.Map(attacksToAdd, func(attack *models.Attack, _ int) int {
			return attack.FlagID
		})

		duplicateFlagIDs, _ := lo.Difference(attackFlagIDs, insertedAttackFlagIDs)
		for _, flagID := range duplicateFlagIDs {
			flag := flagByID[flagID]
			resp.Responses = append(resp.Responses, &receiverpb.FlagResponse{
				Flag:      flag.Flag,
				Verdict:   receiverpb.FlagResponse_VERDICT_DUPLICATE,
				ServiceId: int64(flag.ServiceID),
				VictimId:  int64(flag.TeamID),
				Message:   duplicateFlagMessage,
			})
		}

		if len(attacksToAdd) == 0 {
			return nil
		}

		s.stateMu.Lock()
		defer s.stateMu.Unlock()

		savedState := s.state.Clone()
		shouldRollbackState := true
		defer func() {
			if shouldRollbackState {
				s.state = savedState
			}
		}()

		for i, attack := range attacksToAdd {
			if err := s.state.ProcessAttack(gameState, serviceByID[attack.ServiceID], attack); err != nil {
				return fmt.Errorf("applying attack #%d: %w", i, err)
			}

			resp.Responses = append(resp.Responses, &receiverpb.FlagResponse{
				Flag:    flagByID[attack.FlagID].Flag,
				Verdict: receiverpb.FlagResponse_VERDICT_ACCEPTED,
				Message: fmt.Sprintf(
					"attacked team %d and gained %.4f points",
					attack.VictimID,
					attack.AttackerDelta,
				),
				ServiceId:     int64(attack.ServiceID),
				VictimId:      int64(attack.VictimID),
				AttackerDelta: attack.AttackerDelta,
				VictimDelta:   attack.VictimDelta,
			})
		}

		if _, err := tx.
			NewUpdate().
			Model(&attacksToAdd).
			Column("attacker_delta", "victim_delta").
			Bulk().
			Exec(ctx); err != nil {
			return fmt.Errorf("updating attacks' deltas: %w", err)
		}

		// Intentionally ignoring the potential issue of a missing state rollback in
		// the improbable case all queries in tx finish but the tx is reverted afterward.
		shouldRollbackState = false
		addedAttacks = attacksToAdd

		return nil
	}); err != nil {
		return nil, fmt.Errorf("in transaction: %w", err)
	}

	if len(addedAttacks) > 0 {
		notification := &receiverpb.AttackNotification_Batch{
			Attacks: lo.Map(addedAttacks, func(attack *models.Attack, _ int) *receiverpb.AttackNotification {
				return attack.ToNotificationProto()
			}),
		}
		if err := s.producer.PublishProto(ctx, notification); err != nil {
			zap.L().Error("publishing attack notification", zap.Error(err))
		}
	}

	s.observeMetrics(resp, attacker, teamByID, serviceByID, addedAttacks, time.Since(start).Seconds())

	return resp, nil
}

func (s *Service) observeMetrics(
	resp *receiverpb.SubmitFlagsResponse,
	attacker *teamspb.Team,
	teamByID map[int]*teamspb.Team,
	serviceByID map[int]*servicespb.Service,
	addedAttacks []*models.Attack,
	durationSeconds float64,
) {
	attackerIDStr := strconv.FormatInt(attacker.GetId(), 10)
	attackerName := attacker.GetName()
	s.metrics.ObserveSubmission(attackerIDStr, attackerName, len(resp.GetResponses()))

	hasAccepted := false
	for _, r := range resp.GetResponses() {
		serviceID := strconv.FormatInt(r.GetServiceId(), 10)
		serviceName := ""
		if svc, ok := serviceByID[int(r.GetServiceId())]; ok {
			serviceName = svc.GetName()
		}

		victimID := strconv.FormatInt(r.GetVictimId(), 10)
		victimName := ""
		if team, ok := teamByID[int(r.GetVictimId())]; ok {
			victimName = team.GetName()
		}

		s.metrics.ObserveVerdict(r.GetVerdict(), serviceID, serviceName)
		s.metrics.ObserveVerdictDetailed(attackerIDStr, attackerName, victimID, victimName, serviceID, serviceName, r.GetVerdict())
		if r.GetVerdict() == receiverpb.FlagResponse_VERDICT_ACCEPTED {
			hasAccepted = true
		}
	}

	for _, attack := range addedAttacks {
		victimIDStr := strconv.Itoa(attack.VictimID)
		victimName := ""
		if team, ok := teamByID[attack.VictimID]; ok {
			victimName = team.GetName()
		}

		serviceIDStr := strconv.Itoa(attack.ServiceID)
		serviceName := ""
		if svc, ok := serviceByID[attack.ServiceID]; ok {
			serviceName = svc.GetName()
		}

		s.metrics.ObserveAttack(
			attackerIDStr,
			attackerName,
			victimIDStr,
			victimName,
			serviceIDStr,
			serviceName,
			attack.AttackerDelta,
		)
	}

	s.metrics.ObserveProcessingTime(durationSeconds, hasAccepted)
}

func (s *Service) GetState(context.Context, *receiverpb.GetStateRequest) (*receiverpb.GetStateResponse, error) {
	s.stateMu.Lock()
	defer s.stateMu.Unlock()

	return &receiverpb.GetStateResponse{
		State: s.state.ToProto(),
	}, nil
}

func (s *Service) RestoreState(ctx context.Context) error {
	serviceList, err := s.servicesClient.List(ctx)
	if err != nil {
		return fmt.Errorf("fetching services: %w", err)
	}

	servicesByID := lo.KeyBy(serviceList, func(service *servicespb.Service) int {
		return int(service.GetId())
	})

	start := time.Now()

	if err := s.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		attackCount, err := tx.NewSelect().Model(&models.Attack{}).Count(ctx)
		if err != nil {
			return fmt.Errorf("counting attacks: %w", err)
		}

		zap.L().Info("restoring state from attacks", zap.Int("attack_count", attackCount))

		lastID := -1
		for {
			const batchSize = 1000
			var batch []*models.Attack
			if err := tx.
				NewSelect().
				Model(&models.Attack{}).
				Where("id > ?", lastID).
				Order("id").
				Limit(batchSize).
				Scan(ctx, &batch); err != nil {
				return fmt.Errorf("fetching attacks batch: %w", err)
			}
			if len(batch) == 0 {
				break
			}
			lastID = batch[len(batch)-1].ID

			zap.L().Info("applying batch of attacks", zap.Int("batch_size", len(batch)))
			if err := s.state.ApplyRaw(servicesByID, batch...); err != nil {
				return fmt.Errorf("applying batch of %d attacks to state: %w", len(batch), err)
			}

			if len(batch) < batchSize {
				break
			}
		}

		return nil
	}); err != nil {
		return fmt.Errorf("in tx: %w", err)
	}

	zap.L().Info("state restored", zap.Duration("duration", time.Since(start)))

	return nil
}
