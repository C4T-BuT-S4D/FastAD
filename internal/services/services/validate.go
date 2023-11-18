package services

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	checkerpb "github.com/c4t-but-s4d/fastad/pkg/proto/checker"
	servicespb "github.com/c4t-but-s4d/fastad/pkg/proto/data/services"
)

func (s *Service) validateCreateBatch(req *servicespb.CreateBatchRequest) error {
	if len(req.Services) == 0 {
		return status.Errorf(codes.InvalidArgument, "services required")
	}
	for i, service := range req.Services {
		if service.Name == "" {
			return status.Errorf(codes.InvalidArgument, "services.%d: name required", i)
		}
		if service.Checker == nil {
			return status.Errorf(codes.InvalidArgument, "services.%d: checker required", i)
		}
		checker := service.Checker
		if checker.Type == checkerpb.Type_TYPE_UNSPECIFIED {
			return status.Errorf(codes.InvalidArgument, "services.%d.checker: type required", i)
		}
		if checker.Path == "" {
			return status.Errorf(codes.InvalidArgument, "services.%d.checker: path required", i)
		}
		if checker.DefaultTimeout.AsDuration() == 0 {
			return status.Errorf(codes.InvalidArgument, "services.%d.checker: default_timeout required", i)
		}
		for j, action := range checker.Actions {
			if action.Action == checkerpb.Action_ACTION_UNSPECIFIED {
				return status.Errorf(codes.InvalidArgument, "services.%d.checker.actions.%d: action required", i, j)
			}
		}
	}
	return nil
}
