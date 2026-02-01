package httpext

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func HealthHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	}
}
