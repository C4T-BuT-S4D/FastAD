package teams

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	teamspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/teams"
)

func (s *Service) validateCreateBatchRequest(req *teamspb.CreateBatchRequest) error {
	if len(req.Teams) == 0 {
		return status.Error(codes.InvalidArgument, "teams required")
	}
	for i, team := range req.Teams {
		if team.Name == "" {
			return status.Errorf(codes.InvalidArgument, "teams.%d: name required", i)
		}
		if team.Address == "" {
			return status.Errorf(codes.InvalidArgument, "teams.%d: address required", i)
		}
	}
	return nil
}
