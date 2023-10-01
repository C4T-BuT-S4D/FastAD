package gamestate

import (
	"context"
	"fmt"

	"github.com/c4t-but-s4d/fastad/internal/version"
	gspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/game_state"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	gspb.UnimplementedGameStateServiceServer

	controller *Controller
}

func NewService(controller *Controller) *Service {
	return &Service{controller: controller}
}

func (s *Service) Get(ctx context.Context, req *gspb.GetRequest) (*gspb.GetResponse, error) {
	logrus.Debugf("GameStateService/Get: %v", req)

	gotVersion, err := s.controller.Versions.Get(ctx, VersionKey)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "getting version: %v", err)
	}

	requestedVersion := req.GetVersion().GetVersion()
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
	logrus.Debugf("GameStateService/Update: %v", req)

	if err := s.sanitizeUpdateRequest(req); err != nil {
		return nil, fmt.Errorf("sanitizing request: %w", err)
	}

	gs, err := s.controller.Update(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("updating game state: %w", err)
	}

	return &gspb.UpdateResponse{
		GameState: gs.ToProto(),
	}, nil
}
