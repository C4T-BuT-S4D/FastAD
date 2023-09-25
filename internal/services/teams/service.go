package teams

import (
	"context"

	"github.com/c4t-but-s4d/fastad/internal/models"
	teamspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/teams"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	lastUpdate, err := s.controller.LastUpdate(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "getting last update time: %v", err)
	}

	if req.LastUpdate != 0 {
		if req.LastUpdate > lastUpdate {
			return nil, status.Errorf(codes.InvalidArgument, "last update time is in the future")
		}
		if req.LastUpdate == lastUpdate {
			return &teamspb.ListResponse{LastUpdate: lastUpdate}, nil
		}
	}

	teams, err := s.controller.List(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "getting teams: %v", err)
	}

	return &teamspb.ListResponse{
		LastUpdate: lastUpdate,
		Teams: lo.Map(teams, func(team *models.Team, _ int) *teamspb.Team {
			return team.ToProto()
		}),
	}, nil
}

func (s *Service) CreateBatch(ctx context.Context, req *teamspb.CreateBatchRequest) (*teamspb.CreateBatchResponse, error) {
	// FIXME: check admin rights.

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
