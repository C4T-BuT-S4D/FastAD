package httpext

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const RequestIDContextKey = "request_id"

func RequestIDMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()
			rid := req.Header.Get(echo.HeaderXRequestID)
			if rid == "" {
				rid = uuid.NewString()
			}
			res.Header().Set(echo.HeaderXRequestID, rid)
			c.Set(RequestIDContextKey, rid)
			return next(c)
		}
	}
}
