package grpcext

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	// Enable gzip compression.
	_ "google.golang.org/grpc/encoding/gzip"
)

func NewServer() *grpc.Server {
	s := grpc.NewServer()
	reflection.Register(s)
	return s
}
