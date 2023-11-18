package receiver

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	receiverpb "github.com/c4t-but-s4d/fastad/pkg/proto/receiver"
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
