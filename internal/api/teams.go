package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/samber/lo"

	"github.com/c4t-but-s4d/fastad/internal/models"
	"github.com/c4t-but-s4d/fastad/pkg/httpext"
	checkerpb "github.com/c4t-but-s4d/fastad/pkg/proto/checker"
	teamspb "github.com/c4t-but-s4d/fastad/pkg/proto/data/teams"
)

func (s *Service) HandleTeamsList() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := httpext.ContextFromEcho(c)

		teams, err := s.teamsClient.List(ctx)
		if err != nil {
			return httpext.NewErrorFromStatus(err, "listing teams")
		}

		resp := &teamspb.Team_Batch{
			Teams: lo.Map(teams, func(team *teamspb.Team, _ int) *teamspb.Team {
				teamCloned := team.CloneVT()
				teamCloned.Token = ""
				return teamCloned
			}),
		}

		return ProtoJSON(c, http.StatusOK, resp)
	}
}

func (s *Service) HandleTeamHistory() echo.HandlerFunc {
	type request struct {
		TeamID    int `param:"team_id"`
		ServiceID int `query:"service_id"`
		Limit     int `query:"limit"`
	}

	const (
		defaultLimit = 100
		maxLimit     = 500
	)

	// TODO: cache this handler.
	return func(c echo.Context) error {
		ctx := httpext.ContextFromEcho(c)

		req := new(request)
		if err := c.Bind(req); err != nil {
			return httpext.NewErrorf(
				http.StatusBadRequest,
				"binding request: %v",
				err,
			)
		}

		if req.TeamID == 0 {
			return httpext.NewErrorf(http.StatusBadRequest, "team_id is required")
		}

		limit := req.Limit
		if limit <= 0 {
			limit = defaultLimit
		}

		query := s.db.
			NewSelect().
			Model(&models.CheckerExecution{}).
			Where("team_id = ?", req.TeamID).
			OrderExpr("created_at DESC").
			Limit(min(maxLimit, limit))

		if req.ServiceID != 0 {
			query.Where("service_id = ?", req.ServiceID)
		}

		var executions []*models.CheckerExecution
		if err := query.Scan(ctx, &executions); err != nil {
			return httpext.NewErrorf(
				http.StatusInternalServerError,
				"listing executions: %v",
				err,
			)
		}

		resp := &checkerpb.Execution_Batch{
			Executions: lo.Map(
				executions,
				func(execution *models.CheckerExecution, _ int) *checkerpb.Execution {
					execution.Private = ""
					execution.Command = ""
					return execution.ToProto()
				},
			),
		}

		return ProtoJSON(c, http.StatusOK, resp)
	}
}

func (s *Service) HandleTeamUpdate() echo.HandlerFunc {
	type request struct {
		AvatarURL *string `json:"avatar_url"`
	}

	return func(c echo.Context) error {
		ctx := httpext.ContextFromEcho(c)

		req := new(request)
		if err := c.Bind(req); err != nil {
			return httpext.NewErrorf(http.StatusBadRequest, "binding request: %v", err)
		}

		teamToken := c.Request().Header.Get(TeamTokenHeader)
		currentTeam, err := s.teamsClient.GetByToken(ctx, teamToken)
		if err != nil {
			return httpext.NewErrorFromStatus(err, "getting current team")
		}

		updateReq := &teamspb.UpdateRequest{Id: currentTeam.GetId()}
		if req.AvatarURL != nil {
			updateReq.AvatarUrl = *req.AvatarURL
		}
		team, err := s.teamsClient.Update(ctx, updateReq)
		if err != nil {
			return httpext.NewErrorFromStatus(err, "updating team")
		}

		team.Token = ""

		return ProtoJSON(c, http.StatusOK, team)
	}
}
