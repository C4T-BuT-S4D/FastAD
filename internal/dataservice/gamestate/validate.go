package gamestate

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	gspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/game_state"
)

func (s *Service) validateCreateRequest(req *gspb.CreateRequest) error {
	if req.GetGameState() == nil {
		return status.Error(codes.InvalidArgument, "game_state required")
	}
	gs := req.GetGameState()
	if gs.GetStartTime() == nil {
		return status.Error(codes.InvalidArgument, "start_time required")
	}
	if end := gs.GetEndTime(); end != nil && end.AsTime().Before(gs.GetStartTime().AsTime()) {
		return status.Error(codes.InvalidArgument, "end_time is before start time")
	}
	if gs.GetRoundDuration().AsDuration() == 0 {
		return status.Error(codes.InvalidArgument, "round_duration required")
	}
	if gs.GetFlagLifetimeRounds() == 0 {
		return status.Error(codes.InvalidArgument, "flag_lifetime_rounds required")
	}
	return nil
}

func (s *Service) validateUpdateRequest(req *gspb.UpdateRequest) error {
	startTime := req.GetStartTime()
	endTime := req.GetEndTime()
	if startTime != nil && endTime != nil && !endTime.AsTime().IsZero() && endTime.AsTime().Before(startTime.AsTime()) {
		return status.Error(codes.InvalidArgument, "end_time is before start time")
	}

	if rd := req.GetRoundDuration(); rd != nil && rd.AsDuration() == 0 {
		return status.Error(codes.InvalidArgument, "round_duration must be non-zero")
	}
	if req.FlagLifetimeRounds != nil && req.GetFlagLifetimeRounds() == 0 {
		return status.Error(codes.InvalidArgument, "flag_lifetime_rounds must be non-zero")
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
