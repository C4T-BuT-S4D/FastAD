package httpext

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"google.golang.org/grpc/status"
)

type Error struct {
	code    int
	message string
}

func NewError(code int, msg string) *Error {
	return &Error{
		code:    code,
		message: msg,
	}
}

func NewErrorf(code int, format string, a ...any) *Error {
	return NewError(code, fmt.Sprintf(format, a...))
}

func NewErrorFromStatus(err error, msg string) *Error {
	if st, ok := status.FromError(err); ok {
		return NewError(StatusFromGRPC(st), fmt.Sprintf("%s: %v", msg, st.Message()))
	}

	return &Error{
		code:    http.StatusInternalServerError,
		message: fmt.Sprintf("%s: %v", msg, err),
	}
}

func NewErrorFromStatusf(err error, format string, a ...any) *Error {
	return NewErrorFromStatus(err, fmt.Sprintf(format, a...))
}

func (e *Error) Error() string {
	return e.message
}

func ErrorHandler() echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		logger := ContextLogger(c)
		c.Response().Header().Set("X-Request-ID", RequestID(c))

		var (
			e  *Error
			he *echo.HTTPError
		)

		switch {
		case errors.As(err, &e):
			if e.code/100 == 5 {
				logger.Warn(
					"internal error",
					zap.Error(err),
					zap.String("message", e.message),
					zap.Int("status", e.code),
				)
			} else {
				logger.Debug(
					"client error",
					zap.Error(err),
					zap.String("message", e.message),
					zap.Int("status", e.code),
				)
			}

			if err := c.JSON(e.code, map[string]string{"error": e.message}); err != nil {
				c.Logger().Error(err)
			}

		case errors.As(err, &he):
			if he.Internal != nil {
				var herr *echo.HTTPError
				if errors.As(he.Internal, &herr) {
					he = herr
				}
			}

			f := logger.Debug
			if he.Code/100 == 5 {
				f = logger.Warn
			}
			f(
				"unexpected error",
				zap.Error(err),
				zap.Any("message", he.Message),
				zap.Int("status", he.Code),
			)
			if c.Request().Method == http.MethodHead {
				if err := c.NoContent(he.Code); err != nil {
					c.Logger().Error(err)
				}
			} else {
				if msg, ok := he.Message.(string); ok {
					if err := c.JSON(he.Code, map[string]string{"error": msg}); err != nil {
						c.Logger().Error(err)
					}
				} else {
					if err := c.JSON(he.Code, map[string]string{"error": http.StatusText(he.Code)}); err != nil {
						c.Logger().Error(err)
					}
				}
			}

		default:
			logger.Error("unexpected error", zap.Error(err))
			if err := c.JSON(http.StatusInternalServerError, map[string]string{"error": "unknown error"}); err != nil {
				c.Logger().Error(err)
			}
		}
	}
}
