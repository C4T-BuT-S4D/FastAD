package api

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/c4t-but-s4d/fastad/pkg/httpext"
)

func (s *Service) HandleGetGameState() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := httpext.ContextFromEcho(c)

		gs, err := s.gameStateClient.Get(ctx)
		if err != nil {
			return httpext.NewErrorFromStatus(err, "getting game state")
		}

		return c.JSON(http.StatusOK, gs)
	}
}
