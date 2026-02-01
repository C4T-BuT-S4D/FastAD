package services

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	checkerpb "github.com/c4t-but-s4d/fastad/pkg/proto/checker"
	servicespb "github.com/c4t-but-s4d/fastad/pkg/proto/data/services"
)

func (s *Service) validateCreateBatch(req *servicespb.CreateBatchRequest) error {
	if len(req.GetServices()) == 0 {
		return status.Errorf(codes.InvalidArgument, "services required")
	}
	for i, service := range req.GetServices() {
		if service.GetName() == "" {
			return status.Errorf(codes.InvalidArgument, "services.%d: name required", i)
		}
		if service.GetChecker() == nil {
			return status.Errorf(codes.InvalidArgument, "services.%d: checker required", i)
		}
		checker := service.GetChecker()
		if checker.GetType() == checkerpb.Type_TYPE_UNSPECIFIED {
			return status.Errorf(codes.InvalidArgument, "services.%d.checker: type required", i)
		}
		if checker.GetPath() == "" {
			return status.Errorf(codes.InvalidArgument, "services.%d.checker: path required", i)
		}
		if checker.GetDefaultTimeout().AsDuration() == 0 {
			return status.Errorf(codes.InvalidArgument, "services.%d.checker: default_timeout required", i)
		}
		for j, action := range checker.GetActions() {
			if action.GetAction() == checkerpb.Action_ACTION_UNSPECIFIED {
				return status.Errorf(codes.InvalidArgument, "services.%d.checker.actions.%d: action required", i, j)
			}
		}
	}
	return nil
}

func (s *Service) validateUpdateRequest(req *servicespb.UpdateRequest) error {
	if req.GetId() == 0 {
		return status.Error(codes.InvalidArgument, "id required")
	}
	return nil
}
