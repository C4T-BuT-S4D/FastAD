package grpctools

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "google.golang.org/grpc/encoding/gzip"
)

func NewServer() *grpc.Server {
	s := grpc.NewServer()
	reflection.Register(s)
	return s
}
