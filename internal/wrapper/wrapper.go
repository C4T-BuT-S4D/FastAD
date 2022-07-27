package wrapper

import (
	"net/http"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"
)

func OriginAllowed(_ string) bool {
	return true
}

func WebSocketOriginAllowed(_ *http.Request) bool {
	return true
}

func NewServer(grpcServer *grpc.Server) *http.Server {
	httpServer := &http.Server{
		Handler: grpcweb.WrapServer(
			grpcServer,
			grpcweb.WithWebsockets(true),
			grpcweb.WithOriginFunc(OriginAllowed),
			grpcweb.WithWebsocketOriginFunc(WebSocketOriginAllowed),
		),
	}

	return httpServer
}
