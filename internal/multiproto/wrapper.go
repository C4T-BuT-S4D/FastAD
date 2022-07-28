package multiproto

import (
	"net/http"
	"strings"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
)

func OriginAllowed(_ string) bool {
	return true
}

func WebSocketOriginAllowed(_ *http.Request) bool {
	return true
}

func NewHandler(grpcServer *grpc.Server) http.Handler {
	return h2c.NewHandler(
		&Handler{
			grpcServer: grpcServer,
			webServer: grpcweb.WrapServer(
				grpcServer,
				grpcweb.WithWebsockets(true),
				grpcweb.WithOriginFunc(OriginAllowed),
				grpcweb.WithWebsocketOriginFunc(WebSocketOriginAllowed),
			),
			httpHandler: http.DefaultServeMux,
		},
		&http2.Server{},
	)
}

type Handler struct {
	grpcServer  *grpc.Server
	webServer   *grpcweb.WrappedGrpcServer
	httpHandler http.Handler
}

func (s *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc"):
		s.grpcServer.ServeHTTP(w, r)
	case s.webServer.IsGrpcWebRequest(r) || s.webServer.IsGrpcWebSocketRequest(r):
		s.webServer.ServeHTTP(w, r)
	default:
		s.httpHandler.ServeHTTP(w, r)
	}
}
