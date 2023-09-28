package services

import (
	servicespb "github.com/c4t-but-s4d/fastad/pkg/proto/data/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) sanitizeCreateBatch(req *servicespb.CreateBatchRequest) error {
	if len(req.Services) == 0 {
		return status.Errorf(codes.InvalidArgument, "empty services list")
	}
	for i, service := range req.Services {
		if service.Checker == nil {
			return status.Errorf(codes.InvalidArgument, "services.%d: checker is required", i)
		}
	}
	return nil
}
