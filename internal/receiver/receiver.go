package receiver

import (
	"context"

	receiverpb "FastAD/pkg/proto/receiver"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func New() *Service {
	return &Service{}
}

type Service struct {
	*receiverpb.UnimplementedReceiverServiceServer
}

func (s *Service) SubmitFlags(context.Context, *receiverpb.SubmitFlagsRequest) (*receiverpb.SubmitFlagsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented yet")
}
