package gamestate

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	gspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/game_state"
)

func (s *Service) validateUpdateRequest(req *gspb.UpdateRequest) error {
	if req.GetStartTime() == nil {
		return status.Error(codes.InvalidArgument, "start_time required")
	}
	if end := req.GetEndTime(); end != nil && end.AsTime().Before(req.GetStartTime().AsTime()) {
		return status.Error(codes.InvalidArgument, "end_time is before start time")
	}
	if req.GetRoundDuration().AsDuration() == 0 {
		return status.Error(codes.InvalidArgument, "round_duration required")
	}
	if req.GetFlagLifetimeRounds() == 0 {
		return status.Error(codes.InvalidArgument, "flag_lifetime_rounds required")
	}
	if req.GetMode() == gspb.GameMode_GAME_MODE_UNSPECIFIED {
		req.Mode = gspb.GameMode_GAME_MODE_CLASSIC
	}
	return nil
}

func (s *Service) validateUpdateRoundRequest(req *gspb.UpdateRoundRequest) error {
	if req.GetRunningRound() == 0 {
		return status.Error(codes.InvalidArgument, "running_round required")
	}
	if req.GetRunningRoundStart() == nil {
		return status.Error(codes.InvalidArgument, "running_round_start required")
	}
	return nil
}
