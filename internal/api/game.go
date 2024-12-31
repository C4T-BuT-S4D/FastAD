package api

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/c4t-but-s4d/fastad/pkg/httpext"
)

func (s *Service) HandleGetGameState() echo.HandlerFunc {
	type request struct {
		Format string `query:"format"`
	}

	return func(c echo.Context) error {
		ctx := httpext.ContextFromEcho(c)

		var req request
		if err := c.Bind(&req); err != nil {
			return httpext.NewErrorf(http.StatusBadRequest, "parsing request: %v", err)
		}

		gs, err := s.gameStateClient.Get(ctx)
		if err != nil {
			return httpext.NewErrorFromStatus(err, "getting game state")
		}

		switch req.Format {
		case "proto":
			return ProtoRaw(c, http.StatusOK, gs)
		default:
			return ProtoJSON(c, http.StatusOK, gs)
		}
	}
}
