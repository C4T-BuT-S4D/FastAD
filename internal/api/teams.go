package api

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/c4t-but-s4d/fastad/pkg/httpext"
)

func (s *Service) HandleTeamsList() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := httpext.ContextFromEcho(c)

		teams, err := s.teamsClient.List(ctx)
		if err != nil {
			return httpext.NewErrorFromStatus(err, "listing teams")
		}

		for _, team := range teams {
			team.Token = ""
		}

		return c.JSON(http.StatusOK, teams)
	}
}
