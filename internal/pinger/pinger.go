package pinger

import (
	"context"

	pingerpb "fastad/pkg/proto/pinger"
)

func New() *Service {
	return &Service{}
}

type Service struct {
	*pingerpb.UnimplementedPingerServiceServer
}

func (*Service) Ping(context.Context, *pingerpb.PingRequest) (*pingerpb.PingResponse, error) {
	return &pingerpb.PingResponse{}, nil
}
