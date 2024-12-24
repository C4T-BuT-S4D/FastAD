package gamestate

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/c4t-but-s4d/fastad/internal/version"
	gspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/game_state"
)

type Service struct {
	gspb.UnimplementedGameStateServiceServer

	controller *Controller
}

func NewService(controller *Controller) *Service {
	return &Service{controller: controller}
}

func (s *Service) Get(ctx context.Context, req *gspb.GetRequest) (*gspb.GetResponse, error) {
	zap.L().Debug("GameStateService/Get", zap.Any("request", req))

	gotVersion, err := s.controller.Versions.Get(ctx, VersionKey)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "getting version: %v", err)
	}

	requestedVersion := int(req.GetVersion().GetVersion())
	if requestedVersion > gotVersion {
		return nil, status.Errorf(codes.FailedPrecondition, "requested version is greater than current")
	}
	if requestedVersion == gotVersion {
		return &gspb.GetResponse{Version: version.NewVersionProto(gotVersion)}, nil
	}

	state, err := s.controller.Get(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "getting game state: %v", err)
	}

	return &gspb.GetResponse{
		GameState: state.ToProto(),
		Version:   version.NewVersionProto(gotVersion),
	}, nil
}

func (s *Service) Update(ctx context.Context, req *gspb.UpdateRequest) (*gspb.UpdateResponse, error) {
	zap.L().Debug("GameStateService/Update", zap.Any("request", req))

	// FiXME: security.

	if err := s.validateUpdateRequest(req); err != nil {
		return nil, fmt.Errorf("validating request: %w", err)
	}

	gs, newVersion, err := s.controller.Update(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("updating game state: %w", err)
	}

	return &gspb.UpdateResponse{
		GameState: gs.ToProto(),
		Version:   version.NewVersionProto(newVersion),
	}, nil
}

func (s *Service) UpdateRound(ctx context.Context, req *gspb.UpdateRoundRequest) (*gspb.UpdateRoundResponse, error) {
	zap.L().Debug("GameStateService/UpdateRound", zap.Any("request", req))

	// FiXME: security.

	if err := s.validateUpdateRoundRequest(req); err != nil {
		return nil, fmt.Errorf("validating request: %w", err)
	}

	gs, newVersion, err := s.controller.UpdateRound(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("updating round: %w", err)
	}

	return &gspb.UpdateRoundResponse{
		GameState: gs.ToProto(),
		Version:   version.NewVersionProto(newVersion),
	}, nil
}
