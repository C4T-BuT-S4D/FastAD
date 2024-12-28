package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/c4t-but-s4d/fastad/pkg/httpext"
	teamspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/teams"
)

func (s *Service) HandleTeamsList() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := httpext.ContextFromEcho(c)

		teams, err := s.teamsClient.List(ctx)
		if err != nil {
			return httpext.NewErrorFromStatus(err, "listing teams")
		}

		// TODO: return proto from clients.
		resp := &teamspb.Team_Batch{
			Teams: lo.Map(teams, func(team *teamspb.Team, _ int) *teamspb.Team {
				teamCloned := team.CloneVT()
				teamCloned.Token = ""
				return teamCloned
			}),
		}

		raw, err := protojson.MarshalOptions{EmitUnpopulated: true}.Marshal(resp)
		if err != nil {
			return httpext.NewErrorf(
				http.StatusInternalServerError,
				"marshaling teams: %v",
				err,
			)
		}

		return c.JSONBlob(http.StatusOK, raw)
	}
}
