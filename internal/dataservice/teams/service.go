package teams

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/c4t-but-s4d/fastad/internal/models"
	"github.com/c4t-but-s4d/fastad/internal/version"
	teamspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/teams"
)

type Service struct {
	teamspb.UnimplementedTeamsServiceServer

	controller *Controller
}

func NewService(controller *Controller) *Service {
	return &Service{
		controller: controller,
	}
}

func (s *Service) List(ctx context.Context, req *teamspb.ListRequest) (*teamspb.ListResponse, error) {
	zap.L().Debug("TeamsService/List", zap.Any("request", req))

	gotVersion, err := s.controller.Versions.Get(ctx, VersionKey)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "getting version: %v", err)
	}

	requestedVersion := int(req.GetVersion().GetVersion())
	if requestedVersion > gotVersion {
		return nil, status.Errorf(codes.FailedPrecondition, "requested version is greater than current")
	}
	if requestedVersion == gotVersion {
		return &teamspb.ListResponse{Version: version.NewVersionProto(gotVersion)}, nil
	}

	teams, err := s.controller.List(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "getting teams: %v", err)
	}

	return &teamspb.ListResponse{
		Teams: lo.Map(teams, func(team *models.Team, _ int) *teamspb.Team {
			return team.ToProto()
		}),
		Version: version.NewVersionProto(gotVersion),
	}, nil
}

func (s *Service) CreateBatch(ctx context.Context, req *teamspb.CreateBatchRequest) (*teamspb.CreateBatchResponse, error) {
	zap.L().Debug("TeamsService/CreateBatch", zap.Any("request", req))

	if err := s.validateCreateBatchRequest(req); err != nil {
		return nil, fmt.Errorf("validating request: %w", err)
	}

	teams := lo.Map(req.Teams, func(team *teamspb.Team, _ int) *models.Team {
		teamModel := models.NewTeamFromProto(team)
		teamModel.ID = 0
		teamModel.Token = uuid.NewString()
		return teamModel
	})

	if err := s.controller.CreateBatch(ctx, teams); err != nil {
		return nil, status.Errorf(codes.Internal, "creating teams: %v", err)
	}

	result := lo.Map(teams, func(team *models.Team, _ int) *teamspb.Team {
		return team.ToProto()
	})
	return &teamspb.CreateBatchResponse{
		Teams: result,
	}, nil
}
