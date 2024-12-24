package httpext

import (
	"context"

	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/metadata"
)

func ContextFromEcho(c echo.Context) context.Context {
	return metadata.NewIncomingContext(
		c.Request().Context(),
		metadata.Pairs(RequestIDContextKey, RequestID(c)),
	)
}

func RequestID(c echo.Context) string {
	if requestID := c.Get(RequestIDContextKey); requestID != nil {
		if requestIDStr, ok := requestID.(string); ok {
			return requestIDStr
		}
	}
	panic("request id unset")
}
