package httpext

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

const RequestIDContextKey = "request_id"

func RequestIDMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Set(RequestIDContextKey, uuid.NewString())
		return next(c)
	}
}

func ContextLogger(c echo.Context) *zap.Logger {
	return zap.L().With(zap.String(RequestIDContextKey, RequestID(c)))
}
