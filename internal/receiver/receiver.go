package receiver

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/uptrace/bun"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/c4t-but-s4d/fastad/internal/clients/gamestate"
	"github.com/c4t-but-s4d/fastad/internal/clients/services"
	"github.com/c4t-but-s4d/fastad/internal/clients/teams"
	"github.com/c4t-but-s4d/fastad/internal/models"
	receiverpb "github.com/c4t-but-s4d/fastad/pkg/proto/receiver"
)

const maxFlagsInRequest = 100

const (
	invalidFlagMessage    = "invalid flag"
	oldFlagMessage        = "flag is too old"
	ownFlagMessage        = "flag is your own"
	serviceInvalidMessage = "service is invalid or disabled"
	duplicateFlagMessage  = "flag already submitted"
)

func New(
	db *bun.DB,
	teamsClient *teams.Client,
	servicesClient *services.Client,
	gameStateClient *gamestate.Client,
) *Service {
	return &Service{
		db:              db,
		teamsClient:     teamsClient,
		servicesClient:  servicesClient,
		gameStateClient: gameStateClient,
		state:           NewState(),
	}
}

type Service struct {
	receiverpb.UnimplementedReceiverServiceServer

	db              *bun.DB
	teamsClient     *teams.Client
	servicesClient  *services.Client
	gameStateClient *gamestate.Client

	stateMu sync.Mutex
	state   *State
}

func (s *Service) SubmitFlags(ctx context.Context, req *receiverpb.SubmitFlagsRequest) (*receiverpb.SubmitFlagsResponse, error) {
	zap.L().Debug("Receiver/SubmitFlags", zap.Any("request", req))

	if len(req.Flags) == 0 {
		return nil, status.Error(codes.InvalidArgument, "no flags")
	}

	if len(req.Flags) > maxFlagsInRequest {
		return nil, status.Errorf(codes.InvalidArgument, "too many flags (max %d)", maxFlagsInRequest)
	}

	tokens := metadata.ValueFromIncomingContext(ctx, "team_token")
	if len(tokens) != 1 {
		return nil, status.Error(codes.InvalidArgument, "team token required")
	}

	gameState, err := s.gameStateClient.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetching game state: %w", err)
	}

	serviceList, err := s.servicesClient.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetching services: %w", err)
	}
	serviceByID := lo.KeyBy(serviceList, func(srv *models.Service) int {
		return srv.ID
	})

	attacker, err := s.teamsClient.GetByToken(ctx, tokens[0])
	if err != nil {
		return nil, fmt.Errorf("fetching attacker: %w", err)
	}

	if gameState.Paused {
		return nil, status.Error(codes.Unavailable, "game is paused")
	}

	attacksRequestID := uuid.NewString()

	uniqueFlags := lo.Uniq(req.Flags)
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

			if service, ok := serviceByID[flag.ServiceID]; !ok || service.Disabled {
				baseResponse.Verdict = receiverpb.FlagResponse_VERDICT_INVALID
				baseResponse.Message = serviceInvalidMessage
				resp.Responses = append(resp.Responses, baseResponse)
				continue
			}

			if flag.TeamID == attacker.ID {
				baseResponse.Verdict = receiverpb.FlagResponse_VERDICT_OWN
				baseResponse.Message = ownFlagMessage
				resp.Responses = append(resp.Responses, baseResponse)
				continue
			}

			// TODO: check flag lifetime (skipped for now for easier manual tests).
			// if gameState.RunningRound-flag.Round > gameState.FlagLifetimeRounds {
			// 	baseResponse.Verdict = receiverpb.FlagResponse_VERDICT_OLD
			// 	baseResponse.Message = oldFlagMessage
			// 	resp.Responses = append(resp.Responses, baseResponse)
			// 	continue
			// }

			attacksToAdd = append(attacksToAdd, &models.Attack{
				ServiceID:  flag.ServiceID,
				AttackerID: attacker.ID,
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

		return nil
	}); err != nil {
		return nil, fmt.Errorf("in transaction: %w", err)
	}

	return resp, nil
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

	servicesByID := lo.KeyBy(serviceList, func(service *models.Service) int {
		return service.ID
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
