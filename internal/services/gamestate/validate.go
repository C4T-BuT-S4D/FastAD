package gamestate

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	gspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/game_state"
)

func (s *Service) validateUpdateRequest(req *gspb.UpdateRequest) error {
	if req.StartTime == nil {
		return status.Error(codes.InvalidArgument, "start_time required")
	}
	if end := req.EndTime; end != nil && end.AsTime().Before(req.StartTime.AsTime()) {
		return status.Error(codes.InvalidArgument, "end_time is before start time")
	}
	if req.RoundDuration.AsDuration() == 0 {
		return status.Error(codes.InvalidArgument, "round_duration required")
	}
	if req.FlagLifetimeRounds == 0 {
		return status.Error(codes.InvalidArgument, "flag_lifetime_rounds required")
	}
	if req.Mode == gspb.GameMode_GAME_MODE_UNSPECIFIED {
		return status.Error(codes.InvalidArgument, "mode required")
	}
	return nil
}
